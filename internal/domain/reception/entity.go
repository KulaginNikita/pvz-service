package reception

import (
	"time"
	"github.com/google/uuid"
)

type Status string

const (
	StatusInProgress Status = "in_progress"
	StatusClosed     Status = "close"
)

type Reception struct {
	ID     uuid.UUID
	PVZID  uuid.UUID
	Status Status
	StartedAt time.Time
}
