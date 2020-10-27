package v1

import (
	"github.com/google/uuid"
)

// UuidGenerator ...
type UuidGenerator interface {
	Generate() string
}

type uuidGenerator struct{}

func (g *uuidGenerator) Generate() string {
	return uuid.New().String()
}

// NewUuidGenerator ...
func NewUuidGenerator() UuidGenerator {
	return &uuidGenerator{}
}
