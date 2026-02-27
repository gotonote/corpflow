import { useState } from 'react'
import './App.css'

type Tab = 'home' | 'flow' | 'chat' | 'agents' | 'settings'

function App() {
  const [activeTab, setActiveTab] = useState<Tab>('home')

  return (
    <div className="app">
      <header className="header">
        <h1>🚀 CorpFlow</h1>
        <nav>
          <button className={activeTab === 'home' ? 'active' : ''} onClick={() => setActiveTab('home')}>
            🏠 Home
          </button>
          <button className={activeTab === 'chat' ? 'active' : ''} onClick={() => setActiveTab('chat')}>
            💬 Chat
          </button>
          <button className={activeTab === 'flow' ? 'active' : ''} onClick={() => setActiveTab('flow')}>
            🔀 Flows
          </button>
          <button className={activeTab === 'agents' ? 'active' : ''} onClick={() => setActiveTab('agents')}>
            🤖 Agents
          </button>
          <button className={activeTab === 'settings' ? 'active' : ''} onClick={() => setActiveTab('settings')}>
            ⚙️ Settings
          </button>
        </nav>
      </header>
      
      <main className="main">
        {activeTab === 'home' && <HomePanel />}
        {activeTab === 'chat' && <ChatPanel />}
        {activeTab === 'flow' && <FlowPanel />}
        {activeTab === 'agents' && <AgentsPanel />}
        {activeTab === 'settings' && <SettingsPanel />}
      </main>
    </div>
  )
}

// Quick Start Templates
const templates = [
  { id: 'simple-chat', name: '💬 Simple Chat', desc: 'Basic AI conversation', icon: '💬' },
  { id: 'multi-voting', name: '🗳️ Multi-Model Vote', desc: 'Multiple AI vote on best answer', icon: '🗳️' },
  { id: 'research', name: '🔍 Research Assistant', desc: 'Search & analyze', icon: '🔍' },
  { id: 'customer-service', name: '🎧 Customer Service', desc: 'AI support bot', icon: '🎧' },
  { id: 'code-review', name: '📝 Code Review', desc: 'Automated code review', icon: '📝' },
  { id: 'content', name: '✍️ Content Creator', desc: 'Social media content', icon: '✍️' },
]

function HomePanel() {
  const [selectedTemplate, setSelectedTemplate] = useState<string | null>(null)

  return (
    <div className="home-container">
      {/* Welcome */}
      <section className="welcome-section">
        <h2>Welcome to CorpFlow</h2>
        <p>Multi-Agent Collaboration Platform | 多智能体协作平台</p>
        <div className="quick-actions">
          <button className="btn-primary" onClick={() => window.location.href = '/chat'}>
            🚀 Start Chatting
          </button>
          <button className="btn-secondary" onClick={() => window.location.href = '/flow'}>
            ➕ Create Flow
          </button>
        </div>
      </section>

      {/* Quick Templates */}
      <section className="templates-section">
        <h3>⚡ Quick Start Templates</h3>
        <p className="section-desc">Click to use, no configuration needed</p>
        
        <div className="templates-grid">
          {templates.map(t => (
            <div key={t.id} className="template-card" onClick={() => setSelectedTemplate(t.id)}>
              <div className="template-icon">{t.icon}</div>
              <div className="template-name">{t.name}</div>
              <div className="template-desc">{t.desc}</div>
            </div>
          ))}
        </div>
      </section>

      {/* Features */}
      <section className="features-section">
        <h3>✨ Features</h3>
        <div className="features-grid">
          <div className="feature-card">
            <span>🤖</span>
            <h4>AI Agents</h4>
            <p>Support GPT-4, Claude, GLM-4, Kimi, Qwen, DeepSeek</p>
          </div>
          <div className="feature-card">
            <span>🗳️</span>
            <h4>Multi-Model Voting</h4>
            <p>Multiple AI discuss and vote on best answer</p>
          </div>
          <div className="feature-card">
            <span>🔀</span>
            <h4>Visual Flow</h4>
            <p>Drag-and-drop workflow automation</p>
          </div>
          <div className="feature-card">
            <span>💬</span>
            <h4>Multi-Channel</h4>
            <p>Feishu, WeChat, Telegram, Discord</p>
          </div>
          <div className="feature-card">
            <span>🧠</span>
            <h4>Memory</h4>
            <p>Supervisors can view subordinate work history</p>
          </div>
          <div className="feature-card">
            <span>📱</span>
            <h4>Mobile App</h4>
            <p>iOS, Android, Windows, Mac supported</p>
          </div>
        </div>
      </section>

      {/* Demo */}
      <section className="demo-section">
        <h3>📖 Demo</h3>
        
        <div className="demo-card">
          <h4>💬 Chat Demo</h4>
          <div className="demo-content">
            <p><strong>You:</strong> What's CorpFlow?</p>
            <p><strong>CorpFlow:</strong> CorpFlow is a multi-agent collaboration platform...</p>
          </div>
        </div>

        <div className="demo-card">
          <h4>🗳️ Voting Demo</h4>
          <div className="demo-content">
            <p><strong>Q:</strong> How to improve user experience?</p>
            <p>GPT-4: Suggestion 1... (score: 85)</p>
            <p>GLM-4: Suggestion 2... (score: 92) ⭐</p>
            <p className="winner">Winner: GLM-4</p>
          </div>
        </div>
      </section>
    </div>
  )
}

