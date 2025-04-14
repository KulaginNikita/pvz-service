package api

import (
	"github.com/KulaginNikita/pvz-service/internal/service"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"
)

type API struct {
	Unimplemented
	receptionService service.ReceptionService
	pvzService service.PVZService
	userService service.UserService
	productService service.ProductService
	jwtManager  *jwtutil.Manager
}

func NewAPI(userService service.UserService, 
			pvzService service.PVZService, 
			receptionService service.ReceptionService, 
			productService service.ProductService,
			jwtManager *jwtutil.Manager) *API {
	return &API{
		userService: userService,
		pvzService: pvzService,
		receptionService: receptionService,
		productService: productService,
		jwtManager:  jwtManager,
	}
}
