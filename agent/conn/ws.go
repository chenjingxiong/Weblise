package conn

import (
	"encoding/base64"
	"log"
	"runtime"
	"sync"
	"time"

	"github.com/chenjingxiong/weblise/server/protocol"
	"github.com/chenjingxiong/weblise/agent/input"
	"github.com/chenjingxiong/weblise/agent/screen"
	"github.com/gorilla/websocket"
)

type Connection struct {
	serverAddr     string
	deviceKey      string
	name           string
	wsConn         *websocket.Conn
	capturer       screen.Capturer
	mouse          *input.Mouse
	keyboard       *input.Keyboard
	mu             sync.Mutex
	sendChan       chan *protocol.Message
	connected      bool
	screenInterval time.Duration
}

type Config struct {
	ServerAddr    string
	DeviceKey     string
	Name          string
	CaptureFPS    int
	ScreenQuality int
}

func New(cfg *Config) (*Connection, error) {
	capturer, err := screen.New(screen.DefaultConfig())
	if err != nil {
		return nil, err
	}

	interval := time.Second / time.Duration(cfg.CaptureFPS)

	return &Connection{
		serverAddr:     cfg.ServerAddr,
		deviceKey:      cfg.DeviceKey,
		name:          cfg.Name,
		capturer:      capturer,
		mouse:         input.NewMouse(),
		keyboard:      input.NewKeyboard(),
		sendChan:      make(chan *protocol.Message, 32),
		screenInterval: interval,
	}, nil
}

func (c *Connection) Connect() error {
	wsConn, _, err := websocket.DefaultDialer.Dial(c.serverAddr, nil)
	if err != nil {
		return err
	}

	c.wsConn = wsConn
	c.connected = true
	log.Printf("[Conn] Connected to %s", c.serverAddr)

	regMsg, _ := protocol.NewMessage(protocol.MessageTypeRegister, protocol.RegisterMessage{
		Key:    c.deviceKey,
		Name:   c.name,
		OSType: getOSType(),
	})

	if err := c.wsConn.WriteJSON(regMsg); err != nil {
		return err
	}

	log.Printf("[Conn] Registered with key: %s", c.deviceKey)
	return nil
}

func (c *Connection) Start() error {
	go c.sendLoop()
	go c.receiveLoop()
	go c.screenLoop()
	go c.heartbeatLoop()
	return nil
}

func (c *Connection) sendLoop() {
	for {
		select {
		case msg := <-c.sendChan:
			if !c.connected {
				return
			}
			if err := c.wsConn.WriteJSON(msg); err != nil {
				log.Printf("[Conn] Send error: %v", err)
				c.Close()
				return
			}
		}
	}
}

func (c *Connection) receiveLoop() {
	for {
		if !c.connected {
			return
		}

		var msg protocol.Message
		if err := c.wsConn.ReadJSON(&msg); err != nil {
			log.Printf("[Conn] Receive error: %v", err)
			c.Close()
			return
		}

		c.handleMessage(&msg)
	}
}

func (c *Connection) handleMessage(msg *protocol.Message) {
	switch msg.Type {
	case protocol.MessageTypeInput:
		c.handleInput(msg)
	case protocol.MessageTypePong:
		// Heartbeat response
	}
}

func (c *Connection) handleInput(msg *protocol.Message) {
	var inputMsg protocol.InputMessage
	if err := msg.ParseData(&inputMsg); err != nil {
		return
	}

	switch inputMsg.Action {
	case "mousemove":
		c.mouse.Move(inputMsg.Data.X, inputMsg.Data.Y)
	case "mousedown", "mouseup":
		button := getButtonName(inputMsg.Data.Button)
		if inputMsg.Action == "mousedown" {
			c.mouse.Press(button)
		} else {
			c.mouse.Release(button)
		}
	case "keydown":
		c.keyboard.Press(inputMsg.Data.Key)
	}
}

func (c *Connection) screenLoop() {
	ticker := time.NewTicker(c.screenInterval)
	defer ticker.Stop()

	for range ticker.C {
		if !c.connected {
			return
		}

		img, err := c.capturer.Capture()
		if err != nil {
			continue
		}

		data, err := screen.EncodeJPEG(img, 80)
		if err != nil {
			continue
		}

		base64Data := base64.StdEncoding.EncodeToString(data)

		frameMsg := protocol.FrameMessage{
			Data:      []byte(base64Data),
			Width:     img.Bounds().Dx(),
			Height:    img.Bounds().Dy(),
			Timestamp: time.Now().Unix(),
		}

		msg, _ := protocol.NewMessage(protocol.MessageTypeFrame, frameMsg)
		c.sendChan <- msg
	}
}

func (c *Connection) heartbeatLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		if !c.connected {
			return
		}

		pingMsg := protocol.HeartbeatMessage{Timestamp: time.Now().Unix()}
		msg, _ := protocol.NewMessage(protocol.MessageTypeHeartbeat, pingMsg)
		c.sendChan <- msg
	}
}

func (c *Connection) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if !c.connected {
		return
	}

	c.connected = false
	if c.wsConn != nil {
		c.wsConn.Close()
	}
	close(c.sendChan)
	log.Printf("[Conn] Connection closed")
}

func getButtonName(button int) string {
	switch button {
	case 0:
		return "left"
	case 1:
		return "middle"
	case 2:
		return "right"
	default:
		return "left"
	}
}

func getOSType() string {
	return runtime.GOOS
}
