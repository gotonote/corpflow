# CorpFlow

<p align="center">
  <img src="docs/logo.jpg" width="200" alt="CorpFlow Logo">
</p>

**多智能体协作平台**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Platform](https://img.shields.io/badge/Platform-Flutter-blue.svg)](https://flutter.dev)

> **English**: [README.md](./README.md)

---

## 概述

CorpFlow 是一个**多智能体协作平台**，支持：
- 创建和管理 AI 智能体
- 可视化流程编排（拖拽操作）
- 多渠道部署（飞书、微信、Telegram、Discord）
- 多模型投票决策

---

## 功能

| 功能 | 说明 |
|------|------|
| 🤖 **智能体** | 创建自定义AI智能体，支持多种模型 |
| 🔀 **流程编排** | 可视化工作流自动化 |
| 💬 **多渠道** | 飞书、微信、Telegram、Discord |
| 🗳️ **多模型投票** | 多AI讨论并投票选择最佳答案 |
| 📱 **移动应用** | iOS、Android、macOS、Windows、iPadOS |

---

## 支持的AI模型

| 模型 | 供应商 | 环境变量 |
|------|--------|----------|
| GPT-4 | OpenAI | `OPENAI_API_KEY` |
| Claude 3 | Anthropic | `ANTHROPIC_API_KEY` |
| GLM-4 | 智谱 | `ZHIPU_API_KEY` |
| Kimi | 月之暗面 | `KIMI_API_KEY` |
| Qwen | 阿里通义千问 | `DASHSCOPE_API_KEY` |
| DeepSeek | 深度求索 | `DEEPSEEK_API_KEY` |
| MiniMax | MiniMax | `MINIMAX_API_KEY` |

---

## 快速开始

### 后端 (Go + Docker)

```bash
# 克隆仓库
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# 复制配置
cp .env.example .env

# 编辑 .env 添加 API Key
vim .env

# 使用 Docker 启动
docker-compose up -d
```

### 移动端 (Flutter)

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

---

## 如何使用

### 1. 打开浏览器

启动 Docker 后，访问：**http://localhost:3000**

---

### 2. 配置 AI 模型

在 **设置** 中：
- 点击选择模型（如 GPT-4、GLM-4、Kimi 等）
- 填入对应的 API Key
- 保存

---

### 3. 使用功能

| 功能 | 说明 |
|------|------|
| 💬 **对话** | 选择智能体，发送消息，获得AI回复 |
| 🔀 **流程** | 创建流程，拖拽节点，连接执行 |
| 🤖 **智能体** | 创建/管理不同模型的AI智能体 |
| ⚙️ **设置** | 配置模型、API Key、投票、渠道 |

---

### 4. 连接移动端（可选）

移动端连接本地服务器：
1. 确保手机和电脑在同一WiFi
2. 查看电脑IP: `ipconfig` (Windows) / `ifconfig` (Mac/Linux)
3. 在移动端设置中填入: `http://你的IP:8080`

---

### 5. 多模型投票

在 **设置** → 多模型投票 中启用

**投票方式：**
- **综合评分**：按准确性+完整性+清晰度+创造性评分
- **交叉评估**：模型互相评估
- **按长度**：简单按回复长度

**评分权重：**
- 准确性 - 30%
- 完整性 - 30%
- 清晰度 - 20%
- 创造性 - 20%

---

## 常见问题

| 问题 | 解决 |
|------|------|
| 无法访问 localhost:3000 | 检查 Docker 是否运行: `docker ps` |
| API 调用失败 | 确认已在设置中配置 API Key |
| 移动端连不上 | 检查防火墙/确保同一网络 |

---

## 环境变量

```bash
# AI 模型
export OPENAI_API_KEY=sk-xxx
export ANTHROPIC_API_KEY=sk-ant-xxx
export ZHIPU_API_KEY=xxx
export KIMI_API_KEY=xxx
export DASHSCOPE_API_KEY=xxx
export DEEPSEEK_API_KEY=xxx
export MINIMAX_API_KEY=xxx

# 渠道
export FEISHU_APP_ID=xxx
export FEISHU_APP_SECRET=xxx
export WECHAT_APP_ID=xxx
export TELEGRAM_BOT_TOKEN=xxx
```

---

## API 接口

| 方法 | 端点 | 说明 |
|------|------|------|
| POST | `/api/agents` | 创建智能体 |
| GET | `/api/agents` | 列出智能体 |
| POST | `/api/flows` | 创建流程 |
| POST | `/api/flows/:id/execute` | 执行流程 |
| POST | `/webhook/feishu` | 飞书 webhook |

---

## 许可证

MIT License
