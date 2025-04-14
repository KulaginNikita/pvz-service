package model

import (
	"github.com/google/uuid"
	"time"
)

type Reception struct {
	ID     uuid.UUID `db:"id"`
	PVZID  uuid.UUID     `db:"pvz_id"`
	Status string    `db:"status"`
	StartedAt time.Time `db:"started_at"`
}

type PVZID struct {
	PVZID uuid.UUID `db:"pvz_id"`
}