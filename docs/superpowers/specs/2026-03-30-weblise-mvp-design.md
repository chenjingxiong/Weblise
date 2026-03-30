# Weblise MVP - 远程桌面系统设计文档

**日期:** 2026-03-30
**版本:** MVP v1.0
**状态:** 设计阶段

---

## 1. 项目概述

Weblise 是一个基于 Web 的远程桌面控制系统，类似 TeamViewer、RustDesk 或 ToDesk，但以 Web 客户端为核心。

### 1.1 MVP 目标

实现最小可行产品，跑通 "Agent → Server → Web客户端" 的完整远程桌面流程。

### 1.2 包含功能

- ✅ Agent 单文件部署（Windows/Linux）
- ✅ Server Docker 部署
- ✅ Web 客户端（浏览器访问）
- ✅ 屏幕实时显示
- ✅ 鼠标/键盘远程控制
- ✅ 简单 Key 认证模式

### 1.3 暂不包含（v2+）

- ❌ 用户账号系统
- ❌ Android/Windows 原生客户端
- ❌ P2P 打洞（先用中转）
- ❌ 节流模式
- ❌ 应用窗口级捕获

---

## 2. 系统架构

### 2.1 整体架构图

```
┌─────────────────────────────────────────────────────────────┐
│                        Docker 容器                          │
│                                                             │
│   ┌─────────────────────────────────────────────────┐      │
│   │                    Server (Go)                   │      │
│   │                                                  │      │
│   │   ┌──────────────┐         ┌──────────────┐     │      │
│   │   │  HTTP Server │         │  WS Server   │     │      │
│   │   │   :8080      │         │   :8443      │     │      │
│   │   │              │         │              │     │      │
│   │   │  • 静态文件  │         │  • Agent 连接│     │      │
│   │   │  • 反向代理  │         │  • 客户端WS  │     │      │
│   │   └──────────────┘         └──────────────┘     │      │
│   └─────────────────────────────────────────────────┘      │
└─────────────────────────────────────────────────────────────┘
           │                           │
           │                           │
    ┌──────▼────────┐          ┌───────▼────────┐
    │  反向代理      │          │   Agent 连接    │
    │  (Nginx等)    │          │   (WS :8443)   │
    │  HTTPS :443   │          │                 │
    └───────────────┘          └─────────────────┘
           │
    ┌──────▼────────┐
    │ 用户浏览器     │
    │ HTTPS 访问    │
    └───────────────┘
```

### 2.2 连接模式

MVP 版本使用服务器中转模式：

```
Agent ←──WebSocket──→ Server ←──WebSocket──→ Web Client
```

### 2.3 端口配置

| 端口 | 协议 | 用途 | 对外暴露 |
|------|------|------|----------|
| **8080** | HTTP | Web 静态文件 + API | ✅ 给反向代理 |
| **8443** | WebSocket | Agent/Client 通讯端点 | ✅ Agent连接 |

---

## 3. 技术栈

### 3.1 技术选型

| 组件 | 技术 | 说明 |
|------|------|------|
| **Agent** | Go | 单文件编译，跨平台部署 |
| **Server** | Go | 高性能 WebSocket 服务 |
| **数据库** | SQLite | 轻量，无需额外服务 |
| **Web客户端** | HTML/JS/CSS | 浏览器原生，无需框架 |

### 3.2 技术栈选择理由

- **Go**: 单文件编译、跨平台、高性能、WebSocket 生态成熟
- **SQLite**: 零配置、文件存储、适合小规模部署
- **纯前端**: MVP 阶段无需复杂框架，快速迭代

---

## 4. 组件设计

### 4.1 Agent (Go 单文件)

**职责:**
- 屏幕捕获
- 输入事件执行（鼠标/键盘）
- 与 Server 的 WebSocket 通信

**启动参数:**
```bash
./agent --server=ws://server.com:8443 --key=DEVICE_KEY
```

**模块结构:**
```
agent/
├── main.go           # 入口，参数解析
├── screen/
│   ├── capture.go    # 屏幕捕获接口
│   ├── windows.go    # Windows 实现
│   └── linux.go      # Linux 实现
├── input/
│   ├── keyboard.go   # 键盘控制
│   └── mouse.go      # 鼠标控制
└── conn/
    └── ws.go         # WebSocket 连接和消息处理
```

