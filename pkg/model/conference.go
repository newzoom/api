package model

// Conference data model
type Conference struct {
	Base
	ID   int    `json:"id"`
	Name string `json:"name"`
}
