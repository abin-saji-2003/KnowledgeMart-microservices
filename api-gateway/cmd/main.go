package main

import (
	"api-gateway/internal/handlers"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	authProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/auth-pb"
	cartProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/cart-pb"
	orderProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/order-pb"
	productProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/product-pb"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func connectGRPC(serviceName, address string) (*grpc.ClientConn, error) {
	conn, err := grpc.Dial(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Failed to connect to %s gRPC server: %v", serviceName, err)
		return nil, err
	}
	log.Printf("✅ Connected to %s gRPC server at %s", serviceName, address)
	return conn, nil
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	r := gin.Default()

	services := map[string]string{
		"Auth":    os.Getenv("AUTH_SERVICE_ADDR"),
		"Product": os.Getenv("PRODUCT_SERVICE_ADDR"),
		"Cart":    os.Getenv("CART_SERVICE_ADDR"),
		"Order":   os.Getenv("ORDER_SERVICE_ADDR"),
	}

	// Connect to gRPC services
	authConn, err := connectGRPC("Auth", services["Auth"])
	if err != nil {
		log.Fatal("Exiting due to connection failure")
	}
	defer authConn.Close()
	authClient := authProto.NewAuthServiceClient(authConn)

	productConn, err := connectGRPC("Product", services["Product"])
	if err != nil {
		log.Fatal("Exiting due to connection failure")
	}
	defer productConn.Close()
	productClient := productProto.NewProductServiceClient(productConn)

	cartConn, err := connectGRPC("Cart", services["Cart"])
	if err != nil {
		log.Fatal("Exiting due to connection failure")
	}
	defer cartConn.Close()
	cartClient := cartProto.NewCartServiceClient(cartConn)

	orderConn, err := connectGRPC("Order", services["Order"])
	if err != nil {
		log.Fatal("Exiting due to connection failure")
	}
	defer orderConn.Close()
	orderClient := orderProto.NewOrderServiceClient(orderConn)

	handlers.RegisterAuthRoutes(r, authClient)

	handlers.AddProduct(r, productClient)

	handlers.RegisterCartRoutes(r, cartClient)

	orderHandler := handlers.NewOrderHandler(orderClient)
	orderHandler.RegisterOrderRoutes(r)

	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"message": "API Gateway is Running!"})
	})

	log.Println("API Gateway running on port 8080...")
	if err := r.Run(":8080"); err != nil {
		log.Fatalf("Failed to start API Gateway: %v", err)
	}
}
