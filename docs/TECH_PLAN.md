# CorpFlow 基于 OpenClaw 二次开发技术方案

## 一、架构设计

```
┌─────────────────────────────────────────────────────────────┐
│                        用户层                                │
│   飞书 │ 微信 │ Telegram │ Discord │ Web 前端               │
└──────────────────────────┬──────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    CorpFlow 服务层                           │
│                                                              │
│  ┌─────────────┐  ┌─────────────┐  ┌─────────────────────┐  │
│  │  Web UI     │  │  流程编辑器  │  │   OKR/绩效管理      │  │
│  │ (React)     │  │ (React Flow)│  │   (Bitable)        │  │
│  └──────┬──────┘  └──────┬──────┘  └──────────┬──────────┘  │
│         │                │                     │              │
│         └────────────────┼─────────────────────┘              │
│                          │                                    │
│                   ┌──────▼──────┐                             │
│                   │  REST API   │                             │
│                   │  (Express)  │                             │
│                   └──────┬──────┘                             │
└──────────────────────────┼────────────────────────────────────┘
                           │
┌──────────────────────────▼──────────────────────────────────┐
│                    OpenClaw Gateway                          │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                    Agent 调度层                          │ │
│  │  ┌──────────┐  ┌──────────┐  ┌──────────┐              │ │
│  │  │   CEO    │  │ Manager  │  │ Worker  │              │ │
│  │  │  Agent   │  │  Agent   │  │  Agent  │              │ │
│  │  └──────────┘  └──────────┘  └──────────┘              │ │
│  └─────────────────────────────────────────────────────────┘ │
│                                                              │
│  ┌─────────────────────────────────────────────────────────┐ │
│  │                    工具层                                │ │
│  │  exec │ browser │ file │ web_search │ feishu_*        │ │
│  └─────────────────────────────────────────────────────────┘ │
└──────────────────────────────────────────────────────────────┘
```

## 二、模块详细设计

### 2.1 核心模块

| 模块 | 技术栈 | 职责 |
|------|--------|------|
| Web 前端 | React + TypeScript | UI 界面 |
| 流程编辑器 | React Flow | 可视化工作流 |
| API 网关 | Express/FastAPI | 业务接口 |
| 数据存储 | 飞书 Bitable | OKR/绩效数据 |
| AI 脑 | OpenClaw | Agent 调度 |

### 2.2 数据模型

```typescript
// 组织架构
interface OrgStructure {
  agent_id: string
  name: string
  role: 'CEO' | 'Manager' | 'Worker'
  parent_id?: string  // 上级 Agent ID
  model: string
  system_prompt: string
}

// OKR
interface OKR {
  id: string
  objective: string
  key_results: KeyResult[]
  owner_id: string
  quarter: string  // 2024-Q1
  status: 'draft' | 'active' | 'completed'
}

interface KeyResult {
  id: string
  metric: string
  target: number
  current: number
}

// 绩效评估
interface PerformanceReview {
  id: string
  agent_id: string
  period: string
  metrics: {
    task_completion: number
    quality_score: number
    collaboration: number
  }
  ai_analysis: string
}

// 工作流
interface Flow {
  id: string
  name: string
  nodes: FlowNode[]
  edges: FlowEdge[]
  created_at: string
}

interface FlowNode {
  id: string
  type: 'agent' | 'tool' | 'condition'
  config: AgentConfig | ToolConfig
}

interface FlowEdge {
  source: string
  target: string
  condition?: string
}
```

## 三、API 接口设计

### 3.1 组织架构

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/org | 获取组织架构 |
| POST | /api/org/agent | 创建 Agent |
| PUT | /api/org/agent/:id | 更新 Agent |
| DELETE | /api/org/agent/:id | 删除 Agent |

### 3.2 OKR 管理

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/okr | 获取 OKR 列表 |
| POST | /api/okr | 创建 OKR |
| PUT | /api/okr/:id | 更新 OKR |
| GET | /api/okr/:id/progress | AI 分析进度 |

