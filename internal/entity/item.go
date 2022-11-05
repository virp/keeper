package entity

import (
	"time"
)

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
