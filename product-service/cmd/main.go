package main

import (
	"fmt"
	"log"
	"net"

	database "product-service/config"
	"product-service/internal/handlers"
	"product-service/internal/repository"
	"product-service/internal/usecase"

	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	database.ConnectDB()

	grpcServer := grpc.NewServer()

	// Enable reflection for debugging with grpcurl
	reflection.Register(grpcServer)

	productRepo := repository.NewProductRepository(database.DB)

	productUC := usecase.NewProductUseCase(productRepo)

	ProductHandler := handlers.NewProductHandler(productUC)

	// Register gRPC Service
	productProto.RegisterProductServiceServer(grpcServer, ProductHandler)

	// Start gRPC Listener
	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("Failed to listen on port 50052: %v", err)
	}

	fmt.Println("gRPC Server is running on port 50052...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
