package chat

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"

	"corpflow/internal/store"
)

// MessageType 消息类型
type MessageType string

const (
	MessageTypeText    MessageType = "text"
	MessageTypeImage   MessageType = "image"
	MessageTypeFile    MessageType = "file"
	MessageTypeCommand MessageType = "command"
	MessageTypeSystem  MessageType = "system"
)

// Message 聊天消息
type Message struct {
	ID         string                 `json:"id"`
	Type       MessageType          `json:"type"`
	Content    string                 `json:"content"`
	Sender     string                 `json:"sender"`      // user/bot
	SenderID   string                 `json:"sender_id"`   // 用户ID
	ReceiverID string                 `json:"receiver_id"` // 接收者ID
	ChannelID  string                 `json:"channel_id"`  // 渠道ID
	Metadata   map[string]interface{} `json:"metadata"`    // 附加数据
	CreatedAt  time.Time             `json:"created_at"`
}

// Conversation 聊天会话
type Conversation struct {
	ID           string                 `json:"id"`
	UserID       string                 `json:"user_id"`
	AgentID      string                 `json:"agent_id"`
	Channel      string                 `json:"channel"`
	ChannelID    string                 `json:"channel_id"`
	Title        string                 `json:"title"`
	Messages     []Message              `json:"messages"`
	Context      map[string]interface{} `json:"context"`
	LastMessage  string                 `json:"last_message"`
	CreatedAt    time.Time             `json:"created_at"`
	UpdatedAt    time.Time             `json:"updated_at"`
}

// Service 聊天服务
type Service struct {
	db           *store.Postgres
	redis        *store.Redis
	hub          *Hub
	mu           sync.RWMutex
	conversations map[string]*Conversation // 内存缓存
}

// NewService 创建聊天服务
func NewService(db *store.Postgres, redis *store.Redis, hub *Hub) *Service {
	return &Service{
		db:           db,
		redis:        redis,
		hub:          hub,
		conversations: make(map[string]*Conversation),
	}
}

// ========== 会话管理 ==========

// CreateConversation 创建会话
func (s *Service) CreateConversation(userID, agentID, channel, channelID string) (*Conversation, error) {
	conv := &Conversation{
		ID:        generateID(),
		UserID:    userID,
		AgentID:   agentID,
		Channel:   channel,
		ChannelID: channelID,
		Title:     "新对话",
		Messages:  []Message{},
		Context:   make(map[string]interface{}),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 缓存到内存
	s.mu.Lock()
	s.conversations[conv.ID] = conv
	s.mu.Unlock()

	// 持久化到Redis
	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("conversation:%s", conv.ID)
		if data, err := json.Marshal(conv); err == nil {
			s.redis.Set(ctx, key, data, 24*time.Hour)
		}
	}

	return conv, nil
}

// GetConversation 获取会话
func (s *Service) GetConversation(conversationID string) (*Conversation, error) {
	s.mu.RLock()
	if conv, ok := s.conversations[conversationID]; ok {
		s.mu.RUnlock()
		return conv, nil
	}
	s.mu.RUnlock()

	// 从Redis加载
	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("conversation:%s", conversationID)
		var conv Conversation
		if err := s.redis.Get(ctx, key, &conv); err == nil {
			s.mu.Lock()
			s.conversations[conv.ID] = &conv
			s.mu.Unlock()
			return &conv, nil
		}
	}

	return nil, fmt.Errorf("conversation not found")
}

// ListConversations 获取用户会话列表
func (s *Service) ListConversations(userID string) ([]*Conversation, error) {
	s.mu.RLock()
	var convs []*Conversation
	for _, conv := range s.conversations {
		if conv.UserID == userID {
			convs = append(convs, conv)
		}
	}
	s.mu.RUnlock()

	// 按更新时间排序
	for i := 0; i < len(convs)-1; i++ {
		for j := i + 1; j < len(convs); j++ {
			if convs[i].UpdatedAt.Before(convs[j].UpdatedAt) {
				convs[i], convs[j] = convs[j], convs[i]
			}
		}
	}

	return convs, nil
}

// DeleteConversation 删除会话
func (s *Service) DeleteConversation(conversationID string) error {
	s.mu.Lock()
	delete(s.conversations, conversationID)
	s.mu.Unlock()

	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("conversation:%s", conversationID)
		s.redis.Del(ctx, key)
	}

	return nil
}

// ========== 消息管理 ==========

