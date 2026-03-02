import express from 'express'
import cors from 'cors'
import axios from 'axios'

const app = express()
app.use(cors())
app.use(express.json())

const PORT = process.env.PORT || 8081

// 飞书配置
const FEISHU_APP_ID = process.env.FEISHU_APP_ID || 'cli_a9f5ba317c391cd0'
const FEISHU_APP_SECRET = process.env.FEISHU_APP_SECRET || 'i3esoDXjBYZG3hamKXsaFd6y3oruICAh'

// 获取飞书 access_token
let tenantAccessToken = ''
let tokenExpire = 0

async function getTenantAccessToken() {
  if (tenantAccessToken && Date.now() < tokenExpire) {
    return tenantAccessToken
  }
  const res = await axios.post('https://open.feishu.cn/open-apis/auth/v3/tenant_access_token/internal', {
    app_id: FEISHU_APP_ID,
    app_secret: FEISHU_APP_SECRET
  })
  tenantAccessToken = res.data.tenant_access_token
  tokenExpire = Date.now() + res.data.expire * 1000 - 60000
  return tenantAccessToken
}

// 流程存储（内存中，生产环境用数据库）
let flows = []

// API: 获取流程列表
app.get('/api/flows', async (req, res) => {
  res.json(flows)
})

// API: 保存流程
app.post('/api/flows', async (req, res) => {
  const flow = req.body
  const existingIndex = flows.findIndex(f => f.id === flow.id)
  if (existingIndex >= 0) {
    flows[existingIndex] = { ...flow, updatedAt: new Date().toISOString() }
  } else {
    flow.id = `flow_${Date.now()}`
    flow.createdAt = new Date().toISOString()
    flows.push(flow)
  }
  res.json({ success: true, flow })
})

// API: 执行流程
app.post('/api/flows/:id/execute', async (req, res) => {
  const { id } = req.params
  const flow = flows.find(f => f.id === id)
  if (!flow) {
    return res.status(404).json({ error: '流程不存在' })
  }
  
  // 简单执行：返回流程信息
  // 生产环境需要调用 OpenClaw 执行
  res.json({
    success: true,
    message: '流程执行成功',
    flow: flow.name
  })
})

// API: OKR 列表（模拟）
app.get('/api/okr', (req, res) => {
  res.json([
    { id: '1', objective: '提升产品体验', progress: 60, status: 'active' },
    { id: '2', objective: '扩大用户规模', progress: 40, status: 'active' }
  ])
})

// API: 聊天（调用 OpenClaw）
app.post('/api/chat/messages', async (req, res) => {
  const { message } = req.body
  // 这里可以调用 OpenClaw Gateway
  res.json({
    content: `收到消息: ${message}`,
    sender: 'bot'
  })
})

app.listen(PORT, () => {
  console.log(`CorpFlow API running on port ${PORT}`)
})
