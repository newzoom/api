package model

// User data model
type User struct {
	Base
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}
