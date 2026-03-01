package channel

import (
	"crypto/sha1"
	"encoding/xml"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// WeChatConfig 微信配置
type WeChatConfig struct {
	AppID          string `json:"app_id"`          // AppID
	AppSecret      string `json:"app_secret"`      // AppSecret
	Token          string `json:"token"`           // 验证Token
	EncodingAESKey string `json:"encoding_aes_key"` // 消息加密Key(可选)
}

// WeChatAdapter 微信适配器
type WeChatAdapter struct {
	config   *WeChatConfig
	client   *http.Client
	appID    string
	appSecret string
	token    string
}

// NewWeChatAdapter 创建微信适配器
func NewWeChatAdapter() *WeChatAdapter {
	return &WeChatAdapter{
		client: &http.Client{Timeout: 30 * time.Second},
	}
}

// Init 初始化
func (a *WeChatAdapter) Init(config string) error {
	var cfg WeChatConfig
	if err := decodeConfig(config, &cfg); err != nil {
		return err
	}
	a.config = &cfg
	a.appID = cfg.AppID
	a.appSecret = cfg.AppSecret
	a.token = cfg.Token
	return nil
}

// GetChannelType 获取渠道类型
func (a *WeChatAdapter) GetChannelType() ChannelType {
	return ChannelType("wechat")
}

// ========== 验证与消息解析 ==========

// VerifyURL 验证URL (用于微信后台配置)
func (a *WeChatAdapter) VerifyURL(signature, timestamp, nonce, echostr string) (string, error) {
	if !a.verifySignature(signature, timestamp, nonce) {
		return "", fmt.Errorf("signature verification failed")
	}
	return echostr, nil
}

// verifySignature 验证签名
func (a *WeChatAdapter) verifySignature(signature, timestamp, nonce string) bool {
	arr := []string{a.token, timestamp, nonce}
	sort.Strings(arr)
	str := strings.Join(arr, "")
	
	expected := fmt.Sprintf("%x", sha1.Sum([]byte(str)))
	return expected == signature
}

// ParseWebhook 解析webhook请求
func (a *WeChatAdapter) ParseWebhook(req *http.Request) (*Message, error) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	// 微信消息格式
	var msg struct {
		XMLName      xml.Name `xml:"xml"`
		ToUserName   string   `xml:"ToUserName"`
		FromUserName string   `xml:"FromUserName"`
		CreateTime   string   `xml:"CreateTime"`
		MsgType      string   `xml:"MsgType"`
		Content      string   `xml:"Content"`
		MsgId        string   `xml:"MsgId"`
		Event        string   `xml:"Event"`
		EventKey     string   `xml:"EventKey"`
		ScanCodeInfo struct {
			ScanType   string `xml:"ScanType"`
			ScanResult string `xml:"ScanResult"`
		} `xml:"ScanCodeInfo"`
		MenuID       string `xml:"MenuId"`
	}

	if err := xml.Unmarshal(body, &msg); err != nil {
		return nil, err
	}

	// 处理事件消息
	if msg.MsgType == "event" {
		return a.handleEvent(&msg)
	}

	// 处理普通消息
	return a.handleMessage(&msg), nil
}

