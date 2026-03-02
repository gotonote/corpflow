# CorpFlow — Multi-Agent Collaboration Platform

<p align="center">
  <strong>Enterprise AI Collaboration Platform Built on OpenClaw</strong>
</p>

<p align="center">
  <a href="README_zh.md">中文</a> | English
</p>

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                         User Layer                               │
│         Feishu │ WeChat │ Telegram │ Web Frontend            │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    CorpFlow Frontend (React)                   │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│   │   Home   │  │   Chat   │  │  Flows   │  │   OKR    │     │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘     │
└────────────────────────────┬────────────────────────────────────┘
                             │ REST API
┌────────────────────────────▼────────────────────────────────────┐
│                   CorpFlow Backend (Express)                    │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│   │   Flow   │  │   OKR    │  │ Executor │                   │
│   │   Management│          │  │          │                   │
│   └──────────┘  └──────────┘  └──────────┘                   │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    OpenClaw Gateway                             │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    Agent Scheduler                       │   │
│   │   🤖 CEO   │   👔 Manager   │   💻 Worker              │   │
│   └─────────────────────────────────────────────────────────┘   │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    Tools Layer                           │   │
│   │   exec │ browser │ file │ web_search │ feishu_*       │   │
│   └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## ✨ Features

| Feature | Description |
|---------|-------------|
| 🤖 **Multi-Agent** | CEO → Manager → Worker three-tier collaboration |
| 🔀 **Flow Orchestration** | Visual drag-and-drop workflow builder |
| 📊 **OKR Management** | Goal setting and progress tracking |
| ⏰ **Scheduled Tasks** | Auto-generated reports and reminders |
| 💬 **Multi-Channel** | Feishu/WeChat/Telegram/Discord |

---

## 🚀 Quick Start

### Choose Deployment Option

| Option | AI Chat | Flow Editor | Complexity |
|--------|---------|-------------|------------|
| **Skill Only** | ✅ | ❌ | Easy |
| **Skill + Frontend** | ✅ | ✅ | Medium |

#### Option 1: Skill Only (AI Chat)

```bash
# Copy skill to OpenClaw
cp -r corpflow-skill /home/admin/.openclaw/workspace/skills/
```

Then chat with OpenClaw in Feishu:
```
@OpenClaw 帮我制定本季度 OKR
```

#### Option 2: Skill + Frontend (Full Features)

```bash
# 1. Install OpenClaw (if not installed)
# See: https://docs.openclaw.ai

# 2. Copy skill
cp -r corpflow-skill /home/admin/.openclaw/workspace/skills/

# 3. Install frontend dependencies
cd frontend
npm install

# 4. Install backend dependencies
cd ../backend
npm install
```

### Requirements

| Component | Version |
|-----------|---------|
| Node.js | 18+ |
| npm | 9+ |
| OpenClaw | Latest |

### Installation

```bash
# Clone project
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# Copy skill to OpenClaw
cp -r corpflow-skill ~/.openclaw/workspace/skills/

# Install frontend dependencies
cd frontend
npm install

# Install backend dependencies
cd ../backend
npm install
```

### Start Services

```bash
# Start backend (Terminal 1)
cd backend
npm run dev
# → Running on http://localhost:8081

# Start frontend (Terminal 2)
cd frontend
npm run dev
# → Running on http://localhost:3000
```

### Usage

1. Open http://localhost:3000
2. Click **💬 Chat** to start conversation
3. Click **🔀 Flows** to create workflow
4. Or chat directly with OpenClaw in Feishu

---

## 📋 Project Structure

```
corpflow/
├── corpflow-skill/     # OpenClaw skill (AI chat)
├── frontend/           # React frontend
│   ├── src/
│   │   ├── App.tsx         # Main app
│   │   ├── Chat.tsx       # Chat component
│   │   ├── FlowEditor.tsx # Flow editor
│   │   └── api.ts         # API calls
│   └── package.json
│
├── backend/            # Express backend
│   ├── server.js       # API service
│   └── package.json
│
├── CORPFLOW.md        # System Prompt
├── docs/
│   └── TECH_PLAN.md    # Technical plan
└── README.md
```

---

## 🔧 Configuration

### Feishu Integration

OpenClaw already has Feishu configured. See [OpenClaw Docs](https://docs.openclaw.ai).

### Environment Variables

```bash
# Backend (backend/.env)
PORT=8081

# Frontend (optional)
VITE_API_URL=http://localhost:8081
```

---

## ⏰ Scheduled Tasks

Auto-configured tasks:

| Task | Time | Description |
|------|------|-------------|
| `corpflow-okr-reminder` | Weekdays 9:00 | OKR progress reminder |
| `corpflow-weekly-summary` | Friday 18:00 | Weekly summary |

---

## 🔌 API Endpoints

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/flows` | Get flow list |
| POST | `/api/flows` | Save flow |
| POST | `/api/flows/:id/execute` | Execute flow |
| GET | `/api/okr` | Get OKR list |
| POST | `/api/chat/messages` | Send message |

---

## 🛠️ Tech Stack

| Layer | Technology |
|--------|------------|
| Frontend | React + TypeScript + Vite + React Flow |
| Backend | Express + Node.js |
| AI Engine | OpenClaw |
| Channels | Feishu/WeChat/Telegram/Discord |

---

## 📄 License

MIT License

---

*Built on OpenClaw*
