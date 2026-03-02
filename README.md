# CorpFlow — 多智能体协作平台

<p align="center">
  <strong>基于 OpenClaw 二次开发的企业级 AI 协作平台</strong>
</p>

---

## 🏗️ 架构图

```
┌─────────────────────────────────────────────────────────────────┐
│                         用户层                                   │
│         飞书 │ 微信 │ Telegram │ Web 前端                     │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    CorpFlow 前端 (React)                        │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐     │
│   │   Home   │  │   Chat   │  │  Flows   │  │  OKR    │     │
│   └──────────┘  └──────────┘  └──────────┘  └──────────┘     │
└────────────────────────────┬────────────────────────────────────┘
                             │ REST API
┌────────────────────────────▼────────────────────────────────────┐
│                   CorpFlow 后端 (Express)                       │
│   ┌──────────┐  ┌──────────┐  ┌──────────┐                   │
│   │  流程管理  │  │   OKR   │  │  执行器  │                   │
│   └──────────┘  └──────────┘  └──────────┘                   │
└────────────────────────────┬────────────────────────────────────┘
                             │
┌────────────────────────────▼────────────────────────────────────┐
│                    OpenClaw Gateway                             │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    Agent 调度层                          │   │
│   │   🤖 CEO   │   👔 Manager   │   💻 Worker              │   │
│   └─────────────────────────────────────────────────────────┘   │
│   ┌─────────────────────────────────────────────────────────┐   │
│   │                    工具层                                  │   │
│   │   exec │ browser │ file │ web_search │ feishu_*       │   │
│   └─────────────────────────────────────────────────────────┘   │
└─────────────────────────────────────────────────────────────────┘
```

---

## ✨ 功能特性

| 功能 | 说明 |
|------|------|
| 🤖 **多智能体** | CEO → Manager → Worker 三层协作 |
| 🔀 **流程编排** | 可视化拖拽工作流 |
| 📊 **OKR 管理** | 目标设定与进度追踪 |
| ⏰ **定时任务** | 自动生成报告和提醒 |
| 💬 **多渠道** | 飞书/微信/Telegram/Discord |

---

## 🚀 快速开始

### 环境要求

| 组件 | 版本 |
|------|------|
| Node.js | 18+ |
| npm | 9+ |

### 安装

```bash
# 克隆项目
git clone https://github.com/gotonote/corpflow.git
cd corpflow

# 安装前端依赖
cd frontend
npm install

# 安装后端依赖
cd ../backend
npm install
```

### 启动

```bash
# 启动后端 (终端 1)
cd backend
npm run dev
# → 运行在 http://localhost:8081

# 启动前端 (终端 2)
cd frontend
npm run dev
# → 运行在 http://localhost:3000
```

### 使用

1. 打开 http://localhost:3000
2. 点击 **💬 Chat** 开始对话
3. 点击 **🔀 Flows** 创建工作流
4. 或直接在飞书与 OpenClaw 对话

---

## 📋 项目结构

```
corpflow/
├── frontend/           # React 前端
│   ├── src/
│   │   ├── App.tsx         # 主应用
│   │   ├── Chat.tsx        # 聊天组件
│   │   ├── FlowEditor.tsx  # 流程编辑器
│   │   └── api.ts          # API 调用
│   └── package.json
│
├── backend/            # Express 后端
│   ├── server.js       # API 服务
│   └── package.json
│
├── CORPFLOW.md        # System Prompt
├── docs/
│   └── TECH_PLAN.md    # 技术方案
└── README.md
```

---

## 🔧 配置

### 飞书集成

OpenClaw 已配置飞书，详见 [OpenClaw 文档](https://docs.openclaw.ai)。

### 环境变量

```bash
# 后端 (backend/.env)
PORT=8081

# 前端 (可选)
VITE_API_URL=http://localhost:8081
```

---

## ⏰ 定时任务

已配置以下自动任务：

| 任务 | 时间 | 说明 |
|------|------|------|
| `corpflow-okr-reminder` | 工作日 9:00 | OKR 进度提醒 |
| `corpflow-weekly-summary` | 周五 18:00 | 周总结 |

---

## 🔌 API 接口

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/api/flows` | 获取流程列表 |
| POST | `/api/flows` | 保存流程 |
| POST | `/api/flows/:id/execute` | 执行流程 |
| GET | `/api/okr` | 获取 OKR 列表 |
| POST | `/api/chat/messages` | 发送消息 |

---

## 🛠️ 技术栈

| 层级 | 技术 |
|------|------|
| 前端 | React + TypeScript + Vite + React Flow |
| 后端 | Express + Node.js |
| AI 引擎 | OpenClaw |
| 渠道 | 飞书/微信/Telegram/Discord |

---

## 📄 许可证

MIT License

---

*基于 OpenClaw 二次开发*
