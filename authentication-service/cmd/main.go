package main

import (
	"fmt"
	"log"
	"net"

	"authentication-service/config"
	"authentication-service/internal/handlers"
	"authentication-service/internal/repository"
	"authentication-service/internal/usecase"

	authProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/auth-pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	database.ConnectDB()

	grpcServer := grpc.NewServer()

	// Enable reflection for debugging with grpcurl
	reflection.Register(grpcServer)

	userRepo := repository.NewUserRepository(database.DB)
	sellerRepo := repository.NewSellerRepository(database.DB)

	authUC := usecase.NewAuthUseCase(userRepo, sellerRepo)

	authHandler := handler.NewAuthHandler(authUC)

	authProto.RegisterAuthServiceServer(grpcServer, authHandler)

	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	fmt.Println("gRPC Server is running on port 50051...")
	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
