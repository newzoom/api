package model

// Conference data model
type Conference struct {
	Base
	ID          string `json:"id"`
	Topic       string `json:"topic"`
	Description string `json:"description"`
	Password    string `json:"-"`
	HostID      string `json:"host_id"`
	IsActive    bool   `json:"is_active"`

	Users           []*User           `json:"users,omitempty" sql:"-"`
	ConferenceUsers []*ConferenceUser `json:"-" gorm:"foreignkey:ConferenceID;association_foreignkey:ID"`
}
