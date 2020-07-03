package ws

import (
	"time"

	"github.com/gorilla/websocket"
	"github.com/phuwn/tools/log"
)

// Client is a middleman between the user's websocket connection and the hub.
type Client struct {
	id   string
	conn *websocket.Conn
	hub  *Hub

	// Buffered channel of outbound messages.
	send chan *Message
}

// readPump pumps messages from the websocket connection to the hub.
func (c *Client) readPump() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			log.Error(err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// handle unexpected close error
			}
			c.hub.broadcast <- &Message{leaveMessage, []byte{}, c}
			break
		}
		c.hub.broadcast <- &Message{normalMessage, message, c}
	}
}

// writePump pumps messages from the hub to the websocket connection.
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
				// The hub closed the channel.
				c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.conn.NextWriter(websocket.BinaryMessage)
			if err != nil {
				log.Error(err)
				return
			}
			w.Write(message.toSendData())

			// Add queued chat messages to the current websocket message.
			// for i := 0; i < len(c.send); i++ {
			// 	w.Write((<-c.send).data)
			// }

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
