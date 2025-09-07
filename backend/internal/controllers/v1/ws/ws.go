package ws

import (
	"backend/internal/redis"
	"backend/utils"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	connections = make(map[string]*websocket.Conn)
	connMutex   sync.RWMutex
)

// setConnState sets connection state in Redis
func setConnState(connID, field, value string) {
	redis.Rdb.HSet(redis.Ctx, "ws:"+connID, field, value)
	redis.Rdb.Expire(redis.Ctx, "ws:"+connID, time.Hour)
}

// getConnState gets connection state from Redis
func getConnState(connID, field string) string {
	val, _ := redis.Rdb.HGet(redis.Ctx, "ws:"+connID, field).Result()
	return val
}

// isAuthenticated checks if connection is authenticated
func isAuthenticated(connID string) bool {
	return getConnState(connID, "auth") == "true"
}

func HandleWS(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "failed to upgrade"})
		return
	}
	defer conn.Close()

	connID := uuid.New().String()

	connMutex.Lock()
	connections[connID] = conn
	connMutex.Unlock()

	defer func() {
		connMutex.Lock()
		delete(connections, connID)
		connMutex.Unlock()
		redis.Rdb.Del(redis.Ctx, "ws:"+connID)
	}()

	conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf(`{"connectionId":"%s"}`, connID)))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		var data map[string]string
		if err := json.Unmarshal(msg, &data); err != nil {
			continue
		}

		switch data["type"] {
		case "AUTH":
			claims, err := utils.ValidateToken(data["token"])
			if err != nil {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"unauthorized"}`))
				continue
			}

			setConnState(connID, "auth", "true")
			setConnState(connID, "userId", claims.ID)
			conn.WriteMessage(websocket.TextMessage, []byte(`{"message":"authenticated"}`))

		case "SUBSCRIBE_TICKET":
			if !isAuthenticated(connID) {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"not authenticated"}`))
				continue
			}

			setConnState(connID, "ticketId", data["ticketId"])
			conn.WriteMessage(websocket.TextMessage, []byte(`{"message":"subscribed"}`))

		case "COMMENT":
			if !isAuthenticated(connID) {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"not authenticated"}`))
				continue
			}

			ticketID := getConnState(connID, "ticketId")
			if data["ticketId"] != ticketID {
				conn.WriteMessage(websocket.TextMessage, []byte(`{"error":"not subscribed"}`))
				continue
			}

			userID := getConnState(connID, "userId")
			reply := fmt.Sprintf(`{"ticketId":"%s","userId":"%s","comment":"%s"}`,
				ticketID, userID, data["comment"])
			conn.WriteMessage(websocket.TextMessage, []byte(reply))
		}
	}
}
