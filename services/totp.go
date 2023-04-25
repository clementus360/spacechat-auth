package services

import (
	"fmt"
	"os"
	"time"

	"github.com/pquerna/otp"
	"github.com/pquerna/otp/totp"
)

type TOTPService interface {
	GenerateCode(Secret string) (string, error)
}

type OTPService struct {
	Issuer string
}

func NewOtpService() *OTPService {
	return &OTPService{
		Issuer: os.Getenv("ISSUER"),
	}
}

func (os OTPService) GenerateTotp(secret string, number string) (string, error) {
	totpConfig := totp.GenerateOpts{
		Issuer:      os.Issuer,
		AccountName: number,
		Period:      120,
		Digits:      6,
		Algorithm:   otp.AlgorithmSHA256,
	}

	totpKey, err := totp.Generate(totpConfig)
	if err != nil {
		fmt.Println("Failed to generate otp key")
		return "", err
	}

	if err != nil {
		fmt.Println("Failed to generate otp code")
		return "", err
	}

	return totpKey.Secret(), nil
}

func GenerateTotpCode(secret string) (string, error) {
	totpCode, err := totp.GenerateCode(secret, time.Now().UTC())
	if err != nil {
		return "", err
	}

	return totpCode, nil
}

func VerifyTotpCode(totpCode string, totpSecret string) bool {
	isValid := totp.Validate(totpCode, totpSecret)
	return isValid
}
