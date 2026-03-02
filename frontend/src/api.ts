// CorpFlow API - 极简实现
// 使用后端 API 存储流程数据

const API_BASE = 'http://localhost:8081'

// 流程管理
export const flowAPI = {
  // 获取流程列表
  list: async () => {
    const res = await fetch(`${API_BASE}/flows`)
    return res.json()
  },

  // 保存流程
  save: async (flow: any) => {
    const res = await fetch(`${API_BASE}/flows`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(flow)
    })
    return res.json()
  },

  // 执行流程
  execute: async (flowId: string, input: any) => {
    const res = await fetch(`${API_BASE}/flows/${flowId}/execute`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(input)
    })
    return res.json()
  }
}

// OKR 管理
export const okrAPI = {
  list: async () => {
    const res = await fetch(`${API_BASE}/okr`)
    return res.json()
  },
  
  create: async (okr: any) => {
    const res = await fetch(`${API_BASE}/okr`, {
      method: 'POST',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(okr)
    })
    return res.json()
  },
  
  update: async (id: string, data: any) => {
    const res = await fetch(`${API_BASE}/okr/${id}`, {
      method: 'PUT',
      headers: { 'Content-Type': 'application/json' },
      body: JSON.stringify(data)
    })
    return res.json()
  }
}

// 执行器 - 调用 OpenClaw
export const executor = {
  run: async (nodes: any[], input: any) => {
    // 找到起始节点
    const startNode = nodes.find(n => n.type === 'trigger')
    if (!startNode) throw new Error('没有触发器')
    
    // 简单的线性执行
    let currentData = input
    for (const node of nodes) {
      if (node.type === 'llm' || node.type === 'agent') {
        // 调用 OpenClaw
        const res = await fetch('/api/chat/messages', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            message: currentData,
            prompt: node.data.prompt,
            model: node.data.model
          })
        })
        currentData = await res.json()
      } else if (node.type === 'tool') {
        // 执行工具
        // TODO: 调用 OpenClaw 工具
      }
    }
    return currentData
  }
}
