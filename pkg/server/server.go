package server

import (
	"github.com/newzoom/api/pkg/service"
	"github.com/newzoom/api/pkg/store"
)

var srv Server

// Server - server core structure
type Server struct {
	store   *store.Store
	service *service.Service
}

// NewServerCfg - create new server
func NewServerCfg(store *store.Store, service *service.Service) {
	srv.store = store
	srv.service = service
}

// GetServerCfg - get server param
func GetServerCfg() *Server {
	return &srv
}

// Store - get store
func (s *Server) Store() *store.Store {
	return s.store
}

// Service - get service
func (s *Server) Service() *service.Service {
	return s.service
}
