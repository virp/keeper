package services

import (
	"github.com/google/uuid"
)

type IdGenerator interface {
	Generate() string
}

type UuidGenerator struct {
}

func (g *UuidGenerator) Generate() string {
	return uuid.NewString()
}
