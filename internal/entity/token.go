package entity

import (
	"time"
)

type Token struct {
	ID        string
	UserLogin string
	CreatedAt time.Time
}
