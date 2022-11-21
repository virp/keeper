package entity

import (
	"time"
)

// Token struct represent business entity for stored user token.
type Token struct {
	ID        string
	UserLogin string
	CreatedAt time.Time
}
