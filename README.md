# Weblise

Web-based remote desktop system - MVP

## Quick Start

### Server

```bash
cd server
go run main.go serve
```

Server runs on:
- HTTP: http://localhost:8080
- WebSocket: ws://localhost:8443

### Agent

```bash
cd agent
go run main.go --server=ws://localhost:8443 --key=test-device-key
```

## Architecture

- **Agent**: Go binary for screen capture and input control
- **Server**: Go server for routing messages between clients and agents
- **Web Client**: Browser-based remote control interface

## Development Status

- [x] Design specification
- [ ] Foundation
- [ ] Agent development
- [ ] Server development
- [ ] Web client development
- [ ] Docker deployment
- [ ] End-to-end testing
