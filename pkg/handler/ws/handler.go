package ws

import (
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo"
	"github.com/phuwn/tools/errors"
	"github.com/phuwn/tools/log"
)

const (
	normalMessage = iota
	welcomeMessage
	leaveMessage

	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer.
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

// Message - transmitted packet data info
type Message struct {
	typ    int
	data   []byte
	client *Client
}

func (m *Message) toSendData() []byte {
	// send data structure: typ - clientID - data
	return append([]byte{byte(m.typ)}, append([]byte(m.client.id), m.data...)...)
}

// CheckOrigin returns true if the request's origin is in allowed list
func CheckOrigin(r *http.Request) bool {
	return true
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin:     CheckOrigin,
}

// Serve handles websocket requests from the peer.
func Serve(c echo.Context, hubID, clientID string) error {
	hub := getHub(hubID)
	if hub == nil {
		return errors.New("invalid hub")
	}
	conn, err := upgrader.Upgrade(c.Response().Writer, c.Request(), nil)
	if err != nil {
		log.Error(err)
		return errors.Customize(500, "failed to upgrade request", err)
	}
	client := &Client{id: clientID, hub: hub, conn: conn, send: make(chan *Message, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
	return nil
}
