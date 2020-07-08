package store

import (
	"github.com/newzoom/api/pkg/store/conference"
	"github.com/newzoom/api/pkg/store/conferenceuser"
	"github.com/newzoom/api/pkg/store/user"
)

// Store - database handling implementation
type Store struct {
	User           user.Store
	Conference     conference.Store
	ConferenceUser conferenceuser.Store
}

// New - create new store variable
func New() *Store {
	return &Store{
		User:           user.NewStore(),
		Conference:     conference.NewStore(),
		ConferenceUser: conferenceuser.NewStore(),
	}
}