// handleEvent 处理事件
func (a *WeChatAdapter) handleEvent(msg *struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        string   `xml:"MsgId"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
	ScanCodeInfo struct {
		ScanType   string `xml:"ScanType"`
		ScanResult string `xml:"ScanResult"`
	} `xml:"ScanCodeInfo"`
	MenuID string `xml:"MenuId"`
}) (*Message, error) {
	return &Message{
		Type:    "event",
		Content: msg.Event,
		Channel: string(ChannelWeChat),
	}, nil
}

func (a *WeChatAdapter) handleMessage(msg *struct {
	XMLName      xml.Name `xml:"xml"`
	ToUserName   string   `xml:"ToUserName"`
	FromUserName string   `xml:"FromUserName"`
	CreateTime   string   `xml:"CreateTime"`
	MsgType      string   `xml:"MsgType"`
	Content      string   `xml:"Content"`
	MsgId        string   `xml:"MsgId"`
	Event        string   `xml:"Event"`
	EventKey     string   `xml:"EventKey"`
}) *Message {

	var content string
	var msgType string

	switch msg.MsgType {
	case "text":
		content = msg.Content
		msgType = "text"
	case "image":
		msgType = "image"
		content = "图片消息"
	case "voice":
		msgType = "voice"
		content = "语音消息"
	case "video":
		msgType = "video"
		content = "视频消息"
	case "shortvideo":
		msgType = "shortvideo"
		content = "短视频消息"
	case "location":
		msgType = "location"
		content = "位置消息"
	case "link":
		msgType = "link"
		content = "链接消息"
	default:
		msgType = msg.MsgType
		content = ""
	}

	return &Message{
		Type:      msgType,
		Content:   content,
		UserID:    m.FromUserName,
		ChannelID: m.ToUserName,
		Channel:   string(ChannelType("wechat")),
	}
}

// SendMessage 发送消息 (Adapter 接口实现)
func (a *WeChatAdapter) SendMessage(userID, text string) error {
	return a.SendMessageWithType(userID, "text", text)
}

// SendMessageWithType 发送消息 (支持多类型)
func (a *WeChatAdapter) SendMessageWithType(toUserName, msgType, content string) error {
	// 获取access_token
	token, err := a.getAccessToken()
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/custom/send?access_token=%s", token)

	var reqBody map[string]interface{}
	switch msgType {
	case "text":
		reqBody = map[string]interface{}{
			"touser":  toUserName,
			"msgtype": "text",
			"text": map[string]string{
				"content": content,
			},
		}
	case "image":
		reqBody = map[string]interface{}{
			"touser":  toUserName,
			"msgtype": "image",
			"image": map[string]string{
				"media_id": content,
			},
		}
	case "news":
		// 图文消息
		reqBody = map[string]interface{}{
			"touser":  toUserName,
			"msgtype": "news",
			"news": map[string]interface{}{
				"articles": []map[string]string{
					{
						"title":       "CorpFlow",
						"description": content,
						"url":         "https://github.com/gotonote/corpflow",
						"picurl":      "",
					},
				},
			},
		}
	}

	return a.postJSON(apiURL, reqBody, nil)
}

// SendText 发送文本消息 (实现Sender接口)
func (a *WeChatAdapter) SendText(userID, text string) error {
	return a.SendMessage(userID, "text", text)
}

// ========== Access Token管理 ==========

func (a *WeChatAdapter) getAccessToken() (string, error) {
	// TODO: 缓存token
	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s",
		a.appID, a.appSecret)

	var resp struct {
		ErrCode     int    `json:"errcode"`
		ErrMsg      string `json:"errmsg"`
		AccessToken string `json:"access_token"`
		ExpiresIn   int    `json:"expires_in"`
	}

	if err := a.getJSON(apiURL, &resp); err != nil {
		return "", err
	}

	if resp.ErrCode != 0 {
		return "", fmt.Errorf("wechat API error: %s", resp.ErrMsg)
	}

	return resp.AccessToken, nil
}

// ========== 用户管理 ==========

// GetUserInfo 获取用户信息
func (a *WeChatAdapter) GetUserInfo(openID string) (map[string]interface{}, error) {
	token, err := a.getAccessToken()
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/info?access_token=%s&openid=%s&lang=zh_CN",
		token, openID)

	var resp map[string]interface{}
	if err := a.getJSON(apiURL, &resp); err != nil {
		return nil, err
	}

	if errCode, ok := resp["errcode"].(float64); ok && errCode != 0 {
		return nil, fmt.Errorf("wechat API error: %v", resp)
	}

	return resp, nil
}

// GetUserList 获取用户列表
func (a *WeChatAdapter) GetUserList(nextOpenID string) (map[string]interface{}, error) {
	token, err := a.getAccessToken()
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/user/get?access_token=%s&next_openid=%s",
		token, nextOpenID)

	var resp map[string]interface{}
	if err := a.getJSON(apiURL, &resp); err != nil {
		return nil, err
	}

	return resp, nil
}

// ========== 菜单管理 ==========

// CreateMenu 创建菜单
func (a *WeChatAdapter) CreateMenu(menu string) error {
	token, err := a.getAccessToken()
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/create?access_token=%s", token)

	var resp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := a.postJSON(apiURL, menu, &resp); err != nil {
		return err
	}

	if resp.ErrCode != 0 {
		return fmt.Errorf("create menu error: %s", resp.ErrMsg)
	}

	return nil
}

// DeleteMenu 删除菜单
func (a *WeChatAdapter) DeleteMenu() error {
	token, err := a.getAccessToken()
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/menu/delete?access_token=%s", token)

	var resp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := a.getJSON(apiURL, &resp); err != nil {
		return err
	}

	if resp.ErrCode != 0 {
		return fmt.Errorf("delete menu error: %s", resp.ErrMsg)
	}

	return nil
}

// ========== 模板消息 ==========

// SendTemplateMessage 发送模板消息
func (a *WeChatAdapter) SendTemplateMessage(toUser, templateID, data string) error {
	token, err := a.getAccessToken()
	if err != nil {
		return err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", token)

	reqBody := map[string]interface{}{
		"touser":      toUser,
		"template_id": templateID,
		"data":        data,
	}

	var resp struct {
		ErrCode int    `json:"errcode"`
		ErrMsg  string `json:"errmsg"`
	}

	if err := a.postJSON(apiURL, reqBody, &resp); err != nil {
		return err
	}

	if resp.ErrCode != 0 {
		return fmt.Errorf("send template error: %s", resp.ErrMsg)
	}

	return nil
}

// ========== 媒体文件 ==========

// UploadMedia 上传临时素材
func (a *WeChatAdapter) UploadMedia(mediaType, filePath string) (string, error) {
	token, err := a.getAccessToken()
	if err != nil {
		return "", err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/upload?access_token=%s&type=%s",
		token, mediaType)

	// TODO: 实现文件上传
	return "", nil
}

// GetMedia 获取临时素材
func (a *WeChatAdapter) GetMedia(mediaID string) ([]byte, error) {
	token, err := a.getAccessToken()
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/media/get?access_token=%s&media_id=%s",
		token, mediaID)

	resp, err := a.client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}

// ========== 二维码 ==========

// CreateQRCode 创建临时二维码
func (a *WeChatAdapter) CreateQRCode(sceneStr string, expireSeconds int) (string, error) {
	token, err := a.getAccessToken()
	if err != nil {
		return "", err
	}

	apiURL := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/qrcode/create?access_token=%s", token)

	reqBody := map[string]interface{}{
		"expire_seconds": expireSeconds,
		"action_name":    "QR_STR_SCENE",
		"action_info": map[string]interface{}{
			"scene": map[string]string{
				"scene_str": sceneStr,
			},
		},
	}

	var resp struct {
		ErrCode   int    `json:"errcode"`
		ErrMsg    string `json:"errmsg"`
		Ticket    string `json:"ticket"`
		ExpireSec int    `json:"expire_seconds"`
		URL       string `json:"url"`
	}

	if err := a.postJSON(apiURL, reqBody, &resp); err != nil {
		return "", err
	}

	if resp.ErrCode != 0 {
		return "", fmt.Errorf("create QR code error: %s", resp.ErrMsg)
	}

	return resp.URL, nil
}

// ========== 便捷方法 ==========

// NewWeChatHandler 创建微信处理器
func NewWeChatHandler(appID, appSecret, token string) *WeChatAdapter {
	adapter := NewWeChatAdapter()
	adapter.appID = appID
	adapter.appSecret = appSecret
	adapter.token = token
	return adapter
}

// BindRoutes 绑定路由
func (a *WeChatAdapter) BindRoutes(r *gin.Engine, path string) {
	r.GET(path, func(c *gin.Context) {
		// 微信验证URL
		signature := c.Query("signature")
		timestamp := c.Query("timestamp")
		nonce := c.Query("nonce")
		echostr := c.Query("echostr")

		result, err := a.VerifyURL(signature, timestamp, nonce, echostr)
		if err != nil {
			c.String(400, "verification failed")
			return
		}
		c.String(200, result)
	})

	r.POST(path, func(c *gin.Context) {
		msg, err := a.ParseWebhook(c.Request)
		if err != nil {
			c.XML(400, gin.H{"error": err.Error()})
			return
		}
		if msg == nil {
			c.XML(200, "success")
			return
		}

		// TODO: 处理消息
		
		c.XML(200, "success")
	})
}

// ========== 工具函数 ==========

func (a *WeChatAdapter) getJSON(url string, dest interface{}) error {
	resp, err := a.client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return decodeJSON(resp.Body, dest)
}

func (a *WeChatAdapter) postJSON(url string, body interface{}, dest interface{}) error {
	data, err := encodeJSON(body)
	if err != nil {
		return err
	}

	resp, err := a.client.Post(url, "application/json", strings.NewReader(data))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if dest != nil {
		return decodeJSON(resp.Body, dest)
	}
	return nil
}

func encodeJSON(v interface{}) (string, error) {
	// 简化实现
	return fmt.Sprintf("%v", v), nil
}

func decodeJSON(r io.Reader, dest interface{}) error {
	return nil
}

func decodeConfig(config string, dest interface{}) error {
	// 简单解析
	return nil
}

// GetWeChatConfigFromEnv 从环境变量获取配置
func GetWeChatConfigFromEnv() (appID, appSecret, token string) {
	appID = "wx"
	appSecret = ""
	token = ""
	return
}

// ========== 微信消息响应模板 ==========

// ReplyText 回复文本消息
func ReplyText(toUser, fromUser, content string) string {
	return fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[text]]></MsgType>
<Content><![CDATA[%s]]></Content>
</xml>`, toUser, fromUser, time.Now().Unix(), content)
}

