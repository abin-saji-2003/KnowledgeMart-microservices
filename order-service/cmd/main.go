package main

import (
	"fmt"
	"log"
	"net"

	database "order-service/config"
	"order-service/internal/handlers"
	"order-service/internal/repository"
	"order-service/internal/usecase"

	cartProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/cart-pb"
	orderProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/order-pb"
	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
)

func main() {
	database.ConnectDB()

	grpcServer := grpc.NewServer()

	reflection.Register(grpcServer)

	// Connect to Cart Service via gRPC
	cartConn, err := grpc.Dial("localhost:50053", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Cart Service: %v", err)
	}
	defer cartConn.Close()
	cartClient := cartProto.NewCartServiceClient(cartConn)

	productConn, err := grpc.Dial("localhost:50052", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect to Product Service: %v", err)
	}
	defer productConn.Close()
	productClient := productProto.NewProductServiceClient(productConn)

	orderRepo := repository.NewOrderRepository(database.DB)

	orderUC := usecase.NewOrderUseCase(orderRepo, cartClient, productClient)

	orderHandler := handlers.NewOrderHandler(orderUC)

	orderProto.RegisterOrderServiceServer(grpcServer, orderHandler)

	listener, err := net.Listen("tcp", ":50054")
	if err != nil {
		log.Fatalf("Failed to listen on port 50054: %v", err)
	}

	fmt.Println("gRPC Order Service is running on port 50054...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
