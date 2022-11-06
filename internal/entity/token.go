package entity

import (
	"time"
)

type Token struct {
	ID        string
	UserID    string
	CreatedAt time.Time
}
