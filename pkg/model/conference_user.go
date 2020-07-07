package model

// ConferenceUser - conference_users relationship data model
type ConferenceUser struct {
	Base
	UserID       string `json:"user_id"`
	ConferenceID string `json:"conference_id"`

	User       *User       `json:"-" gorm:"foreignkey:ID;association_foreignkey:UserID"`
	Conference *Conference `json:"-" gorm:"foreignkey:ID;association_foreignkey:ConferenceID"`
}
