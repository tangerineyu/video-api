package ws

import (
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	ID     uint
	Socket *websocket.Conn
	Send   chan []byte
}
type Manager struct {
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Lock       sync.RWMutex
}

var WebManager = Manager{
	Clients:    make(map[uint]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (manager *Manager) Start() {
	for {
		select {
		case client := <-manager.Register:
			manager.Lock.Lock()
			manager.Clients[client.ID] = client
			manager.Lock.Unlock()
			log.Println("用户上线:", client.ID)

		case client := <-manager.Unregister:
			manager.Lock.Lock()
			if _, ok := manager.Clients[client.ID]; ok {
				close(client.Send)
				delete(manager.Clients, client.ID)
			}
			manager.Lock.Unlock()
			log.Println("用户断开连接:", client.ID)
		}
	}
}

// 发送信息`
func (manager *Manager) Send(userID uint, message []byte) {
	manager.Lock.Lock()
	defer manager.Lock.Unlock()
	if client, ok := manager.Clients[userID]; ok {
		select {
		case client.Send <- message:
		default:
			close(client.Send)
			delete(manager.Clients, userID)
		}
	}
}

// 允许跨域连接
// websocket.Upgrader将一个普通的Http请求，升级为websocket连接
var Upgrader = websocket.Upgrader{
	//CheckOrigin：安全配置，如果客户端地址与我的服务器的地址不一样，
	//默认的Upgrader会拒绝连接，如果没有这段代码，不同接口的前端连接后端时会被拒绝
	CheckOrigin: func(r *http.Request) bool {
		//返回true就是不管从哪里来都允许连接
		return true
	},
}
