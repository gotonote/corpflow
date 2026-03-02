import { useState, useCallback } from 'react'
import {
  ReactFlow,
  MiniMap,
  Controls,
  Background,
  useNodesState,
  useEdgesState,
  addEdge,
  Connection,
  Edge,
  Node,
  BackgroundVariant,
  Panel,
  Position,
  Handle,
} from '@xyflow/react'
import '@xyflow/react/dist/style.css'
import './FlowEditor.css'
import { flowAPI, executor } from './api'

// 节点数据类型
type NodeData = {
  label: string
  description?: string
  model?: string
  triggerType?: string
  condition?: string
  toolType?: string
  toolName?: string
  prompt?: string
  outputType?: string
  [key: string]: any
}

// 自定义节点组件
function AgentNode({ data, selected }: { data: NodeData; selected?: boolean }) {
  return (
    <div className={`custom-node agent-node ${selected ? 'selected' : ''}`}>
      <Handle type="target" position={Position.Top} />
      <div className="node-icon">🤖</div>
      <div className="node-content">
        <div className="node-label">{data.label}</div>
        <div className="node-desc">{data.description || 'AI智能体'}</div>
        <div className="node-model">{data.model || 'GPT-4'}</div>
      </div>
      <Handle type="source" position={Position.Bottom} />
    </div>
  )
}

function TriggerNode({ data, selected }: { data: NodeData['data']; selected?: boolean }) {
  return (
    <div className={`custom-node trigger-node ${selected ? 'selected' : ''}`}>
      <Handle type="target" position={Position.Top} />
      <div className="node-icon">⚡</div>
      <div className="node-content">
        <div className="node-label">{data.label}</div>
        <div className="node-desc">{data.triggerType || '消息触发'}</div>
      </div>
      <Handle type="source" position={Position.Bottom} />
    </div>
  )
}

function ConditionNode({ data, selected }: { data: NodeData; selected?: boolean }) {
  return (
    <div className={`custom-node condition-node ${selected ? 'selected' : ''}`}>
      <Handle type="target" position={Position.Top} />
      <Handle type="target" position={Position.Left} id="true" />
      <div className="node-icon">🔀</div>
      <div className="node-content">
        <div className="node-label">条件分支</div>
        <div className="node-desc">{data.condition || '条件判断'}</div>
      </div>
      <Handle type="source" position={Position.Bottom} id="true" />
      <Handle type="source" position={Position.Right} id="false" />
    </div>
  )
}

function ToolNode({ data, selected }: { data: NodeData; selected?: boolean }) {
  return (
    <div className={`custom-node tool-node ${selected ? 'selected' : ''}`}>
      <Handle type="target" position={Position.Top} />
      <div className="node-icon">🔧</div>
      <div className="node-content">
        <div className="node-label">{data.label}</div>
        <div className="node-desc">{data.toolName || data.toolType || '工具'}</div>
      </div>
      <Handle type="source" position={Position.Bottom} />
    </div>
  )
}

function LLMNode({ data, selected }: { data: NodeData; selected?: boolean }) {
  return (
    <div className={`custom-node llm-node ${selected ? 'selected' : ''}`}>
      <Handle type="target" position={Position.Top} />
      <div className="node-icon">🧠</div>
      <div className="node-content">
        <div className="node-label">大模型</div>
        <div className="node-desc">{data.model || 'GPT-4'}</div>
      </div>
      <Handle type="source" position={Position.Bottom} />
    </div>
  )
}

function OutputNode({ data, selected }: { data: NodeData; selected?: boolean }) {
  return (
    <div className={`custom-node output-node ${selected ? 'selected' : ''}`}>
      <Handle type="target" position={Position.Top} />
      <div className="node-icon">📤</div>
      <div className="node-content">
        <div className="node-label">输出</div>
        <div className="node-desc">{data.outputType || '返回结果'}</div>
      </div>
    </div>
  )
}

// 节点类型映射
const nodeTypes = {
  agent: AgentNode,
  trigger: TriggerNode,
  condition: ConditionNode,
  tool: ToolNode,
  llm: LLMNode,
  output: OutputNode,
}

// 初始节点
const initialNodes: Node[] = [
  {
    id: '1',
    type: 'trigger',
    position: { x: 250, y: 50 },
    data: { label: '消息触发', triggerType: '用户消息' },
  },
  {
    id: '2',
    type: 'agent',
    position: { x: 250, y: 200 },
    data: { label: '主Agent', description: '处理用户请求', model: 'GPT-4' },
  },
  {
    id: '3',
    type: 'output',
    position: { x: 250, y: 350 },
    data: { label: '返回结果', outputType: '文本' },
  },
]

const initialEdges: Edge[] = [
  { id: 'e1-2', source: '1', target: '2', animated: true },
  { id: 'e2-3', source: '2', target: '3', animated: true },
]

