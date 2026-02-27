# CorpFlow

**Multi-Agent Collaboration Platform** | 多智能体协作平台

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Flutter-blue.svg)](https://flutter.dev)
[![AI Models](https://img.shields.io/badge/AI-Models-GPT--4%20%7C%20Claude%20%7C%20GLM%20%7C%20Kimi-green.svg)](https://github.com/gotonote/corpflow)

---

## English

### What is CorpFlow?

CorpFlow is a **multi-agent collaboration platform** that enables you to:
- Create and manage AI agents
- Build visual workflows with drag-and-drop
- Deploy across multiple channels (Feishu, WeChat, Telegram, Discord)
- Use multiple AI models with intelligent voting

### Features

| Feature | Description |
|---------|-------------|
| 🤖 **AI Agents** | Create custom AI agents with different models |
| 🔀 **Flow Builder** | Visual workflow automation |
| 💬 **Multi-Channel** | Deploy on Feishu, WeChat, Telegram, Discord |
| 🗳️ **Multi-Model Voting** | Let multiple AI models discuss and vote |
| 📱 **Mobile App** | iOS, Android, macOS, Windows, iPadOS |

### Supported AI Models

| Model | Provider | Env Variable |
|-------|----------|--------------|
| GPT-4 | OpenAI | `OPENAI_API_KEY` |
| Claude 3 | Anthropic | `ANTHROPIC_API_KEY` |
| GLM-4 | Zhipu (智谱) | `ZHIPU_API_KEY` |
| Kimi | Moonshot (月之暗面) | `KIMI_API_KEY` |
| Qwen | Alibaba (通义千问) | `DASHSCOPE_API_KEY` |
| DeepSeek | DeepSeek | `DEEPSEEK_API_KEY` |
| MiniMax | MiniMax | `MINIMAX_API_KEY` |

### Quick Start

#### 1. Backend (Go + Docker)

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

#### 2. Mobile App (Flutter)

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

# Build for Windows
flutter build windows --release
```

### How to Use Each Feature

#### 💬 Chat / 对话

1. Tap **"New Chat"** button
2. Type your message in the input field
3. AI responds instantly
4. Conversation is saved automatically

**Multi-channel**: Connect Feishu/WeChat/Telegram in Settings → Channels

#### 🔀 Flow / 流程编排

1. Go to **Flow** tab
2. Tap **"+"** to create new flow
3. **Add nodes**:
   - **Trigger**: Message trigger, schedule, webhook
   - **Agent**: AI agent node
   - **Tool**: Browser, search, calculator
   - **Condition**: Branch logic
4. **Connect nodes** by dragging from output to input
5. **Save** your flow
6. **Execute** by tapping play button

#### 🤖 Agents / 智能体

1. Go to **Agents** tab
2. Tap **"+"** to create new agent
3. Configure:
   - Name your agent
   - Select AI model (GPT-4/Claude/GLM/Kimi/Qwen/DeepSeek)
   - Set system prompt
   - Enable tools (browser, search, etc.)
4. Save and use in flows or chat

#### 🗳️ Multi-Model Voting / 多模型投票

Enable in **Settings** → Multi-Model Voting

**How it works:**
1. Enable voting toggle
2. Select voting method:
   - **Comprehensive**: Scores by Accuracy + Completeness + Clarity + Creativity
   - **Cross-evaluation**: Models evaluate each other
   - **Length**: Simple by response length
3. When enabled, multiple AI models will respond
4. System automatically selects the best response

**Scoring weights:**
- Accuracy (准确性) - 30%
- Completeness (完整性) - 30%
- Clarity (清晰度) - 20%
- Creativity (创造性) - 20%

### API Documentation

#### Create Agent
```bash
POST /api/agents
{
  "name": "Assistant",
  "model_provider": "openai",
  "model_name": "gpt-4",
  "tools": ["search", "browser"]
}
```

#### Create Flow
```bash
POST /api/flows
{
  "name": "User Support",
  "nodes": [...],
  "edges": [...]
}
```

#### Execute Flow
```bash
POST /api/flows/:id/execute
{
  "input": "User question",
  "user_id": "user123"
}
```

---

## 中文

### 什么是 CorpFlow？

CorpFlow 是一个**多智能体协作平台**，支持：

- 创建和管理 AI 智能体
- 可视化流程编排（拖拽操作）
- 多渠道部署（飞书、微信、Telegram、Discord）
- 多模型投票决策

### 功能一览

| 功能 | 说明 |
|------|------|
| 🤖 **智能体** | 创建自定义 AI 智能体，支持多种模型 |
| 🔀 **流程编排** | 可视化工作流自动化 |
| 💬 **多渠道** | 飞书、微信、Telegram、Discord |
| 🗳️ **多模型投票** | 多AI模型讨论并投票选择最佳答案 |
| 📱 **移动应用** | iOS、Android、macOS、Windows、iPadOS |

### 快速开始

#### 1. 后端 (Go + Docker)

```bash
# 克隆仓库
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# 复制配置
cp .env.example .env

# 编辑 .env 添加你的 API Key
vim .env

# 使用 Docker 启动
docker-compose up -d
```

#### 2. 移动端 (Flutter)

```bash
cd mobile

# 安装依赖
flutter pub get

# 开发运行
flutter run

# 构建 Android
flutter build apk --release

# 构建 iOS (仅 macOS)
flutter build ios --release
```

### 各功能使用指南

#### 💬 对话

1. 点击 **"新建对话"** 按钮
2. 在输入框输入消息
3. AI 即时回复
4. 对话自动保存

**多渠道配置**：设置 → 渠道 → 开启飞书/微信/Telegram

#### 🔀 流程编排

1. 进入 **流程** 标签
2. 点击 **"+"** 创建新流程
3. **添加节点**：
   - **触发器**：消息触发、定时任务、Webhook
   - **智能体**：AI 节点
   - **工具**：浏览器、搜索、计算器
   - **条件**：分支逻辑
4. **连接节点**：从输出拖拽到输入
5. **保存**流程
6. 点击播放按钮**执行**

#### 🤖 智能体

1. 进入 **智能体** 标签
2. 点击 **"+"** 创建新智能体
3. 配置：
   - 设置名称
   - 选择 AI 模型 (GPT-4/Claude/GLM/Kimi/Qwen/DeepSeek)
   - 设置系统提示词
   - 启用工具（浏览器、搜索等）
4. 保存后在流程或对话中使用

#### 🗳️ 多模型投票

在 **设置** → 多模型投票 中启用

**工作原理：**
1. 开启投票开关
2. 选择投票方式：
   - **综合评分**：按准确性+完整性+清晰度+创造性评分
   - **交叉评估**：模型互相评估
   - **按长度**：简单按回复长度
3. 启用后，多个 AI 模型会同时响应
4. 系统自动选择最佳答案

**评分权重：**
- 准确性 (Accuracy) - 30%
- 完整性 (Completeness) - 30%
- 清晰度 (Clarity) - 20%
- 创造性 (Creativity) - 20%

### 环境变量配置

```bash
# OpenAI
export OPENAI_API_KEY=sk-xxx

# Anthropic
export ANTHROPIC_API_KEY=sk-ant-xxx

# 智谱 GLM
export ZHIPU_API_KEY=xxx

# Kimi (月之暗面)
export KIMI_API_KEY=xxx

# 通义千问 (阿里)
export DASHSCOPE_API_KEY=xxx

# DeepSeek
export DEEPSEEK_API_KEY=xxx

# MiniMax
export MINIMAX_API_KEY=xxx

# 飞书
export FEISHU_APP_ID=xxx
export FEISHU_APP_SECRET=xxx

# 微信
export WECHAT_APP_ID=xxx
export WECHAT_APP_SECRET=xxx

# Telegram
export TELEGRAM_BOT_TOKEN=xxx
```

### Docker 
# docker-compose部署

```yaml.yml 已配置以下服务：
# - server: Go 后端 (端口 8080)
# - frontend: React 前端 (端口 3000)
# - db: PostgreSQL 数据库
# - redis: Redis 缓存
```

启动所有服务：
```bash
docker-compose up -d
```

访问：
- 前端：http://localhost:3000
- API：http://localhost:8080/api

---

## License

MIT License - feel free to use and modify!