// SendMessage 发送消息
func (s *Service) SendMessage(convID, msgType, content, sender, senderID string, metadata map[string]interface{}) (*Message, error) {
	conv, err := s.GetConversation(convID)
	if err != nil {
		return nil, err
	}

	msg := &Message{
		ID:         generateID(),
		Type:       MessageType(msgType),
		Content:    content,
		Sender:     sender,
		SenderID:   senderID,
		ReceiverID: conv.UserID,
		ChannelID:  conv.ChannelID,
		Metadata:   metadata,
		CreatedAt:  time.Now(),
	}

	// 添加到会话
	conv.Messages = append(conv.Messages, *msg)
	conv.LastMessage = content
	conv.UpdatedAt = time.Now()

	// 更新标题（第一条消息）
	if len(conv.Messages) == 1 && sender == "user" {
		conv.Title = truncate(content, 50)
	}

	// 保存到Redis
	if s.redis != nil {
		ctx := context.Background()
		key := fmt.Sprintf("conversation:%s", conv.ID)
		if data, err := json.Marshal(conv); err == nil {
			s.redis.Set(ctx, key, data, 24*time.Hour)
		}
	}

	// 通过WebSocket推送
	if s.hub != nil {
		s.hub.SendToUser(conv.UserID, msg)
	}

	return msg, nil
}

// GetMessages 获取会话消息
func (s *Service) GetMessages(convID string, limit, offset int) ([]Message, error) {
	conv, err := s.GetConversation(convID)
	if err != nil {
		return nil, err
	}

	messages := conv.Messages
	if offset > len(messages) {
		offset = len(messages)
	}
	if limit <= 0 || offset+limit > len(messages) {
		limit = len(messages) - offset
	}

	return messages[offset : offset+limit], nil
}

// ========== 对话处理 ==========

// ProcessUserMessage 处理用户消息
func (s *Service) ProcessUserMessage(convID, content string) (*Message, error) {
	conv, err := s.GetConversation(convID)
	if err != nil {
		return nil, err
	}

	// 保存用户消息
	userMsg, err := s.SendMessage(convID, string(MessageTypeText), content, "user", conv.UserID, nil)
	if err != nil {
		return nil, err
	}

	// TODO: 调用AI处理
	// 这里模拟AI响应
	botResponse := "收到你的消息: " + content

	// 保存机器人消息
	botMsg, err := s.SendMessage(convID, string(MessageTypeText), botResponse, "bot", conv.AgentID, nil)
	if err != nil {
		return nil, err
	}

	return botMsg, nil
}

// ========== WebSocket Hub ==========

// Hub WebSocket中心
type Hub struct {
	// 注册用户
	register chan *Client
	// 注销用户
	unregister chan *Client
	// 广播消息
	broadcast chan *BroadcastMessage
	// 用户连接
	clients map[string]map[*Client]bool
	// 互斥锁
	mu sync.RWMutex
}

// BroadcastMessage 广播消息
type BroadcastMessage struct {
	Message   *Message
	TargetUser string // 空表示广播给所有用户
}

// Client WebSocket客户端
type Client struct {
	ID       string
	UserID   string
	Send     chan []byte
	Hub      *Hub
}

// NewHub 创建Hub
func NewHub() *Hub {
	return &Hub{
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *BroadcastMessage, 256),
		clients:    make(map[string]map[*Client]bool),
	}
}

// Run 运行Hub
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			if h.clients[client.UserID] == nil {
				h.clients[client.UserID] = make(map[*Client]bool)
			}
			h.clients[client.UserID][client] = true
			h.mu.Unlock()
			log.Printf("Client connected: %s (%s)", client.ID, client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if clients, ok := h.clients[client.UserID]; ok {
				if _, ok := clients[client]; ok {
					close(client.Send)
					delete(clients, client)
				}
				if len(clients) == 0 {
					delete(h.clients, client.UserID)
				}
			}
			h.mu.Unlock()
			log.Printf("Client disconnected: %s (%s)", client.ID, client.UserID)

		case msg := <-h.broadcast:
			h.mu.RLock()
			if msg.TargetUser != "" {
				// 发送给指定用户
				if clients, ok := h.clients[msg.TargetUser]; ok {
					data, _ := json.Marshal(msg.Message)
					for client := range clients {
						select {
						case client.Send <- data:
						default:
							close(client.Send)
							delete(clients, client)
						}
					}
				}
			} else {
				// 广播给所有用户
				data, _ := json.Marshal(msg.Message)
				for _, clients := range h.clients {
					for client := range clients {
						select {
						case client.Send <- data:
						default:
							close(client.Send)
							delete(clients, client)
						}
					}
				}
			}
			h.mu.RUnlock()
		}
	}
}

// SendToUser 发送给指定用户
func (h *Hub) SendToUser(userID string, msg *Message) {
	h.broadcast <- &BroadcastMessage{
		Message:   msg,
		TargetUser: userID,
	}
}

// Broadcast 广播消息
func (h *Hub) Broadcast(msg *Message) {
	h.broadcast <- &BroadcastMessage{
		Message: msg,
	}
}

// GetOnlineUsers 获取在线用户
func (h *Hub) GetOnlineUsers() []string {
	h.mu.RLock()
	defer h.mu.RUnlock()

	users := make([]string, 0, len(h.clients))
	for userID := range h.clients {
		users = append(users, userID)
	}
	return users
}

// ========== 工具函数 ==========

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), randomString(8))
}

func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[time.Now().UnixNano()%int64(len(chars))]
	}
	return string(result)
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen] + "..."
}
