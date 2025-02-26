package handler

import (
	"context"
	"log"

	authProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/auth-pb"
)

func (h *AuthHandler) SellerRegister(ctx context.Context, req *authProto.SellerRegisterRequest) (*authProto.SellerRegisterResponse, error) {
	log.Println("Seller signup request received for user ID:", req.UserId) // Using ctx indirectly for logging

	_, err := h.authUC.SellerSignup(uint(req.UserId), req.Name, req.Description, req.Password)
	if err != nil {
		return &authProto.SellerRegisterResponse{Success: false, Message: err.Error()}, err
	}
	return &authProto.SellerRegisterResponse{
		Success: true,
		Message: "Signup successful, please login",
	}, nil
}

func (h *AuthHandler) SellerLogin(ctx context.Context, req *authProto.SellerLoginRequest) (*authProto.SellerLoginResponse, error) {
	token, seller, err := h.authUC.SellerLogin(uint(req.UserId), req.Username, req.Password)
	if err != nil {
		return &authProto.SellerLoginResponse{
			Success: false,
			Token:   "",
			Seller:  nil,
		}, err
	}

	// Convert models.Seller to authProto.SellerResponse
	sellerResponse := &authProto.SellerResponse{
		Id:          uint32(seller.ID),
		UserId:      uint32(seller.UserID),
		Name:        seller.UserName,
		Description: seller.Description,
	}

	return &authProto.SellerLoginResponse{
		Success: true,
		Token:   token,
		Seller:  sellerResponse,
	}, nil
}
