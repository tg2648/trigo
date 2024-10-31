package trivia

import "log"

// The `Hub` maintains a set of registered clients and
// broadcasts messages to the clients.
type Hub struct {
	clients    map[*Client]bool // Registered clients
	broadcast  chan []byte      // Inbound messages from clients
	register   chan *Client     // Requests from clients to join a room
	unregister chan *Client     // Requests from clients to leave a room
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			log.Println("Client registered")
			h.clients[client] = true
		case client := <-h.unregister:
			log.Println("Client unregistered")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println("Client sent message")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					log.Println("Could not broadcast to client")
					// Could not send to client
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