### 3.3 流程编排

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/flows | 获取流程列表 |
| POST | /api/flows | 创建流程 |
| PUT | /api/flows/:id | 更新流程 |
| POST | /api/flows/:id/execute | 执行流程 |

### 3.4 绩效评估

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | /api/performance | 评估列表 |
| POST | /api/performance/generate | AI 生成评估 |
| GET | /api/performance/:id | 查看评估详情 |

## 四、与 OpenClaw 集成

### 4.1 Agent 调用

```typescript
import { sessions_spawn } from '@openclaw/sdk'

// 调用 CEO Agent
const result = await sessions_spawn({
  agentId: 'corpflow-ceo',
  task: '分析本季度销售 OKR 进度',
  timeoutSeconds: 120
})

// 多 Agent 协作
async function multiAgentTask(task: string) {
  // 1. CEO 分解任务
  const ceoResult = await sessions_spawn({
    agentId: 'corpflow-ceo',
    task: `分解任务: ${task}`
  })
  
  // 2. 分发给 Manager
  const managerTasks = ceoResult.subtasks
  const managerResults = await Promise.all(
    managerTasks.map(t => sessions_spawn({
      agentId: 'corpflow-manager',
      task: t
    }))
  )
  
  // 3. Worker 执行
  // ...
  
  return aggregateResults(managerResults)
}
```

### 4.2 工具封装

```typescript
// 自定义工具封装
const customTools = {
  // 查询 OKR
  get_okr: async (params: { quarter?: string }) => {
    return await bitable.query('OKR', params)
  },
  
  // 更新进度
  update_progress: async (params: { kr_id: string, value: number }) => {
    return await bitable.update('KeyResults', params)
  },
  
  // 发送通知
  notify: async (params: { user_id: string, message: string }) => {
    return await message.send({
      channel: 'feishu',
      to: params.user_id,
      message: params.message
    })
  }
}
```

### 4.3 消息通道

```typescript
// 飞书消息处理
gateway.on('feishu:message', async (msg) => {
  if (msg.text.includes('OKR')) {
    const okrList = await api.getOKR()
    await message.reply(msg, formatOKR(okrList))
  }
})
```

## 五、开发计划

### Phase 1: 基础能力 (1-2 周)
- [ ] 搭建 React + Vite 项目
- [ ] 对接 OpenClaw API
- [ ] 实现基础聊天界面
- [ ] 集成飞书登录

### Phase 2: 组织架构 (1 周)
- [ ] 创建 CEO/Manager/Worker Agent
- [ ] 实现 Agent 层级调用
- [ ] 飞书 Bitable 存储结构

### Phase 3: OKR 管理 (1 周)
- [ ] OKR CRUD 界面
- [ ] 进度可视化
- [ ] AI 辅助分析

### Phase 4: 流程编排 (2 周)
- [ ] React Flow 集成
- [ ] 节点配置面板
- [ ] 执行引擎对接

### Phase 5: 绩效评估 (1 周)
- [ ] 评估数据采集
- [ ] AI 分析生成
- [ ] 报告导出

## 六、技术栈清单

| 类别 | 技术 | 版本 |
|------|------|------|
| 前端框架 | React | 18.x |
| 流程编辑器 | @xyflow/react | 12.x |
| 构建工具 | Vite | 5.x |
| HTTP 客户端 | Axios | 1.6.x |
| 后端 API | Express | 4.x |
| 数据存储 | 飞书 Bitable | - |
| AI 引擎 | OpenClaw | latest |

## 七、注意事项

1. **Agent Prompt 需精心设计** - CEO/Manager/Worker 职责分离
2. **Bitable 权限** - 确保有读写权限
3. **API 限流** - 合理控制调用频率
4. **错误处理** - Agent 调用失败需有兜底
5. **日志追踪** - 关键操作记录日志

---

*文档版本: v1.0*  
*创建时间: 2026-03-02*