// 侧边栏组件
function Sidebar({ onDrag }: { onDrag: any }) {
  return (
    <aside className="sidebar">
      <h3>📦 节点库</h3>
      <div className="node-palette">
        <div className="palette-section">
          <div className="palette-title">触发器</div>
          <div className="palette-item" onClick={() => onDrag('trigger', '消息触发')}>
            <span>⚡</span> 消息触发
          </div>
          <div className="palette-item" onClick={() => onDrag('trigger', '定时任务')}>
            <span>⏰</span> 定时任务
          </div>
          <div className="palette-item" onClick={() => onDrag('trigger', 'Webhook')}>
            <span>🔗</span> Webhook
          </div>
        </div>
        
        <div className="palette-section">
          <div className="palette-title">智能体</div>
          <div className="palette-item" onClick={() => onDrag('agent', 'AI智能体')}>
            <span>🤖</span> AI智能体
          </div>
          <div className="palette-item" onClick={() => onDrag('llm', '大模型')}>
            <span>🧠</span> 大模型
          </div>
        </div>

        <div className="palette-section">
          <div className="palette-title">流程控制</div>
          <div className="palette-item" onClick={() => onDrag('condition', '条件分支')}>
            <span>🔀</span> 条件分支
          </div>
        </div>

        <div className="palette-section">
          <div className="palette-title">工具</div>
          <div className="palette-item" onClick={() => onDrag('tool', '浏览器')}>
            <span>🌐</span> 浏览器
          </div>
          <div className="palette-item" onClick={() => onDrag('tool', '搜索')}>
            <span>🔍</span> 网页搜索
          </div>
          <div className="palette-item" onClick={() => onDrag('tool', '计算器')}>
            <span>🧮</span> 计算器
          </div>
          <div className="palette-item" onClick={() => onDrag('tool', '代码执行')}>
            <span>💻</span> 代码执行
          </div>
        </div>

        <div className="palette-section">
          <div className="palette-title">输出</div>
          <div className="palette-item" onClick={() => onDrag('output', '返回结果')}>
            <span>📤</span> 返回结果
          </div>
        </div>
      </div>
    </aside>
  )
}

// 节点配置面板
function PropertiesPanel({ 
  selectedNode, 
  setNodes
}: { 
  selectedNode: any
  setNodes: any
}) {
  if (!selectedNode) {
    return (
      <aside className="properties-panel">
        <h3>⚙️ 节点配置</h3>
        <div className="no-selection">
          选择一个节点进行配置
        </div>
      </aside>
    )
  }

  const handleChange = (key: string, value: string) => {
    if (selectedNode && setNodes) {
      const node = selectedNode as Node<any>
      node.data[key] = value
      setNodes((nds: any) =>
        nds.map((n: any) => (n.id === selectedNode.id ? { ...n, data: node.data } : n))
      )
    }
  }

  return (
    <aside className="properties-panel">
      <h3>⚙️ 节点配置</h3>
      <div className="property-group">
        <label>节点名称</label>
        <input 
          type="text" 
          value={selectedNode.data.label} 
          onChange={(e) => handleChange('label', e.target.value)}
        />
      </div>

      {selectedNode.type === 'agent' && (
        <>
          <div className="property-group">
            <label>描述</label>
            <input 
              type="text" 
              value={selectedNode.data.description || ''} 
              onChange={(e) => handleChange('description', e.target.value)}
            />
          </div>
          <div className="property-group">
            <label>模型</label>
            <select 
              value={selectedNode.data.model || 'gpt-4'}
              onChange={(e) => handleChange('model', e.target.value)}
            >
              <option value="gpt-4">GPT-4</option>
              <option value="gpt-3.5-turbo">GPT-3.5 Turbo</option>
              <option value="claude-3-opus">Claude 3 Opus</option>
              <option value="claude-3-sonnet">Claude 3 Sonnet</option>
              <option value="glm-4">GLM-4</option>
            </select>
          </div>
        </>
      )}

      {selectedNode.type === 'trigger' && (
        <div className="property-group">
          <label>触发类型</label>
          <select 
            value={selectedNode.data.triggerType || '用户消息'}
            onChange={(e) => handleChange('triggerType', e.target.value)}
          >
            <option value="用户消息">用户消息</option>
            <option value="定时任务">定时任务</option>
            <option value="Webhook">Webhook</option>
          </select>
        </div>
      )}

      {selectedNode.type === 'tool' && (
        <>
          <div className="property-group">
            <label>工具类型</label>
            <select 
              value={selectedNode.data.toolType || 'browser'}
              onChange={(e) => handleChange('toolType', e.target.value)}
            >
              <option value="browser">浏览器</option>
              <option value="search">网页搜索</option>
              <option value="fetch">获取网页</option>
              <option value="calculator">计算器</option>
              <option value="code">代码执行</option>
            </select>
          </div>
          <div className="property-group">
            <label>工具名称</label>
            <input 
              type="text" 
              value={selectedNode.data.toolName || ''} 
              onChange={(e) => handleChange('toolName', e.target.value)}
            />
          </div>
        </>
      )}

      {selectedNode.type === 'condition' && (
        <div className="property-group">
          <label>条件表达式</label>
          <input 
            type="text" 
            value={selectedNode.data.condition || ''} 
            onChange={(e) => handleChange('condition', e.target.value)}
            placeholder="例如: input contains 'hello'"
          />
        </div>
      )}

      {selectedNode.type === 'llm' && (
        <div className="property-group">
          <label>System Prompt</label>
          <textarea 
            value={selectedNode.data.prompt || ''} 
            onChange={(e) => handleChange('prompt', e.target.value)}
            rows={4}
            placeholder="设置AI的系统提示词..."
          />
        </div>
      )}
    </aside>
  )
}

