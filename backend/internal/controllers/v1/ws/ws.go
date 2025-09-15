package ws

import (
	"backend/internal/services"
	"backend/utils"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub      *Hub
	conn     *websocket.Conn
	send     chan *Message
	room     string
	userId   string
	username string
	userRole string
}

type Message struct {
	Content   string `json:"content"`
	Room      string `json:"room"`
	Username  string `json:"username"`
	UserRole  string `json:"userRole"`
	CreatedAt string `json:"createdAt"`
}

func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}

		var msg Message
		if err := json.Unmarshal(message, &msg); err != nil {
			log.Printf("error unmarshalling message: %v", err)
			continue
		}
		msg.Room = c.room
		msg.Username = c.username
		msg.UserRole = c.userRole
		msg.CreatedAt = time.Now().Format(time.RFC3339)

		// Save comment to DB
		ticketService := services.NewTicketService()
		_, err = ticketService.CreateComment(context.Background(), msg.Room, msg.Content, c.userId, msg.UserRole)
		if err != nil {
			log.Printf("failed to save comment: %v", err)
		}

		c.hub.broadcast <- &msg
	}
}

func (c *Client) writePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			json.NewEncoder(w).Encode(message)

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func ServeWs(hub *Hub, c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade failed:", err)
		return
	}

	// 1. Try to read from query param
	tokenString := c.Query("token")

	// 2. Fallback: check Authorization header
	if tokenString == "" {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			log.Println("Authorization header not found")
			conn.Close()
			return
		}

		tokenParts := strings.Split(authHeader, " ")
		if len(tokenParts) != 2 || tokenParts[0] != "Bearer" {
			log.Println("Invalid Authorization header format")
			conn.Close()
			return
		}
		tokenString = tokenParts[1]
	}

	// 3. Validate token
	claims, err := utils.ValidateToken(tokenString)
	if err != nil {
		log.Println("Token invalid:", err)
		conn.Close()
		return
	}

	// 4. Register client
	client := &Client{
		hub:      hub,
		conn:     conn,
		send:     make(chan *Message, 256),
		room:     c.Param("ticketId"),
		userId:   claims.Data.ID,
		username: claims.Data.Username,
		userRole: claims.Data.Role,
	}
	client.hub.register <- client

	go client.writePump()
	go client.readPump()
}
