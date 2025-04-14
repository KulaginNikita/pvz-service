package converter

import (
	"github.com/KulaginNikita/pvz-service/internal/domain/product"
	"github.com/KulaginNikita/pvz-service/internal/repository/productrepo/model"
)

func ToDB(p *product.Product) *model.Product {
	return &model.Product{
		ID:          p.ID,
		ReceptionID: p.ReceptionID,
		Type:        model.ProductType(p.Type),
		DateTime:    p.DateTime,
	}
}

