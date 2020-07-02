package ws

var server *Server

// NewServer - create and run a new ws server
func NewServer() {
	server = &Server{
		hubs:       make(map[string]*Hub),
		register:   make(chan *Hub),
		unregister: make(chan string),
	}
	server.run()
}

// Server maintains the set of active hubs
type Server struct {
	// Registered hubs
	hubs map[string]*Hub

	// Register hub
	register chan *Hub

	// Unregister hub
	unregister chan string
}

func (s *Server) run() {
	for {
		select {
		case hub := <-s.register:
			s.hubs[hub.id] = hub
		case id := <-s.unregister:
			if hub, ok := s.hubs[id]; ok {
				hub.terminate <- true
				close(hub.broadcast)
				close(hub.register)
				close(hub.unregister)
				close(hub.terminate)
				delete(s.hubs, id)
			}
		}
	}
}

func getHub(id string) *Hub {
	if hub, ok := server.hubs[id]; ok {
		return hub
	}
	return nil
}
