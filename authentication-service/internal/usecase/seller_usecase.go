package usecase

import (
	"authentication-service/internal/models"
	"authentication-service/utils"
	"errors"
	"log"

	"gorm.io/gorm"
)

func (uc *AuthUseCase) SellerSignup(userId uint, name string, description string, password string) (*models.Seller, error) {

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	_, err = uc.userRepo.GetByID(userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("database error")
	}

	// otp, otpExpiry := utils.GenerateOTP()

	newSeller := &models.Seller{
		UserID:      userId,
		UserName:    name,
		Password:    hashedPassword,
		Description: description,
	}

	existingSeller, err := uc.sellerRepo.GetByUserId(userId)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, errors.New("database error")
	}

	log.Print(existingSeller)
	if existingSeller != nil && existingSeller.ID != 0 {
		return nil, errors.New("user already exists")
	}

	_, err = uc.sellerRepo.Create(newSeller)
	if err != nil {
		return nil, errors.New("failed to create user")
	}

	// err = utils.SendOTPEmail(email, otp)
	// if err != nil {
	// 	return nil, err
	// }
	return newSeller, nil
}

func (uc *AuthUseCase) SellerLogin(Id uint, userName, password string) (string, *models.Seller, error) {
	seller, err := uc.sellerRepo.GetByUserName(userName)

	if err != nil {
		return "", nil, errors.New("database error")
	}

	if seller == nil {
		return "", nil, errors.New("invalid email or password")
	}

	if err := utils.CheckPassword(seller.Password, password); err != nil {
		return "", nil, errors.New("incorrect password")
	}

	token, err := utils.GenerateJWT(seller.ID, "seller")
	if err != nil {
		return "", nil, errors.New("failed to generate token")
	}

	return token, seller, nil
}
