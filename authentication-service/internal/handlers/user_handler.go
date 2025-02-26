package handler

import (
	"context"
	"log"

	"authentication-service/internal/usecase"

	authProto "github.com/abin-saji-2003/GRPC-knowledgemart/proto/auth-pb"
)

type AuthHandler struct {
	authUC *usecase.AuthUseCase
	authProto.UnimplementedAuthServiceServer
}

func NewAuthHandler(authUC *usecase.AuthUseCase) *AuthHandler {
	return &AuthHandler{authUC: authUC}
}

func (h *AuthHandler) EmailSignup(ctx context.Context, req *authProto.EmailSignupRequest) (*authProto.EmailSignupResponse, error) {
	log.Print("pass:", req.Password)
	log.Print("con:", req.ConfirmPassword)
	_, err := h.authUC.EmailSignup(req.Name, req.Email, req.PhoneNumber, req.Password, req.ConfirmPassword)
	if err != nil {
		return &authProto.EmailSignupResponse{Success: false, Message: err.Error()}, err
	}

	return &authProto.EmailSignupResponse{
		Success: true,
		Message: "Signup successful, Please login",
	}, nil
}

func (h *AuthHandler) EmailLogin(ctx context.Context, req *authProto.EmailLoginRequest) (*authProto.EmailLoginResponse, error) {
	token, user, err := h.authUC.EmailLogin(req.Email, req.Password)
	if err != nil {
		return &authProto.EmailLoginResponse{
			Success: false,
			Token:   "",
			User:    nil,
		}, err
	}
	return &authProto.EmailLoginResponse{
		Success: true,
		Token:   token,
		User: &authProto.UserResponse{
			Id:          uint32(user.ID),
			Name:        user.Name,
			Email:       user.Email,
			PhoneNumber: user.PhoneNumber,
			Picture:     user.Picture,
			Blocked:     user.Blocked,
		},
	}, nil
}