// 主流程编辑器组件
export default function FlowEditor() {
  const [nodes, setNodes, onNodesChange] = useNodesState(initialNodes)
  const [edges, setEdges, onEdgesChange] = useEdgesState(initialEdges)
  const [selectedNode, setSelectedNode] = useState<Node | null>(null)

  const onConnect = useCallback(
    (params: Connection) => setEdges((eds) => addEdge({ 
      ...params, 
      animated: true,
      style: { stroke: '#667eea', strokeWidth: 2 }
    }, eds)),
    [setEdges],
  )

  const onDragOver = useCallback((event: React.DragEvent) => {
    event.preventDefault()
    event.dataTransfer.dropEffect = 'move'
  }, [])

  const onDrop = useCallback(
    (event: React.DragEvent) => {
      event.preventDefault()
      
      const type = event.dataTransfer.getData('application/reactflow')
      const label = event.dataTransfer.getData('application/label')

      if (!type) return

      const reactFlowBounds = event.currentTarget.getBoundingClientRect()
      const position = {
        x: event.clientX - reactFlowBounds.left - 300,
        y: event.clientY - reactFlowBounds.top - 50,
      }

      const newNode: Node = {
        id: `node_${Date.now()}`,
        type,
        position,
        data: { label },
      }

      setNodes((nds) => nds.concat(newNode))
    },
    [setNodes],
  )

  const onNodeClick = useCallback((_: React.MouseEvent, node: Node) => {
    setSelectedNode(node)
  }, [])

  const onPaneClick = useCallback(() => {
    setSelectedNode(null)
  }, [])

  const handleDragStart = (event: React.DragEvent, type: string, label: string) => {
    event.dataTransfer.setData('application/reactflow', type)
    event.dataTransfer.setData('application/label', label)
    event.dataTransfer.effectAllowed = 'move'
  }

  // 保存流程
  const handleSave = async () => {
    const flowData = {
      name: '新流程',
      nodes: nodes.map(n => ({
        id: n.id,
        type: n.type,
        position: n.position,
        data: n.data,
      })),
      edges: edges.map(e => ({
        id: e.id,
        source: e.source,
        target: e.target,
      })),
    }
    try {
      await flowAPI.save(flowData)
      alert('流程已保存!')
    } catch (err) {
      console.error(err)
      alert('保存失败: ' + err)
    }
  }

  // 测试运行
  const handleRun = async () => {
    try {
      alert('流程执行中...')
      const result = await executor.run(nodes, { text: '测试输入' })
      alert('执行完成: ' + JSON.stringify(result))
    } catch (err) {
      console.error(err)
      alert('执行失败: ' + err)
    }
  }

  return (
    <div className="flow-editor">
      <div className="toolbar">
        <button className="btn-primary" onClick={handleSave}>💾 保存</button>
        <button className="btn-success" onClick={handleRun}>▶️ 执行</button>
        <button className="btn-secondary">📥 导入</button>
        <button className="btn-secondary">📤 导出</button>
      </div>
      
      <div className="editor-container">
        <Sidebar onDrag={handleDragStart} />
        
        <div className="reactflow-wrapper" onDrop={onDrop} onDragOver={onDragOver}>
          <ReactFlow
            nodes={nodes}
            edges={edges}
            onNodesChange={onNodesChange}
            onEdgesChange={onEdgesChange}
            onConnect={onConnect}
            onNodeClick={onNodeClick}
            onPaneClick={onPaneClick}
            nodeTypes={nodeTypes}
            fitView
            snapToGrid
            snapGrid={[15, 15]}
          >
            <Controls />
            <MiniMap 
              nodeColor={(node) => {
                switch (node.type) {
                  case 'agent': return '#667eea'
                  case 'trigger': return '#f59e0b'
                  case 'condition': return '#10b981'
                  case 'tool': return '#ef4444'
                  default: return '#999'
                }
              }}
            />
            <Background variant={BackgroundVariant.Dots} gap={12} size={1} />
            <Panel position="top-right">
              <div className="flow-info">
                节点: {nodes.length} | 连线: {edges.length}
              </div>
            </Panel>
          </ReactFlow>
        </div>

        <PropertiesPanel 
          selectedNode={selectedNode} 
          setNodes={setNodes}
        />
      </div>
    </div>
  )
}
