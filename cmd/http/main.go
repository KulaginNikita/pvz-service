package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/KulaginNikita/pvz-service/internal/api"
	"github.com/KulaginNikita/pvz-service/internal/middleware"
	"github.com/KulaginNikita/pvz-service/internal/repository/productrepo"
	"github.com/KulaginNikita/pvz-service/internal/repository/pvzrepo"
	"github.com/KulaginNikita/pvz-service/internal/repository/receptionrepo"
	"github.com/KulaginNikita/pvz-service/internal/repository/userrepo"
	productservice "github.com/KulaginNikita/pvz-service/internal/service/product"
	pvzservice "github.com/KulaginNikita/pvz-service/internal/service/pvz"
	receptionservice "github.com/KulaginNikita/pvz-service/internal/service/reception"
	userservice "github.com/KulaginNikita/pvz-service/internal/service/user"
	"github.com/KulaginNikita/pvz-service/pkg/jwtutil"

	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func buildDSN() string {
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	db := os.Getenv("PG_DATABASE_NAME")

	fmt.Println(user, pass, host, port, db)
	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, db, user, pass)
}

func main() {
	dsn := buildDSN()
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("couldn't connect to db: %v", err)
	}
	defer pool.Close()

	jwtManager := jwtutil.NewManager("my-secret-key") 
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	userRepo := userrepo.NewUserRepository(pool)
	pvzRepo := pvzrepo.NewPVZRepository(pool)
	receptionRepo := receptionrepo.NewReceptionRepository(pool)
	productrepo := productrepo.NewProductRepository(pool)

	userService := userservice.NewUserService(userRepo, jwtManager)
	receptionService := receptionservice.NewReceptionService(receptionRepo, jwtManager)
	pvzService := pvzservice.NewPVZService(pvzRepo, jwtManager)
	productService := productservice.NewProductService(productrepo, jwtManager)

	ap := api.NewAPI(userService, pvzService, receptionService, productService, jwtManager)

	wrapper := api.ServerInterfaceWrapper{
		Handler: ap,
		ErrorHandlerFunc: func(w http.ResponseWriter, r *http.Request, err error) {
			http.Error(w, err.Error(), http.StatusBadRequest)
		},
	}

	r := chi.NewRouter()

	r.Use(middleware.LoggerMiddleware(logger))

	r.Group(func(r chi.Router) {
		r.Post("/login", wrapper.PostLogin)
		r.Post("/register", wrapper.PostRegister)
		r.Post("/dummyLogin", wrapper.PostDummyLogin)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.JWTAuthMiddleware(jwtManager))
		r.Post("/products", wrapper.PostProducts)
		r.Get("/pvz", wrapper.GetPvz)
		r.Post("/pvz", wrapper.PostPvz)
		r.Post("/pvz/{pvzId}/close_last_reception", wrapper.PostPvzPvzIdCloseLastReception)
		r.Post("/pvz/{pvzId}/delete_last_product", wrapper.PostPvzPvzIdDeleteLastProduct)
		r.Post("/receptions", wrapper.PostReceptions)
	})

	log.Println("🚀 Server started at http://localhost:8080")
	if err := http.ListenAndServe(":8080", r); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}
