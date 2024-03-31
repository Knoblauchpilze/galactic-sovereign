package persistence

import "github.com/google/uuid"

type ApiKey struct {
	Id      uuid.UUID
	Key     uuid.UUID
	ApiUser uuid.UUID
}
