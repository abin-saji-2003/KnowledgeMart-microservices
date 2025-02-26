package usecase

import (
	"authentication-service/internal/models"
	"authentication-service/internal/repository"
	"authentication-service/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

type AuthUseCase struct {
	userRepo   repository.UserRepository
	sellerRepo repository.SellerRepository
}

func NewAuthUseCase(userRepo repository.UserRepository, sellerRepo repository.SellerRepository) *AuthUseCase {
	return &AuthUseCase{
		userRepo:   userRepo,
		sellerRepo: sellerRepo,
	}
}

func (uc *AuthUseCase) EmailSignup(name, email, phone, password, confirmPassword string) (*models.User, error) {
	if password != confirmPassword {
		return nil, errors.New("passwords do not match")
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// otp, otpExpiry := utils.GenerateOTP()

	newUser := &models.User{
		Name:        name,
		Email:       email,
		PhoneNumber: phone,
		Password:    hashedPassword,
		LoginMethod: "email",
		Blocked:     false,
	}

	existingUser, err := uc.userRepo.GetByEmail(email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("database error")
	}

	log.Print(existingUser)
	if existingUser != nil && existingUser.ID != 0 {
		return nil, errors.New("user already exists")
	}

	if err := uc.userRepo.Create(newUser); err != nil {
		return nil, errors.New("failed to create user")
	}

	// err = utils.SendOTPEmail(email, otp)
	// if err != nil {
	// 	return nil, err
	// }
	return newUser, nil
}

func (uc *AuthUseCase) EmailLogin(email, password string) (string, *models.User, error) {
	user, err := uc.userRepo.GetByEmail(email)

	if err != nil {
		return "", nil, errors.New("database error")
	}

	if user == nil {
		return "", nil, errors.New("invalid email or password")
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return "", nil, errors.New("incorrect password")
	}

	if user.Blocked {
		return "", nil, errors.New("user is not authorized to access")
	}

	token, err := utils.GenerateJWT(user.ID, "user")
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, user, nil
}
