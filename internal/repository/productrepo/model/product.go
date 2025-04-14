package model

import (
	"time"

	"github.com/google/uuid"
)

type ProductType string

const (
	ProductElectronics ProductType = "electronics"
	ProductClothes     ProductType = "clothes"
	ProductShoes       ProductType = "shoes"
)

type Product struct {
	ID          uuid.UUID
	ReceptionID uuid.UUID
	Type        ProductType
	DateTime    time.Time
}
