package controllers

import (
	"encoding/json"
	"net/http"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
	"gorm.io/gorm"
)

type AuthData struct {
	PhoneNumber string `json:"phoneNumber"`
	OtpCode string `json:"otpCode"`
}

type ResponseData struct {
	Token string `json:"token"`
}

func VerifyHandler(UserDB *gorm.DB) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var AuthData AuthData
		err := json.NewDecoder(req.Body).Decode(&AuthData)
		if err!=nil {
			HandleError(err, "Failed to decode request body", res)
			return
		}

		var user models.User
		result := UserDB.Where("phone_hash = ?", services.Hash(AuthData.PhoneNumber)).First(&user)
		if result.Error != nil {
			HandleError(err, "Failed to get user from DB", res)
		}

		var encryption models.EncryptionKey
		result = UserDB.Where("user_id = ?", user.ID).First(&encryption)
		if result.Error!=nil {
			HandleError(result.Error, "Failed to get encryption key from DB", res)
		}

		TotpSecret,err := services.Decrypt(user.TotpSecret, encryption.Key)
		if err!=nil {
			HandleError(err, "Failed to decrypt totp secret", res)
		}

		valid := services.VerifyTotpCode(AuthData.OtpCode,TotpSecret)

		if !valid {
			http.Error(res, "TotpCode is invalid", http.StatusUnauthorized)
			return
		}

		jwtToken,err := services.GenerateJWTToken(user.PhoneHash, encryption.Key)
		if err!=nil {
			HandleError(err, "Failed to generate token", res)
		}

		ResponseData := ResponseData{
			Token: jwtToken,
		}

		json.NewEncoder(res).Encode(ResponseData)

	}
}
