package repository

import (
	"errors"
)

var (
	ErrUserAlreadyExist = errors.New("user already exist")
	ErrUserNotFound     = errors.New("user not found")

	ErrTokenNotFound = errors.New("token not found")
	ErrTokenExpired  = errors.New("token expired")

	ErrItemAlreadyExist = errors.New("item already exist")
	ErrItemNotFound     = errors.New("item not found")
)
