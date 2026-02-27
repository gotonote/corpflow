import { useState, useEffect, useRef } from 'react'
import './Chat.css'

interface Message {
  id: string
  type: string
  content: string
  sender: 'user' | 'bot'
  sender_id: string
  created_at: string
}

interface Conversation {
  id: string
  title: string
  last_message: string
  updated_at: string
  messages: Message[]
}

interface ChatProps {
  userId: string
  conversationId?: string
  onConversationChange?: (id: string) => void
}

export default function Chat({ userId, conversationId, onConversationChange }: ChatProps) {
  const [conversations, setConversations] = useState<Conversation[]>([])
  const [currentConv, setCurrentConv] = useState<Conversation | null>(null)
  const [messages, setMessages] = useState<Message[]>([])
  const [input, setInput] = useState('')
  const [loading, setLoading] = useState(false)
  const messagesEndRef = useRef<HTMLDivElement>(null)
  const wsRef = useRef<WebSocket | null>(null)

  // 加载会话列表
  useEffect(() => {
    loadConversations()
  }, [userId])

  // 加载指定会话
  useEffect(() => {
    if (conversationId) {
      loadConversation(conversationId)
    }
  }, [conversationId])

  // 滚动到底部
  useEffect(() => {
    messagesEndRef.current?.scrollIntoView({ behavior: 'smooth' })
  }, [messages])

  // WebSocket连接
  useEffect(() => {
    if (currentConv?.id) {
      connectWebSocket(currentConv.id)
    }
    return () => {
      wsRef.current?.close()
    }
  }, [currentConv?.id])

  const loadConversations = async () => {
    try {
      const res = await fetch(`/api/chat/conversations?user_id=${userId}`)
      const data = await res.json()
      setConversations(data)
    } catch (err) {
      console.error('Failed to load conversations:', err)
    }
  }

  const loadConversation = async (id: string) => {
    try {
      const res = await fetch(`/api/chat/conversations/${id}`)
      const data = await res.json()
      setCurrentConv(data)
      setMessages(data.messages || [])
    } catch (err) {
      console.error('Failed to load conversation:', err)
    }
  }

  const createConversation = async () => {
    try {
      const res = await fetch('/api/chat/conversations', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ user_id: userId, agent_id: 'default' })
      })
      const conv = await res.json()
      setConversations([conv, ...conversations])
      setCurrentConv(conv)
      setMessages([])
      onConversationChange?.(conv.id)
    } catch (err) {
      console.error('Failed to create conversation:', err)
    }
  }

  const sendMessage = async () => {
    if (!input.trim() || !currentConv?.id) return

    const userMessage: Message = {
      id: `temp-${Date.now()}`,
      type: 'text',
      content: input,
      sender: 'user',
      sender_id: userId,
      created_at: new Date().toISOString()
    }

    setMessages(prev => [...prev, userMessage])
    setInput('')
    setLoading(true)

    try {
      const res = await fetch('/api/chat/messages', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          conversation_id: currentConv.id,
          type: 'text',
          content: input,
          sender: 'user',
          sender_id: userId
        })
      })
      const botMessage = await res.json()
      setMessages(prev => [...prev.filter(m => m.id !== userMessage.id), botMessage])
    } catch (err) {
      console.error('Failed to send message:', err)
      setMessages(prev => prev.filter(m => m.id !== userMessage.id))
    } finally {
      setLoading(false)
    }
  }

  const connectWebSocket = (convId: string) => {
    const ws = new WebSocket(`ws://localhost:8080/ws?user_id=${userId}&conversation_id=${convId}`)
    
    ws.onmessage = (event) => {
      const message = JSON.parse(event.data)
      setMessages(prev => [...prev, message])
    }

    ws.onerror = (err) => {
      console.error('WebSocket error:', err)
    }

    wsRef.current = ws
  }

  const handleKeyPress = (e: React.KeyboardEvent) => {
    if (e.key === 'Enter' && !e.shiftKey) {
      e.preventDefault()
      sendMessage()
    }
  }

  const formatTime = (time: string) => {
    return new Date(time).toLocaleTimeString('zh-CN', { 
      hour: '2-digit', 
      minute: '2-digit' 
    })
  }

  return (
    <div className="chat-container">
      {/* 侧边栏 */}
      <aside className="chat-sidebar">
        <div className="sidebar-header">
          <h2>💬 对话</h2>
          <button className="new-chat-btn" onClick={createConversation}>
            + 新建
          </button>
        </div>
        <div className="conversation-list">
          {conversations.map(conv => (
            <div 
              key={conv.id}
              className={`conversation-item ${currentConv?.id === conv.id ? 'active' : ''}`}
              onClick={() => {
                setCurrentConv(conv)
                loadConversation(conv.id)
                onConversationChange?.(conv.id)
              }}
            >
              <div className="conv-title">{conv.title || '新对话'}</div>
              <div className="conv-time">{formatTime(conv.updated_at)}</div>
            </div>
          ))}
        </div>
      </aside>

      {/* 聊天区域 */}
      <main className="chat-main">
        {currentConv ? (
          <>
            <header className="chat-header">
              <h3>{currentConv.title || '新对话'}</h3>
            </header>
            
            <div className="messages-container">
              {messages.length === 0 && (
                <div className="empty-messages">
                  <p>👋 开始对话吧！</p>
                </div>
              )}
              {messages.map(msg => (
                <div 
                  key={msg.id} 
                  className={`message ${msg.sender === 'user' ? 'user' : 'bot'}`}
                >
                  <div className="message-avatar">
                    {msg.sender === 'user' ? '👤' : '🤖'}
                  </div>
                  <div className="message-content">
                    <div className="message-text">{msg.content}</div>
                    <div className="message-time">{formatTime(msg.created_at)}</div>
                  </div>
                </div>
              ))}
              {loading && (
                <div className="message bot">
                  <div className="message-avatar">🤖</div>
                  <div className="message-content">
                    <div className="message-text typing">正在思考...</div>
                  </div>
                </div>
              )}
              <div ref={messagesEndRef} />
            </div>

            <div className="input-container">
              <textarea
                value={input}
                onChange={(e) => setInput(e.target.value)}
                onKeyPress={handleKeyPress}
                placeholder="输入消息... (Enter 发送, Shift+Enter 换行)"
                rows={1}
              />
              <button onClick={sendMessage} disabled={!input.trim() || loading}>
                发送
              </button>
            </div>
          </>
        ) : (
          <div className="no-conversation">
            <p>选择一个对话或创建新对话</p>
            <button onClick={createConversation}>开始新对话</button>
          </div>
        )}
      </main>
    </div>
  )
}
