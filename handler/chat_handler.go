package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"video-api/pkg/errno"
	"video-api/pkg/log"
	"video-api/pkg/utils"
	"video-api/pkg/ws"
	"video-api/service"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type ChatHandler struct {
	msgService service.IMessageService
}

func NewChatHandler(svc service.IMessageService) *ChatHandler {
	return &ChatHandler{msgService: svc}
}

// 消息发送结构体
type SendMsgRequest struct {
	ToUserID uint   `json:"to_user_id"`
	Content  string `json:"content"`
}

// 消息接受结构体
type ReplyMsg struct {
	FromUserID uint   `json:"from_user_id"`
	Content    string `json:"content"`
}

// 建立WebSocket连接
func (h *ChatHandler) Connect(c *gin.Context) {
	//websocket无法自定义Header，不能在bearer里写token，token直接写在URL里
	token := c.Query("token")
	claims, err := utils.ParseToken(token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "鉴权失败"})
		return
	}
	userID := claims.UserID
	//调用upgrade，将http协议升级为Websocket
	conn, err := ws.Upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}
	client := &ws.Client{
		ID:     userID,
		Socket: conn,
		Send:   make(chan []byte),
	}
	//注册到管理器
	ws.WebManager.Register <- client
	go h.WriteLoop(client)
	go h.ReadLoop(client)
}
func (h *ChatHandler) ReadLoop(client *ws.Client) {
	//readloop是一个死循环，退出的方式一般为用户直接关闭或网络错误
	//defer保证一定能注销和关闭
	defer func() {
		ws.WebManager.Unregister <- client
		client.Socket.Close()
	}()
	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			break
		}
		//解析JSON，发过来的数据不一定合法，做一个检查
		var req SendMsgRequest
		if err := json.Unmarshal(message, &req); err != nil {
			//不能因为用户一条格式错误的信息就断开连接，选择忽略此信息
			continue
		}
		//防止消息丢失，比如服务器重启或发送过程出错。存在数据库中
		h.msgService.SaveMessage(client.ID, req.ToUserID, req.Content)
		//组装回复，前端发送的是{toUserID:?,content:?},但是目标用户要得是{fromUserID:?,content:?}
		reply := ReplyMsg{
			FromUserID: client.ID,
			Content:    req.Content,
		}
		replyBytes, _ := json.Marshal(reply)
		//ReadLoop只有当前用户的连接，没有目标用户的连接，
		//webManager有全局的Clients名单，委托Manager
		ws.WebManager.Send(req.ToUserID, replyBytes)
	}
}
func (h *ChatHandler) WriteLoop(client *ws.Client) {
	defer func() {
		client.Socket.Close()
	}()
	for {
		select {
		case message, ok := <-client.Send:
			if !ok {
				return
			}
			client.Socket.WriteMessage(websocket.TextMessage, message)
		}
	}
}

// 获取聊天记录
func (h *ChatHandler) GetHistory(client *gin.Context) {
	userID, ok := getUserID(client)
	if !ok {
		return
	}

	toUserIDStr := client.Query("to_user_id")
	toUserID, _ := strconv.ParseUint(toUserIDStr, 10, 64)

	msgs, err := h.msgService.GetChatHistory(userID, uint(toUserID))
	if err != nil {
		log.Log.Error("获取聊天记录失败")
		Error(client, errno.ServiceErr)
		return
	}
	Success(client, msgs)
}
