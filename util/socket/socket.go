package socket

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

// upgrader 用于将 HTTP 连接升级为 WebSocket 连接
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024, // 读缓冲区大小
	WriteBufferSize: 1024, // 写缓冲区大小
	// 允许跨域（生产环境需根据实际情况限制）
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// Client 表示一个 WebSocket 客户端连接
type Client struct {
	conn *websocket.Conn // 底层 WebSocket 连接
	send chan []byte     // 消息发送通道
}

// 全局客户端管理器（用于管理所有连接的客户端）
type ClientManager struct {
	clients    map[*Client]bool // 所有在线客户端
	broadcast  chan []byte      // 广播消息通道
	register   chan *Client     // 客户端注册通道
	unregister chan *Client     // 客户端注销通道
}

// 初始化管理器
var manager = ClientManager{
	broadcast:  make(chan []byte),
	register:   make(chan *Client),
	unregister: make(chan *Client),
	clients:    make(map[*Client]bool),
}

// 启动管理器（处理客户端注册、注销和消息广播）
func (manager *ClientManager) start() {
	for {
		select {
		case client := <-manager.register:
			// 客户端注册
			manager.clients[client] = true
			log.Printf("新客户端连接，当前在线: %d", len(manager.clients))
		case client := <-manager.unregister:
			// 客户端注销
			if _, ok := manager.clients[client]; ok {
				delete(manager.clients, client)
				close(client.send)
				log.Printf("客户端断开连接，当前在线: %d", len(manager.clients))
			}
		case message := <-manager.broadcast:
			// 广播消息到所有客户端
			for client := range manager.clients {
				select {
				case client.send <- message:
					// 消息发送成功
				default:
					// 发送失败，移除客户端
					close(client.send)
					delete(manager.clients, client)
				}
			}
		}
	}
}

// 读取客户端消息（循环读取，直到连接关闭）
func (c *Client) readPump() {
	defer func() {
		manager.unregister <- c
		c.conn.Close()
	}()

	// 设置读取超时（可选，用于心跳检测）
	c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	c.conn.SetPongHandler(func(string) error {
		// 收到客户端的 pong 响应，更新超时时间
		c.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		// 读取客户端消息
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("读取消息错误: %v", err)
			}
			break
		}
		println(message)
		// 将消息发送到广播通道
		manager.broadcast <- message
	}
}

// 向客户端发送消息（循环发送，直到连接关闭）
func (c *Client) writePump() {
	// 定时发送心跳包（ping）
	ticker := time.NewTicker(30 * time.Second)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.send:
			// 设置写入超时
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if !ok {
				// 发送通道关闭，向客户端发送关闭帧
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			// 准备消息帧（文本类型）
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write(message)

			// 批量写入所有待发送消息（如果有）
			n := len(c.send)
			for i := 0; i < n; i++ {
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// 发送心跳 ping
			c.conn.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// WebSocket 处理函数（升级 HTTP 连接为 WebSocket）
func serveWs(w http.ResponseWriter, r *http.Request) {
	// 升级 HTTP 连接到 WebSocket
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	// 创建客户端实例
	client := &Client{
		conn: conn,
		send: make(chan []byte, 256), // 消息缓冲通道
	}

	// 注册客户端
	manager.register <- client

	// 启动读写协程（并发处理消息读写）
	go client.writePump()
	go client.readPump()
}

func Init() {
	// 启动客户端管理器
	go manager.start()

	// 注册 WebSocket 处理路由
	http.HandleFunc("/ws", serveWs)

	// 启动 HTTP 服务器（WebSocket 基于 HTTP 握手）
	log.Println("服务器启动，监听端口: 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
