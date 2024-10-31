package trivia

import (
	"log"

	"github.com/gorilla/websocket"
)

// A `Client` acts as an intermediary between the websocket connection and a single
// instance of the `Hub` type.
type Client struct {
	hub  *Hub
	send chan []byte     // Buffered channel of outbound messages
	conn *websocket.Conn // The Websocket connection
}

// readPump pumps messages from the Websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("Error: %v", err)
			}
			break
		}

		c.hub.broadcast <- message
	}
}

// writePump pumps messages from the hub to the Websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
func (c *Client) writePump() {
	for message := range c.send {
		c.conn.WriteMessage(websocket.TextMessage, message)
	}
}
