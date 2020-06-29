package server

import "github.com/newzoom/api/pkg/store"

var srv Server

// Server - server core structure
type Server struct {
	store *store.Store
}

// NewServerCfg - create new server
func NewServerCfg(store *store.Store) {
	srv.store = store
}

// GetServerCfg - get server param
func GetServerCfg() *Server {
	return &srv
}

// Store - get store
func (s *Server) Store() *store.Store {
	return s.store
}
