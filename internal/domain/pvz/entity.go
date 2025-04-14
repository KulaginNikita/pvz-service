package pvz

import (
	"time"

	"github.com/google/uuid"
)

type City string

const (
	CityMoscow City = "Москва"
	CitySPB    City = "Санкт-Петербург"
	CityKazan  City = "Казань"
)

type PVZ struct {
	ID   uuid.UUID
	City City
	RegisteredAt time.Time
}
