package websocket

import (
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/aiqoder/monitor-lite-api/internal/pkg/log"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // 允许所有来源，您可能需要根据实际情况修改这个检查
	},
}

type Connection struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

type Hub struct {
	connections map[*Connection]bool
	Broadcast   chan []byte
	register    chan *Connection
	unregister  chan *Connection
}

func NewHub() *Hub {
	return &Hub{
		connections: make(map[*Connection]bool),
		Broadcast:   make(chan []byte),
		register:    make(chan *Connection),
		unregister:  make(chan *Connection),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case conn := <-h.register:
			h.connections[conn] = true
		case conn := <-h.unregister:
			if _, ok := h.connections[conn]; ok {
				delete(h.connections, conn)
				conn.conn.Close()
			}
		case message := <-h.Broadcast:
			for conn := range h.connections {
				conn.Write(message)
			}
		}
	}
}

func (c *Connection) Write(message []byte) {
	c.mu.Lock()
	defer c.mu.Unlock()
	err := c.conn.WriteMessage(websocket.TextMessage, message)
	if err != nil {
		log.Errorf("Error writing message: %v", err)
	}
}

func (h *Hub) HandleWebSocketGin(c *gin.Context) {
	h.HandleWebSocket(c.Writer, c.Request)
}

func (h *Hub) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Errorf("Error upgrading to WebSocket: %v", err)
		return
	}

	c := &Connection{conn: conn}
	h.register <- c

	defer func() {
		h.unregister <- c
	}()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("Error reading message: %v", err)
			}
			break
		}
		h.Broadcast <- message
	}
}
