package channel

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"agent-flow/internal/store"
)

// ChannelType 渠道类型
type ChannelType string

const (
	ChannelFeishu   ChannelType = "feishu"
	ChannelTelegram ChannelType = "telegram"
	ChannelDiscord  ChannelType = "discord"
	ChannelWhatsApp ChannelType = "whatsapp"
	ChannelWeChat   ChannelType = "wechat"
	ChannelSignal   ChannelType = "signal"
	ChannelWebAPI   ChannelType = "webapi"
)

// Message 统一消息格式
type Message struct {
	Type      string      `json:"type"`       // text/image/file
	Content   string      `json:"content"`    // 文本内容或文件URL
	UserID    string      `json:"user_id"`    // 用户ID
	ChannelID string      `json:"channel_id"` // 渠道ID
	Channel   string      `json:"channel"`    // 渠道类型
	RawData   interface{} `json:"raw_data"`   // 原始数据
}

// Sender 消息发送者接口
type Sender interface {
	SendMessage(msg Message) error
	SendText(userID, text string) error
}

// Manager 渠道管理器
type Manager struct {
	db     *store.Postgres
	redis  *store.Redis
	adapters map[ChannelType]Adapter
}

// NewManager 创建渠道管理器
func NewManager(db *store.Postgres, redis *store.Redis) *Manager {
	m := &Manager{
		db:     db,
		redis:  redis,
		adapters: make(map[ChannelType]Adapter),
	}

	// 注册渠道适配器
	m.adapters[ChannelFeishu] = NewFeishuAdapter()
	m.adapters[ChannelTelegram] = NewTelegramAdapter()
	m.adapters[ChannelDiscord] = NewDiscordAdapter()
	m.adapters[ChannelWhatsApp] = NewWhatsAppAdapter()
	m.adapters[ChannelWeChat] = NewWeChatAdapter()

	return m
}

// HandleWebhook 处理各渠道webhook
func (m *Manager) HandleWebhook(c *gin.Context) {
	channelType := c.Param("channel_type")

	adapter, ok := m.adapters[ChannelType(channelType)]
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported channel type"})
		return
	}

	// 解析消息
	msg, err := adapter.ParseWebhook(c.Request)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// 处理消息
	response, err := m.processMessage(msg)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// 发送回复
	if response != "" {
		if err := adapter.SendMessage(msg.UserID, response); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

// processMessage 处理消息（核心逻辑）
func (m *Manager) processMessage(msg *Message) (string, error) {
	// TODO: 
	// 1. 查找或创建会话
	// 2. 获取关联的Agent/Flow
	// 3. 调用AI处理
	// 4. 返回响应

	return "收到消息: " + msg.Content, nil
}

// Adapter 渠道适配器接口
type Adapter interface {
	ParseWebhook(req *http.Request) (*Message, error)
	SendMessage(userID, text string) error
	GetChannelType() ChannelType
}
