package ws

import (
	"encoding/json"
	"fmt"

	"github.com/phuwn/tools/log"
)

var (
	hub *Hub
)

const (
	welcomeMsg = iota
	offerMsg
	answerMsg
	exchangeCandidateMsg
	leaveMsg
)

// Message - data packet that transfered between hub and client through ws conn
type Message struct {
	Typ        int         `json:"typ"`
	Data       interface{} `json:"data"`
	ReceiverID string      `json:"receiver_id"`
	SenderID   string      `json:"sender_id"`

	client *Client
}

// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	conferences map[string]map[string]*Client

	// Inbound messages from the clients.
	broadcast chan *Message

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client
}

// NewHub - create and run a new hub
func NewHub() {
	hub = &Hub{
		broadcast:   make(chan *Message),
		register:    make(chan *Client),
		unregister:  make(chan *Client),
		conferences: make(map[string]map[string]*Client),
	}
	hub.run()
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if con, ok := h.conferences[client.conID]; ok {
				con[client.id] = client
			} else {
				h.conferences[client.conID] = make(map[string]*Client)
				h.conferences[client.conID][client.id] = client
			}
		case client := <-h.unregister:
			if conference, ok := h.conferences[client.conID]; ok {
				if client, ok := conference[client.id]; ok {
					close(client.send)
					delete(conference, client.id)
				}
			}
		case message := <-h.broadcast:
			h.routeMessage(message)
		}
	}
}

func (h *Hub) welcomeRequest(msg *Message) {
	con, ok := h.conferences[msg.client.conID]
	if !ok {
		return
	}
	for _, client := range con {
		if client != msg.client {
			msg.ReceiverID = client.id
			b, err := json.Marshal(msg)
			if err != nil {
				log.Errorf(err, "failed to marshal message %v", msg)
				return
			}
			select {
			case client.send <- b:
			default:
				close(client.send)
				delete(con, client.id)
			}
		}
	}
}

func (h *Hub) offerRequest(msg *Message) {
	con, ok := h.conferences[msg.client.conID]
	if !ok {
		return
	}
	if client, ok := con[msg.ReceiverID]; ok {
		b, err := json.Marshal(msg)
		if err != nil {
			log.Errorf(err, "failed to marshal message %v", msg)
			return
		}
		client.send <- b
		return
	}
	log.Error(fmt.Errorf("invalid receiver %v", msg.ReceiverID))
}

func (h *Hub) answerRequest(msg *Message) {
	con, ok := h.conferences[msg.client.conID]
	if !ok {
		return
	}
	if client, ok := con[msg.ReceiverID]; ok {
		b, err := json.Marshal(msg)
		if err != nil {
			log.Errorf(err, "failed to marshal message %v", msg)
			return
		}
		client.send <- b
		return
	}
	log.Error(fmt.Errorf("invalid receiver %v", msg.ReceiverID))
}

func (h *Hub) exchangeCandidateRequest(msg *Message) {
	con, ok := h.conferences[msg.client.conID]
	if !ok {
		return
	}
	if client, ok := con[msg.ReceiverID]; ok {
		b, err := json.Marshal(msg)
		if err != nil {
			log.Errorf(err, "failed to marshal message %v", msg)
			return
		}
		client.send <- b
		return
	}
	log.Error(fmt.Errorf("invalid receiver %v", msg.ReceiverID))
}

func (h *Hub) leaveRequest(msg *Message) {
	con, ok := h.conferences[msg.client.conID]
	if !ok {
		return
	}
	for _, client := range con {
		if client != msg.client {
			msg.ReceiverID = client.id
			b, err := json.Marshal(msg)
			if err != nil {
				log.Errorf(err, "failed to marshal message %v", msg)
				return
			}
			select {
			case client.send <- b:
			default:
				close(client.send)
				delete(con, client.id)
			}
		}
	}
}

func (h *Hub) routeMessage(msg *Message) {
	switch msg.Typ {
	case welcomeMsg:
		h.welcomeRequest(msg)
	case answerMsg:
		h.answerRequest(msg)
	case offerMsg:
		h.offerRequest(msg)
	case exchangeCandidateMsg:
		h.exchangeCandidateRequest(msg)
	case leaveMsg:
		h.leaveRequest(msg)
	default:
		log.Error(fmt.Errorf("invalid msg type %v", msg.Typ))
	}
}
