package store

import (
	"github.com/newzoom/api/pkg/store/conference"
	"github.com/newzoom/api/pkg/store/user"
)

// Store - server store struct
type Store struct {
	User       user.Store
	Conference conference.Store
}

// New - create new store variable
func New() *Store {
	return &Store{
		User:       user.NewStore(),
		Conference: conference.NewStore(),
	}
}
