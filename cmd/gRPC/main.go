package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"

	pvzgrpc "github.com/KulaginNikita/pvz-service/internal/gRPC/pvz"
	pb "github.com/KulaginNikita/pvz-service/pkg/pvz_v1/pvz" 
)

func buildDSN() string {
	user := os.Getenv("PG_USER")
	pass := os.Getenv("PG_PASSWORD")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_PORT")
	db := os.Getenv("PG_DATABASE_NAME")

	return fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=disable",
		host, port, db, user, pass)
}

func main() {
	dsn := buildDSN()
	pool, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer pool.Close()

	listener, err := net.Listen("tcp", ":3000")
	if err != nil {
		log.Fatalf("failed to listen on port 3000: %v", err)
	}

	server := grpc.NewServer()
	handler := pvzgrpc.NewHandler(pool) 
	pb.RegisterPVZServiceServer(server, handler)

	log.Println("🚀 gRPC server started on :3000")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

