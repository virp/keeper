package entity

import (
	"time"
)

// Item struct represent business entity for stored items in Keeper.
type Item struct {
	ID        string
	UserID    string
	Name      string
	Type      string
	Data      []byte
	Metadata  []Metadata
	CreatedAt time.Time
	UpdatedAt time.Time
}