**屏幕捕获方案:**
- **Windows**: `github.com/kbinani/screenshot` 或 Win32 API
- **Linux**: X11 或 Wayland 协议

**编码方案:**
- MVP: JPEG/PNG 压缩，简单高效
- 后续: WebRTC/H.264 硬件编码

---

### 4.2 Server (Go)

**职责:**
- HTTP 服务（静态文件）
- WebSocket 服务（Agent 和 Client 连接）
- 消息路由（Client ↔ Agent 桥接）
- Key 认证和设备管理

**模块结构:**
```
server/
├── main.go           # 入口，启动 HTTP 和 WS 服务
├── cmd/
│   └── root.go       # 命令行参数
├── http/
│   ├── server.go     # HTTP Server (:8080)
│   └── static.go     # 嵌入的静态文件
├── ws/
│   ├── server.go     # WebSocket Server (:8443)
│   ├── agent.go      # Agent 连接处理
│   └── client.go     # Client 连接处理
├── router/
│   └── router.go     # 消息路由 (Client ↔ Agent)
└── db/
    ├── db.go         # SQLite 连接
    └── schema.sql    # 数据库表结构
```

**数据库表结构:**
```sql
-- 设备表
CREATE TABLE devices (
    id TEXT PRIMARY KEY,
    key TEXT UNIQUE NOT NULL,
    name TEXT,
    os_type TEXT,
    last_seen INTEGER,
    created_at INTEGER
);

-- 连接会话
CREATE TABLE sessions (
    id TEXT PRIMARY KEY,
    device_id TEXT,
    client_ip TEXT,
    connected_at INTEGER,
    FOREIGN KEY (device_id) REFERENCES devices(id)
);
```

---

### 4.3 Web 客户端

**职责:**
- 显示远程屏幕
- 捕获用户输入并发送
- WebSocket 连接管理

**文件结构:**
```
web/
├── index.html        # 单页面应用
├── app.js            # 前端逻辑
├── screen.js         # 屏幕渲染
├── input.js          # 输入处理
└── style.css         # 样式
```

**界面元素:**
- 设备 Key 输入框
- 连接/断开按钮
- 屏幕显示区域（Canvas 或 IMG）
- 连接状态指示

---

## 5. 通讯协议

### 5.1 连接建立流程

```
┌─────────────────────────────────────────────────────────────┐
│                      连接流程                                │
├─────────────────────────────────────────────────────────────┤
│                                                             │
│  1. Agent 启动                                              │
│     └─> WebSocket 连接到 ws://server:8443/agent            │
│     └─> 发送注册消息: { "type": "register", "key": "xxx" }  │
│                                                             │
│  2. Server 验证 Key                                         │
│     └─> 记录 Agent 在线，分配连接ID                         │
│                                                             │
│  3. 用户浏览器访问                                          │
│     └─> https://reverse-proxy.com (Nginx等)               │
│         └─> 反向代理到 http://server:8080                  │
│                                                             │
│  4. Web Client 建立 WebSocket                              │
│     └─> wss://reverse-proxy.com/ws/client (通过反向代理)   │
│     └─> 发送连接请求: { "type": "connect", "agent_key": "xxx" } │
│                                                             │
│  5. Server 桥接 Client 和 Agent                             │
│     └─> 开始转发消息                                        │
│                                                             │
└─────────────────────────────────────────────────────────────┘
```

### 5.2 消息协议

**Agent → Server → Client (屏幕数据):**
```json
{
  "type": "frame",
  "data": "base64_encoded_image_data",
  "width": 1920,
  "height": 1080,
  "timestamp": 1234567890
}
```

**Client → Server → Agent (输入事件):**
```json
{
  "type": "input",
  "action": "mousemove|mousedown|mouseup|keydown|keyup",
  "data": {
    "x": 100,
    "y": 200,
    "button": 1,
    "key": "a",
    "keyCode": 65
  }
}
```

**控制消息:**
```json
// 心跳
{ "type": "ping", "timestamp": 1234567890 }

// 断开连接
{ "type": "disconnect", "reason": "user_initiated" }

// 错误
{ "type": "error", "message": "Invalid key" }
```

---

## 6. 部署方案

### 6.1 Docker Compose 配置

