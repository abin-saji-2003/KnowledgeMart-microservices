package utils

import (
	"math/rand"
	"time"
)

func GenerateOTP() (uint64, time.Time) {
	otp := uint64(rand.Intn(900000) + 100000)
	otpExpiry := time.Now().Add(3 * time.Minute)
	return otp, otpExpiry
}

func GenerateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	randomString := make([]byte, length)

	rand.Seed(time.Now().UnixNano())

	for i := range randomString {
		randomString[i] = charset[rand.Intn(len(charset))]
	}

	return string(randomString)
}
