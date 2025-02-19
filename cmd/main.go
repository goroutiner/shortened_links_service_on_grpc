package main

import (
	"fmt"
	"log"
	"net"
	"shortened_links_service_on_grpc/internal/config"
	"shortened_links_service_on_grpc/internal/entities"
	"shortened_links_service_on_grpc/internal/handlers"
	"shortened_links_service_on_grpc/internal/services"
	"shortened_links_service_on_grpc/internal/storage"
	"shortened_links_service_on_grpc/internal/storage/database"
	"shortened_links_service_on_grpc/internal/storage/memory"
	pb "shortened_links_service_on_grpc/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	var (
		store storage.StorageInterface
		err   error
	)

	go handlers.СleanupVisitors()

	switch config.Mode {
	case "in-memory":
		store = memory.NewMemoryStore()
		log.Println("Using in-memory storage")
	case "postgres":
		entities.Db, err = database.NewDatabaseStore(config.PsqlUrl)
		if err != nil {
			log.Fatal(err.Error())
		}
		defer entities.Db.Close()

		store = database.NewDatabaseConection(entities.Db)
		log.Println("Using PostgreSQL store")
	default:
		log.Fatalf("config.Mode is empty in /internal/config/setting.go")
	}

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	grpcServer := grpc.NewServer(grpc.UnaryInterceptor(handlers.RateLimitInterceptor()))
	service := services.NewShortenerService(store)
	handler := handlers.RegisterShortenerHandler(service)

	pb.RegisterShortenerServiceServer(grpcServer, handler)
	reflection.Register(grpcServer)

	fmt.Println("Service is running ...")
	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("Failed to serve gRPC server: %v", err)
	}
}
