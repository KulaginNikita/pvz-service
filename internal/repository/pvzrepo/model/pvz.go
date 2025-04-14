package model

import (
	"time"

	"github.com/google/uuid"
)

type PVZ struct {
	ID           uuid.UUID `db:"id"`
	City         string    `db:"city"`
	RegisteredAt time.Time `db:"registered_at"`
}
