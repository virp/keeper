package entity

import (
	"time"
)

// User struct represent business entity for stored Keeper user.
type User struct {
	ID           string
	Login        string
	PasswordHash string
	CreatedAt    time.Time
}
