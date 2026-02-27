package chat

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler HTTP处理器
type Handler struct {
	service *Service
}

// NewHandler 创建处理器
func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// ========== 会话API ==========

// CreateConversation 创建会话
func (h *Handler) CreateConversation(c *gin.Context) {
	var req struct {
		UserID   string `json:"user_id" binding:"required"`
		AgentID  string `json:"agent_id"`
		Channel  string `json:"channel"`
		ChannelID string `json:"channel_id"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.service.CreateConversation(
		req.UserID,
		req.AgentID,
		req.Channel,
		req.ChannelID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, conv)
}

// GetConversation 获取会话
func (h *Handler) GetConversation(c *gin.Context) {
	convID := c.Param("id")

	conv, err := h.service.GetConversation(convID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	c.JSON(http.StatusOK, conv)
}

// ListConversations 获取会话列表
func (h *Handler) ListConversations(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	convs, err := h.service.ListConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, convs)
}

// DeleteConversation 删除会话
func (h *Handler) DeleteConversation(c *gin.Context) {
	convID := c.Param("id")

	if err := h.service.DeleteConversation(convID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "deleted"})
}

// UpdateConversationTitle 更新会话标题
func (h *Handler) UpdateConversationTitle(c *gin.Context) {
	convID := c.Param("id")

	var req struct {
		Title string `json:"title" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	conv, err := h.service.GetConversation(convID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "conversation not found"})
		return
	}

	conv.Title = req.Title

	c.JSON(http.StatusOK, conv)
}

// ========== 消息API ==========

// SendMessage 发送消息
func (h *Handler) SendMessage(c *gin.Context) {
	var req struct {
		ConversationID string                 `json:"conversation_id" binding:"required"`
		Type           string                 `json:"type"`
		Content        string                 `json:"content" binding:"required"`
		Sender         string                 `json:"sender"`
		SenderID       string                 `json:"sender_id"`
		Metadata       map[string]interface{} `json:"metadata"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	msgType := req.Type
	if msgType == "" {
		msgType = "text"
	}
	sender := req.Sender
	if sender == "" {
		sender = "user"
	}

	msg, err := h.service.SendMessage(
		req.ConversationID,
		msgType,
		req.Content,
		sender,
		req.SenderID,
		req.Metadata,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, msg)
}

// GetMessages 获取消息列表
func (h *Handler) GetMessages(c *gin.Context) {
	convID := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.service.GetMessages(convID, limit, offset)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, messages)
}

// ========== WebSocket ==========

// WebSocketHandler WebSocket处理
func (h *Handler) WebSocket(c *gin.Context) {
	userID := c.Query("user_id")
	if userID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "user_id is required"})
		return
	}

	// 升级为WebSocket
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	client := &Client{
		ID:     generateID(),
		UserID: userID,
		Send:   make(chan []byte, 256),
		Hub:    h.service.hub,
	}

	// 注册客户端
	h.service.hub.register <- client

	// 启动读写协程
	go client.WritePump()
	go client.ReadPump(c.Param("conversation_id"))
}

// ReadPump 读取客户端消息
func (c *Client) ReadPump(convID string) {
	defer func() {
		c.Hub.unregister <- c
		c.Hub.conn.Close()
	}()

	for {
		_, message, err := c.Hub.conn.ReadMessage()
		if err != nil {
			break
		}

		// 解析消息
		var msg struct {
			Type      string `json:"type"`
			Content   string `json:"content"`
			ChannelID string `json:"channel_id"`
		}

		if err := json.Unmarshal(message, &msg); err != nil {
			continue
		}

		// 发送消息到服务
		if msg.Type == "message" {
			// 这里调用服务处理消息
			// h.service.ProcessUserMessage(convID, msg.Content)
		}
	}
}

// WritePump 向客户端写消息
func (c *Client) WritePump() {
	ticker := time.NewTicker(54 * time.Second)
	defer func() {
		ticker.Stop()
		c.Hub.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			if !ok {
				c.Hub.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Hub.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 添加排队消息
			n := len(c.Send)
			for i := 0; i < n; i++ {
				w.Write([]byte{'\n'})
				w.Write(<-c.Send)
			}

			if err := w.Close(); err != nil {
				return
			}

		case <-ticker.C:
			if err := c.Hub.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ========== 导入需要的包 ==========

// upgrader WebSocket升级器
import (
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize:  1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源
	},
}
