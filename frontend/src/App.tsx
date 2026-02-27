import { useState } from 'react'
import FlowEditor from './FlowEditor'
import Chat from './Chat'
import './App.css'

type Tab = 'home' | 'flow' | 'chat' | 'agents' | 'channels' | 'settings'

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('home')
  const [currentUser] = useState('user-001')

  return (
    <div className="app">
      <header className="header">
        <h1>🚀 CorpFlow</h1>
        <nav>
          <button 
            className={activeTab === 'home' ? 'active' : ''} 
            onClick={() => setActiveTab('home')}
          >
            🏠 首页
          </button>
          <button 
            className={activeTab === 'chat' ? 'active' : ''} 
            onClick={() => setActiveTab('chat')}
          >
            💬 对话
          </button>
          <button 
            className={activeTab === 'flow' ? 'active' : ''} 
            onClick={() => setActiveTab('flow')}
          >
            流程编排
          </button>
          <button 
            className={activeTab === 'agents' ? 'active' : ''} 
            onClick={() => setActiveTab('agents')}
          >
            智能体
          </button>
          <button 
            className={activeTab === 'channels' ? 'active' : ''} 
            onClick={() => setActiveTab('channels')}
          >
            渠道
          </button>
          <button 
            className={activeTab === 'settings' ? 'active' : ''} 
            onClick={() => setActiveTab('settings')}
          >
            设置
          </button>
        </nav>
      </header>
      
      <main className="main">
        {activeTab === 'home' && <HomePanel />}
        {activeTab === 'chat' && <Chat userId={currentUser} />}
        {activeTab === 'flow' && <FlowEditor />}
        {activeTab === 'agents' && <AgentsPanel />}
        {activeTab === 'channels' && <ChannelsPanel />}
        {activeTab === 'settings' && <SettingsPanel />}
      </main>
    </div>
  )
}

function HomePanel() {
  return (
    <div className="home-panel">
      <div className="welcome-card">
        <h2>Welcome to CorpFlow</h2>
        <p>Multi-Agent Collaboration Platform</p>
        <p className="subtitle">多智能体协作平台</p>
      </div>

      <div className="demo-section">
        <h3>Demo / 示例</h3>
        
        <div className="demo-card">
          <h4>💬 Chat Demo</h4>
          <div className="demo-content">
            <p><strong>You:</strong> 什么是CorpFlow?</p>
            <p><strong>CorpFlow:</strong> CorpFlow是一个多智能体协作平台...</p>
          </div>
        </div>

        <div className="demo-card">
          <h4>🔀 Flow Demo</h4>
          <div className="demo-content">
            <p>流程: 触发器 → 智能体 → 条件分支 → 工具 → 输出</p>
            <p className="demo-desc">用户发送消息 → AI处理 → 判断是否需要工具 → 执行 → 返回结果</p>
          </div>
        </div>

        <div className="demo-card">
          <h4>🗳️ Multi-Model Voting Demo</h4>
          <div className="demo-content">
            <p><strong>问题:</strong> 如何提升产品用户体验?</p>
            <p>GPT-4: 建议1... (得分: 85)</p>
            <p>GLM-4: 建议2... (得分: 92) ⭐</p>
            <p>Kimi: 建议3... (得分: 78)</p>
            <p className="demo-winner">最终选择: GLM-4 (综合得分最高)</p>
          </div>
        </div>
      </div>

      <div className="features-section">
        <h3>Features / 功能</h3>
        <div className="feature-grid">
          <div className="feature-item">🤖 AI智能体</div>
          <div className="feature-item">🔀 流程编排</div>
          <div className="feature-item">💬 多渠道</div>
          <div className="feature-item">🗳️ 多模型投票</div>
          <div className="feature-item">📱 移动端</div>
          <div className="feature-item">🧠 记忆系统</div>
        </div>
      </div>
    </div>
  )
}

function AgentsPanel() {
  return (
    <div className="panel">
      <h2>智能体管理</h2>
      <p>创建和管理AI智能体</p>
    </div>
  )
}

function ChannelsPanel() {
  return (
    <div className="panel">
      <h2>渠道管理</h2>
      <p>配置消息接收渠道</p>
    </div>
  )
}

function SettingsPanel() {
  return (
    <div className="panel">
      <h2>系统设置</h2>
      <p>配置API密钥，大模型参数等</p>
    </div>
  )
}

export default App
