package product

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

var (
	ErrUnauthorized    = errors.New("user does not have 'employee' role")
	ErrNoOpenReception = errors.New("no open reception found")
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