```yaml
services:
  weblise-server:
    image: weblise/server:latest
    container_name: weblise
    ports:
      - "8080:8080"   # HTTP for reverse proxy
      - "8443:8443"   # WebSocket for agents
    volumes:
      - ./data:/app/data  # SQLite 数据持久化
    environment:
      - HTTP_PORT=8080
      - WS_PORT=8443
      - DB_PATH=/app/data/weblise.db
    restart: unless-stopped
```

### 6.2 反向代理配置 (Nginx 示例)

```nginx
server {
    listen 443 ssl;
    server_name remote.example.com;

    ssl_certificate /path/to/cert.pem;
    ssl_certificate_key /path/to/key.pem;

    # Web 客户端
    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    # WebSocket 升级
    location /ws/ {
        proxy_pass http://localhost:8080/ws/;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

### 6.3 Agent 部署

**Windows:**
```powershell
# 下载单文件
wget https://releases.weblise.com/agent-windows-amd64.exe -O weblise-agent.exe

# 启动
.\weblise-agent.exe --server=ws://your-server.com:8443 --key=YOUR_DEVICE_KEY
```

**Linux:**
```bash
# 下载单文件
wget https://releases.weblise.com/agent-linux-amd64 -O weblise-agent
chmod +x weblise-agent

# 启动
./weblise-agent --server=ws://your-server.com:8443 --key=YOUR_DEVICE_KEY
```

---

## 7. 安全考虑

### 7.1 MVP 阶段安全措施

- **设备 Key**: 预共享密钥认证
- **HTTPS**: 通过反向代理提供 TLS 加密
- **Key 生成**: 强随机生成，至少 32 字符

### 7.2 后续增强（v2+）

- 端到端加密
- 用户账号系统
- 设备授权管理
- 连接日志和审计
- 速率限制

---

## 8. 项目结构

```
weblise/
├── agent/              # Agent 端 (Go)
│   ├── main.go
│   ├── screen/
│   │   ├── capture.go
│   │   ├── windows.go
│   │   └── linux.go
│   ├── input/
│   │   ├── keyboard.go
│   │   └── mouse.go
│   └── conn/
│       └── ws.go
│
├── server/             # Server 端 (Go)
│   ├── main.go
│   ├── cmd/
│   │   └── root.go
│   ├── http/
│   │   ├── server.go
│   │   └── static.go
│   ├── ws/
│   │   ├── server.go
│   │   ├── agent.go
│   │   └── client.go
│   ├── router/
│   │   └── router.go
│   └── db/
│       ├── db.go
│       └── schema.sql
│
├── web/                # Web 客户端 (打包进 Server)
│   ├── index.html
│   ├── app.js
│   ├── screen.js
│   ├── input.js
│   └── style.css
│
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
│
├── docs/
│   └── superpowers/specs/
│       └── 2026-03-30-weblise-mvp-design.md
│
└── README.md           # 部署文档
```

---

## 9. 开发阶段

### Phase 1: 基础设施
- [ ] 项目脚手架搭建
- [ ] Server HTTP 服务和静态文件
- [ ] WebSocket 服务框架
- [ ] 数据库初始化

### Phase 2: Agent 开发
- [ ] 屏幕捕获（Windows）
- [ ] 屏幕捕获（Linux）
- [ ] WebSocket 连接
- [ ] 输入控制（鼠标/键盘）

### Phase 3: Server 逻辑
- [ ] Agent 注册和认证
- [ ] Client 连接处理
- [ ] 消息路由（桥接）
- [ ] 心跳和断线处理

### Phase 4: Web 客户端
- [ ] 基础 UI
- [ ] WebSocket 连接
- [ ] 屏幕显示（Canvas）
- [ ] 输入捕获和发送

### Phase 5: 部署和测试
- [ ] Docker 打包
- [ ] 端到端测试
- [ ] 文档完善

---

## 10. 后续版本规划

### v2.0 功能
- 用户账号系统
- 设备管理和授权
- P2P 打洞直连
- WebRTC 视频编码
- 节流模式

### v3.0 功能
- Android 原生客户端
- Windows 原生客户端
- 应用窗口级捕获
- 文件传输
- 剪贴板同步

---

**文档版本:** 1.0
**最后更新:** 2026-03-30
