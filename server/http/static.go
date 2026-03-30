package http

const indexHTML = `<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
    <title>Weblise - Remote Desktop</title>
    <style>
        * { margin: 0; padding: 0; box-sizing: border-box; }
        body { font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', Roboto, sans-serif; background: #1a1a2e; color: #eee; min-height: 100vh; }
        .container { max-width: 1200px; margin: 0 auto; padding: 20px; }
        h1 { text-align: center; margin-bottom: 30px; color: #4ecca3; }
        #connection-panel { background: #16213e; padding: 20px; border-radius: 8px; margin-bottom: 20px; display: flex; gap: 10px; align-items: center; }
        #agent-key { flex: 1; padding: 12px; border: 1px solid #0f3460; border-radius: 4px; background: #0f3460; color: #eee; font-size: 14px; }
        #agent-key:focus { outline: none; border-color: #4ecca3; }
        #connect-btn { padding: 12px 24px; background: #4ecca3; color: #1a1a2e; border: none; border-radius: 4px; cursor: pointer; font-weight: bold; font-size: 14px; }
        #connect-btn:hover { background: #3db892; }
        #connect-btn:disabled { background: #555; cursor: not-allowed; }
        #status { padding: 10px; border-radius: 4px; margin-bottom: 20px; text-align: center; font-size: 14px; }
        .status-disconnected { background: #e94560; }
        .status-connecting { background: #f39c12; }
        .status-connected { background: #4ecca3; color: #1a1a2e; }
        #screen-container { display: none; background: #000; border-radius: 8px; overflow: hidden; }
        #screen-canvas { display: block; width: 100%; height: auto; cursor: crosshair; }
        .info { text-align: center; color: #666; margin-top: 20px; font-size: 12px; }
    </style>
</head>
<body>
    <div class="container">
        <h1>Weblise Remote Desktop</h1>
        <div id="connection-panel">
            <input type="text" id="agent-key" placeholder="Enter Agent Key" autocomplete="off">
            <button id="connect-btn">Connect</button>
        </div>
        <div id="status" class="status-disconnected">Disconnected</div>
        <div id="screen-container">
            <canvas id="screen-canvas"></canvas>
        </div>
        <p class="info">Enter your agent key to connect to a remote device</p>
    </div>
    <script>
        const wsUrl = window.location.protocol === 'https:'
            ? 'wss://' + window.location.host + '/client'
            : 'ws://' + window.location.hostname + ':8443/client';

        class WebliseClient {
            constructor() {
                this.ws = null;
                this.connected = false;
                this.canvas = document.getElementById('screen-canvas');
                this.ctx = this.canvas.getContext('2d');
                this.bindEvents();
            }

            bindEvents() {
                document.getElementById('connect-btn').addEventListener('click', () => this.connect());
                document.getElementById('agent-key').addEventListener('keypress', (e) => {
                    if (e.key === 'Enter') this.connect();
                });
                this.canvas.addEventListener('mousemove', (e) => this.sendMouseMove(e));
                this.canvas.addEventListener('mousedown', (e) => this.sendMouseClick(e, 'down'));
                this.canvas.addEventListener('mouseup', (e) => this.sendMouseClick(e, 'up'));
            }

            connect() {
                const key = document.getElementById('agent-key').value.trim();
                if (!key) return;

                this.updateStatus('connecting', 'Connecting...');
                document.getElementById('connect-btn').disabled = true;

                try {
                    this.ws = new WebSocket(wsUrl);

                    this.ws.onopen = () => {
                        this.ws.send(JSON.stringify({ type: 'connect', data: { agent_key: key } }));
                    };

                    this.ws.onmessage = (event) => this.handleMessage(event);

                    this.ws.onclose = () => {
                        this.connected = false;
                        this.updateStatus('disconnected', 'Disconnected');
                        document.getElementById('connect-btn').disabled = false;
                        document.getElementById('screen-container').style.display = 'none';
                    };

                    this.ws.onerror = () => {
                        this.updateStatus('disconnected', 'Connection error');
                    };

                    this.startHeartbeat();
                } catch (error) {
                    this.updateStatus('disconnected', 'Connection failed');
                    document.getElementById('connect-btn').disabled = false;
                }
            }

            handleMessage(event) {
                try {
                    const msg = JSON.parse(event.data);
                    switch (msg.type) {
                        case 'connect':
                            this.connected = true;
                            this.updateStatus('connected', 'Connected');
                            document.getElementById('screen-container').style.display = 'block';
                            break;
                        case 'frame':
                            this.renderFrame(msg.data);
                            break;
                        case 'error':
                            this.updateStatus('disconnected', 'Error');
                            break;
                    }
                } catch (error) {}
            }

            renderFrame(data) {
                const img = new Image();
                img.onload = () => {
                    this.canvas.width = img.width;
                    this.canvas.height = img.height;
                    this.ctx.drawImage(img, 0, 0);
                };
                img.src = 'data:image/jpeg;base64,' + data;
            }

            sendMouseMove(e) {
                if (!this.connected) return;
                const rect = this.canvas.getBoundingClientRect();
                const scaleX = this.canvas.width / rect.width;
                const scaleY = this.canvas.height / rect.height;
                const x = Math.round((e.clientX - rect.left) * scaleX);
                const y = Math.round((e.clientY - rect.top) * scaleY);
                this.ws.send(JSON.stringify({ type: 'input', data: { action: 'mousemove', data: { x, y } } }));
            }

            sendMouseClick(e, action) {
                if (!this.connected) return;
                const rect = this.canvas.getBoundingClientRect();
                const scaleX = this.canvas.width / rect.width;
                const scaleY = this.canvas.height / rect.height;
                const x = Math.round((e.clientX - rect.left) * scaleX);
                const y = Math.round((e.clientY - rect.top) * scaleY);
                this.ws.send(JSON.stringify({ type: 'input', data: { action: 'mouse' + action, data: { x, y, button: e.button } } }));
            }

            updateStatus(state, message) {
                const statusEl = document.getElementById('status');
                statusEl.textContent = message;
                statusEl.className = 'status-' + state;
            }

            startHeartbeat() {
                setInterval(() => {
                    if (this.ws && this.ws.readyState === WebSocket.OPEN) {
                        this.ws.send(JSON.stringify({ type: 'ping' }));
                    }
                }, 30000);
            }
        }

        const client = new WebliseClient();
    </script>
</body>
</html>
`
