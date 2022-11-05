package repository

import (
	"errors"
)

var (
	ErrItemAlreadyExist = errors.New("item already exist")
	ErrItemNotFound     = errors.New("item not found")
)