function ChatPanel() {
  const [messages, setMessages] = useState<{role: string, content: string}[]>([])
  const [input, setInput] = useState('')

  const send = () => {
    if (!input.trim()) return
    setMessages([...messages, {role: 'user', content: input}])
    setInput('')
    // Simulate response
    setTimeout(() => {
      setMessages(prev => [...prev, {role: 'bot', content: 'Configure API key in Settings to start chatting!'}])
    }, 500)
  }

  return (
    <div className="chat-container">
      <div className="chat-messages">
        {messages.length === 0 ? (
          <div className="chat-empty">
            <p>💬 Start a conversation</p>
            <p className="tip">Configure API key in Settings first</p>
          </div>
        ) : (
          messages.map((m, i) => (
            <div key={i} className={`message ${m.role}`}>
              <span className="msg-role">{m.role === 'user' ? '👤' : '🤖'}</span>
              <span className="msg-content">{m.content}</span>
            </div>
          ))
        )}
      </div>
      <div className="chat-input">
        <input 
          value={input} 
          onChange={e => setInput(e.target.value)} 
          onKeyPress={e => e.key === 'Enter' && send()}
          placeholder="Type message..."
        />
        <button onClick={send}>Send</button>
      </div>
    </div>
  )
}

function FlowPanel() {
  return (
    <div className="flow-container">
      <h3>🔀 Flow Editor</h3>
      <p>Visual workflow automation</p>
      <div className="flow-placeholder">
        <p>📝 Flow editor coming soon</p>
        <p className="tip">Use templates to get started!</p>
      </div>
    </div>
  )
}

function AgentsPanel() {
  return (
    <div className="agents-container">
      <h3>🤖 AI Agents</h3>
      <p>Manage your AI agents</p>
      <div className="agents-list">
        <div className="agent-item">
          <span>🤖</span>
          <div>
            <h4>Default Assistant</h4>
            <p>Model: GLM-4</p>
          </div>
        </div>
      </div>
      <button className="btn-primary">+ Add Agent</button>
    </div>
  )
}

function SettingsPanel() {
  const [apiKeys, setApiKeys] = useState({
    openai: '',
    zhipu: '',
    kimi: '',
  })

  return (
    <div className="settings-container">
      <h3>⚙️ Settings</h3>
      
      <div className="settings-section">
        <h4>🔑 API Keys</h4>
        <p className="tip">Enter API key to enable model</p>
        
        <div className="api-key-input">
          <label>OpenAI (GPT-4)</label>
          <input type="password" placeholder="sk-..." />
        </div>
        
        <div className="api-key-input">
          <label>Zhipu (GLM-4)</label>
          <input type="password" placeholder="Enter API key" />
        </div>
        
        <div className="api-key-input">
          <label>Kimi</label>
          <input type="password" placeholder="Enter API key" />
        </div>
        
        <div className="api-key-input">
          <label>Qwen</label>
          <input type="password" placeholder="Enter API key" />
        </div>
        
        <div className="api-key-input">
          <label>DeepSeek</label>
          <input type="password" placeholder="Enter API key" />
        </div>
        
        <button className="btn-save">Save Keys</button>
      </div>

      <div className="settings-section">
        <h4>🎯 Default Model</h4>
        <select>
          <option value="glm-4">GLM-4 (Recommended)</option>
          <option value="gpt-4">GPT-4</option>
          <option value="kimi">Kimi</option>
          <option value="qwen-turbo">Qwen Turbo</option>
          <option value="deepseek-chat">DeepSeek Chat</option>
        </select>
      </div>

      <div className="settings-section">
        <h4>🗳️ Multi-Model Voting</h4>
        <label className="toggle">
          <input type="checkbox" defaultChecked />
          <span>Enable voting (use multiple models)</span>
        </label>
      </div>
    </div>
  )
}

export default App
