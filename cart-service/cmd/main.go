package main

import (
	"fmt"
	"log"
	"net"

	database "cart-service/config"
	"cart-service/internal/handlers"
	"cart-service/internal/repository"
	"cart-service/internal/usecase"

	cartProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/cart-pb"
	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	database.ConnectDB()

	grpcServer := grpc.NewServer()

	// Enable reflection for debugging with grpcurl
	reflection.Register(grpcServer)

	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Product Service: %v", err)
	}
	defer productConn.Close()
	productClient := productProto.NewProductServiceClient(productConn)

	cartRepo := repository.NewCartRepository(database.DB)

	cartUC := usecase.NewCartUsecase(cartRepo, productClient)

	cartHandler := handlers.NewCartHandler(cartUC)

	cartProto.RegisterCartServiceServer(grpcServer, cartHandler)

	listener, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Fatalf("Failed to listen on port 50053: %v", err)
	}

	fmt.Println("gRPC Server is running on port 50053...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
