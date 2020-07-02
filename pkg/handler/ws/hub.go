package ws

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	id string
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// termiante signal channel
	terminate chan bool
}

// NewHub - create and run a new hub
func NewHub(id string) {
	server.register <- &Hub{
		id:         id,
		broadcast:  make(chan *Message),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		terminate:  make(chan bool),
	}
}

func (h *Hub) broadcastMessage(m *Message) {
	for client := range h.clients {
		if client.id != m.client.id {
			select {
			case client.send <- m:
			default:
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			h.broadcastMessage(&Message{welcomeMessage, []byte{}, client})
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				h.broadcastMessage(&Message{leaveMessage, []byte{}, client})
				close(client.send)
				delete(h.clients, client)
			}
		case message := <-h.broadcast:
			h.broadcastMessage(message)
		case <-h.terminate:
			for client := range h.clients {
				close(client.send)
				delete(h.clients, client)
			}
		}
	}
}
