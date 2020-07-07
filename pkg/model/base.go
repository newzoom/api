package model

import "time"

// Base - base model for quick add base field
type Base struct {
	CreatedAt *time.Time `json:"-"`
	UpdatedAt *time.Time `json:"-"`
	DeletedAt *time.Time `json:"-"`
}
