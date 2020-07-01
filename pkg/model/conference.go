package model

// Conference data model
type Conference struct {
	Base
	ID   string `json:"id"`
	Name string `json:"name"`
}
