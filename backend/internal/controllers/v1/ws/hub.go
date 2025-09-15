package ws

import "fmt"

type Hub struct {
	clients    map[string]map[*Client]bool
	broadcast  chan *Message
	register   chan *Client
	unregister chan *Client
}

func NewHub() *Hub {
	return &Hub{
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]map[*Client]bool),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			if h.clients[client.room] == nil {
				h.clients[client.room] = make(map[*Client]bool)
			}
			h.clients[client.room][client] = true
			fmt.Println("size of conn", len(h.clients[client.room]))
		case client := <-h.unregister:
			if _, ok := h.clients[client.room][client]; ok {
				delete(h.clients[client.room], client)
				close(client.send)
			}
		case message := <-h.broadcast:
			for client := range h.clients[message.Room] {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients[message.Room], client)
				}
			}
		}
	}
}
