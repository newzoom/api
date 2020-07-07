package ws

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Client is a middleman between the websocket connection and the hub.
type Client struct {
	id string

	conID string

	// The websocket connection.
	conn *websocket.Conn

	// Buffered channel of outbound messages.
	send chan []byte
}

// readPump pumps messages from the websocket connection to the hub.
//
// The application runs readPump in a per-connection goroutine. The application
// ensures that there is at most one reader on a connection by executing all
// reads from this goroutine.
func (c *Client) readPump() {
	defer func() {
		hub.broadcast <- &Message{Typ: leaveMsg, SenderID: c.id, client: c}
		hub.unregister <- c
		c.conn.Close()
	}()
	// c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error { c.conn.SetReadDeadline(time.Now().Add(pongWait)); return nil })
	for {
		_, raw, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Println("failed to read message", err)
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				// handle unexpected close error
			}
			hub.broadcast <- &Message{Typ: leaveMsg, SenderID: c.id, client: c}
			break
		}

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		msg := &Message{client: c}
		err = json.Unmarshal(raw, msg)
		if err != nil {
			fmt.Println("failed to parse msg", string(raw))
			continue
		}
		msg.SenderID = c.id
		hub.broadcast <- msg
	}
}

// writePump pumps messages from the hub to the websocket connection.
//
// A goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by
// executing all writes from this goroutine.
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

			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				fmt.Println("failed to create writer", string(message))
				continue
			}

			w.Write(message)

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

// Serve handles websocket requests from the peer.
func Serve(c echo.Context, conID, userID string) error {
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		return errors.Customize(500, "failed to upgrade connection", err)
	}
	if con, ok := hub.conferences[conID]; ok && con[userID] != nil {
		return errors.Customize(400, "you've already joined this conference on a different browser please end the current session first", fmt.Errorf("session is in used"))
	}
	client := &Client{userID, conID, conn, make(chan []byte, 256)}
	hub.register <- client
	hub.broadcast <- &Message{Typ: welcomeMsg, SenderID: userID, client: client}

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}
