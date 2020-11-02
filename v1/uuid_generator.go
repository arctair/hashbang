package v1

import uuid "github.com/google/uuid"

// UUIDGenerator ...
type UUIDGenerator interface {
	Generate() string
}

type uuidGenerator struct{}

func (g *uuidGenerator) Generate() string {
	return uuid.New().String()
}

// NewUUIDGenerator ...
func NewUUIDGenerator() UUIDGenerator {
	return &uuidGenerator{}
}
