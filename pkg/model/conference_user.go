package model

// ConferenceUser - conference_users relationship data model
type ConferenceUser struct {
	Base
	UserID       string `json:"user_id"`
	ConferenceID string `json:"conference_id"`

	Users      []*User     `json:"users" sql:"-"`
	Conference *Conference `json:"conference" sql:"-"`
}
