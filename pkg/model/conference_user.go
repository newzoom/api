package model

// ConferenceUser - conference_users relationship data model
type ConferenceUser struct {
	Base
	UserID       int `json:"user_id"`
	ConferenceID int `json:"conference_id"`

	Users      []*User     `json:"users" sql:"-"`
	Conference *Conference `json:"conference" sql:"-"`
}
