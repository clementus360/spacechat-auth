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

func NewOtpService() *OTPService{
	return &OTPService{
		Issuer: os.Getenv("ISSUER"),
	}
}

func (os OTPService) GenerateTotp(secret string, number string) (string, string, error) {
	totpConfig := totp.GenerateOpts{
		Issuer: os.Issuer,
		AccountName: number,
		Period: 30,
		Digits: 6,
		Algorithm: otp.AlgorithmSHA256,
	}

	totpKey,err := totp.Generate(totpConfig)
	if err!=nil {
		fmt.Println("Failed to generate otp key")
		return "","",err
	}

	totpCode,err := totp.GenerateCode(totpKey.Secret(), time.Now())

	if err!=nil {
		fmt.Println("Failed to generate otp code")
		return "","",err
	}

	fmt.Println(totpKey)

	return totpCode, totpKey.Secret(), nil
}
