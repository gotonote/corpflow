# CorpFlow

<p align="center">
  <img src="docs/logo.jpg" width="200" alt="CorpFlow Logo">
</p>

**Multi-Agent Collaboration Platform**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Flutter-blue.svg)](https://flutter.dev)

> **中文**: [README_zh.md](./README_zh.md)

---

## Overview

CorpFlow is a **multi-agent collaboration platform** that enables you to:
- Create and manage AI agents
- Build visual workflows with drag-and-drop
- Deploy across multiple channels (Feishu, WeChat, Telegram, Discord)
- Use multiple AI models with intelligent voting

---

## Features

| Feature | Description |
|---------|-------------|
| 🤖 **AI Agents** | Create custom AI agents with different models |
| 🔀 **Flow Builder** | Visual workflow automation |
| 💬 **Multi-Channel** | Feishu, WeChat, Telegram, Discord |
| 🗳️ **Multi-Model Voting** | Multiple AI models discuss and vote |
| 📱 **Mobile App** | iOS, Android, macOS, Windows, iPadOS |

---

## Supported AI Models

| Model | Provider | Env Variable |
|-------|----------|--------------|
| GPT-4 | OpenAI | `OPENAI_API_KEY` |
| Claude 3 | Anthropic | `ANTHROPIC_API_KEY` |
| GLM-4 | Zhipu | `ZHIPU_API_KEY` |
| Kimi | Moonshot | `KIMI_API_KEY` |
| Qwen | Alibaba | `DASHSCOPE_API_KEY` |
| DeepSeek | DeepSeek | `DEEPSEEK_API_KEY` |
| MiniMax | MiniMax | `MINIMAX_API_KEY` |

---

## Quick Start

### Backend (Go + Docker)

```bash
# Clone the repo
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# Copy configuration
cp .env.example .env

# Edit .env with your API keys
vim .env

# Start with Docker
docker-compose up -d
```

### Mobile App (Flutter)

```bash
cd mobile

# Install dependencies
flutter pub get

# Run in development
flutter run

# Build for Android
flutter build apk --release

# Build for iOS (macOS only)
flutter build ios --release
```

---

## How to Use

### 1. Open in Browser

After starting Docker, open: **http://localhost:3000**

---

### 2. Configure AI Models

In **Settings**:
- Click to select a model (GPT-4, GLM-4, Kimi, etc.)
- Enter your API Key for that model
- Save

---

### 3. Use Features

| Feature | How to Use |
|---------|-------------|
| 💬 **Chat** | Select an agent, type message, get AI response |
| 🔀 **Flow** | Create flow, drag nodes, connect them, execute |
| 🤖 **Agents** | Create/manage AI agents with different models |
| ⚙️ **Settings** | Configure models, API keys, voting, channels |

---

### 4. Connect Mobile App (Optional)

To connect mobile app to local server:

1. Ensure phone and computer are on the same WiFi
2. Get your computer's IP:
   - Windows: `ipconfig`
   - Mac/Linux: `ifconfig`
3. In mobile app Settings, enter: `http://YOUR_IP:8080`

---

### 5. Multi-Model Voting

Enable in **Settings** → Multi-Model Voting

**Voting Methods:**
- **Comprehensive**: Scores by Accuracy + Completeness + Clarity + Creativity
- **Cross-evaluation**: Models evaluate each other
- **Length**: By response length

**Scoring weights:**
- Accuracy - 30%
- Completeness - 30%
- Clarity - 20%
- Creativity - 20%

---

## Environment Variables

```bash
# AI Models
export OPENAI_API_KEY=sk-xxx
export ANTHROPIC_API_KEY=sk-ant-xxx
export ZHIPU_API_KEY=xxx
export KIMI_API_KEY=xxx
export DASHSCOPE_API_KEY=xxx
export DEEPSEEK_API_KEY=xxx
export MINIMAX_API_KEY=xxx

# Channels
export FEISHU_APP_ID=xxx
export FEISHU_APP_SECRET=xxx
export WECHAT_APP_ID=xxx
export TELEGRAM_BOT_TOKEN=xxx
```

---

## Troubleshooting

| Problem | Solution |
|---------|----------|
| Can't access localhost:3000 | Check if Docker is running: `docker ps` |
| API calls fail | Verify API Key is configured in Settings |
| Mobile can't connect | Check firewall / ensure same network |

---

## API Endpoints

| Method | Endpoint | Description |
|--------|----------|-------------|
| POST | `/api/agents` | Create agent |
| GET | `/api/agents` | List agents |
| POST | `/api/flows` | Create flow |
| POST | `/api/flows/:id/execute` | Execute flow |
| POST | `/webhook/feishu` | Feishu webhook |

---

## License

MIT License
