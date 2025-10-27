package utils

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
    "regexp"
    "strings"
)

func IsValidPhoneNumber(phone string) bool {
    phone = strings.TrimSpace(phone)

    nigerianIntl := regexp.MustCompile(`^\+234[0-9]{10}$`)
    nigerianLocal := regexp.MustCompile(`^0[0-9]{10}$`)

    return nigerianIntl.MatchString(phone) || nigerianLocal.MatchString(phone)
}

// registration verification code
func GenerateVerificationCode() string {
	// create byte
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		log.Fatal("Failed to generate verification code: ", err)
	}
	return hex.EncodeToString(bytes)
}

// 2FA code
func Generate2FACode() string {
	// random 6 digit code 0 - 999999
	num, err := rand.Int(rand.Reader, big.NewInt(1000000))
	if err != nil {
		log.Fatal("Failed to generate OTP code: ", err)
	}
	return fmt.Sprintf("%06d", num.Int64())
}