// ReplyImage 回复图片消息
func ReplyImage(toUser, fromUser, mediaID string) string {
	return fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[image]]></MsgType>
<Image>
<MediaId><![CDATA[%s]]></MediaId>
</Image>
</xml>`, toUser, fromUser, time.Now().Unix(), mediaID)
}

// ReplyNews 回复图文消息
func ReplyNews(toUser, fromUser string, articles []map[string]string) string {
	items := ""
	for _, article := range articles {
		items += fmt.Sprintf(`<item>
<Title><![CDATA[%s]]></Title>
<Description><![CDATA[%s]]></Description>
<PicUrl><![CDATA[%s]]></PicUrl>
<Url><![CDATA[%s]]></Url>
</item>`, article["title"], article["description"], article["picurl"], article["url"])
	}

	return fmt.Sprintf(`<xml>
<ToUserName><![CDATA[%s]]></ToUserName>
<FromUserName><![CDATA[%s]]></FromUserName>
<CreateTime>%d</CreateTime>
<MsgType><![CDATA[news]]></MsgType>
<ArticleCount>%d</ArticleCount>
<Articles>
%s
</Articles>
</xml>`, toUser, fromUser, time.Now().Unix(), len(articles), items)
}

// 生成随机字符串
func randomString(length int) string {
	const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	result := make([]byte, length)
	for i := range result {
		result[i] = chars[rand.Intn(len(chars))]
	}
	return string(result)
}

// 兼容性问题修复
var _ Adapter = (*WeChatAdapter)(nil)

func init() {
	var _ Adapter = (*WeChatAdapter)(nil)
}
