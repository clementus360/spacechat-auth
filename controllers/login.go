package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
	"gorm.io/gorm"
)

// Handles Error logs in the console and the response
func HandleError(err error,message string, res http.ResponseWriter) {
		fmt.Println(message)
		fmt.Println(err)
		http.Error(res, message, http.StatusInternalServerError)
}

// Encrypts user data
func EncryptUserData(user models.User, hashSecret string) (models.User, error) {

	var encryptedUser models.User

	encryptedName,err := services.Encrypt(user.Name,hashSecret)
	if err!=nil {
		return encryptedUser,fmt.Errorf("failed to encrypt user's name")
	}

	encryptedEmail,err := services.Encrypt(user.Email,hashSecret)
	if err!=nil {
		return encryptedUser,fmt.Errorf("failed to encrypt user's Email")
	}

	encryptedPhone,err := services.Encrypt(user.Phone,hashSecret)
	if err!=nil {
		return encryptedUser,fmt.Errorf("failed to encrypt user's Number")
	}

	encryptedTotpCode,err := services.Encrypt(user.TotpCode,hashSecret)
	if err!=nil {
		return encryptedUser,fmt.Errorf("failed to encrypt user's TotpCode")
	}

	encryptedTotpSecret,err := services.Encrypt(user.TotpSecret,hashSecret)
	if err!=nil {
		return encryptedUser,fmt.Errorf("failed to encrypt user's TotpSecret")
	}

	encryptedUser.Name = encryptedName
	encryptedUser.Email = encryptedEmail
	encryptedUser.Phone = encryptedPhone
	encryptedUser.TotpCode = encryptedTotpCode
	encryptedUser.TotpSecret = encryptedTotpSecret

	return encryptedUser, nil

}


// Handles requests to /login route
func LoginHandler(UserDB *gorm.DB) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request)  {
		var user models.User
		err := json.NewDecoder(req.Body).Decode(&user)
		if err!=nil {
			HandleError(err, "Failed to decode request body", res)
			return
		}

		// Generate a random secret for TOTP and encryption
		secret,err := services.GenerateSecret()
		if err!=nil {
			HandleError(err, "Failed to generate secret", res)
			return
		}

		// Generate a TOTP key and secret using the user's phone number and the random secret
		OtpKey,OtpSecret,err := services.NewOtpService().GenerateTotp(secret,user.Phone)
		if err!=nil {
			HandleError(err, "Failed to generate otp", res)
			return
		}

		// Generate a random secret for encrypting the user data
		hashSecret,err := services.GenerateSecret()
		if err!=nil {
			HandleError(err, "Failed to generate hash secret", res)
			return
		}

		// Update the user struct with the TOTP key and secret
		user.TotpCode = OtpKey
		user.TotpSecret = OtpSecret

		// Encrypt the user data using the hash secret
		EncryptedUser,err := EncryptUserData(user,hashSecret)
		if err!=nil {
			HandleError(err, "Failed to encrypt user", res)
			return
		}

		EncryptedUser.PhoneHash = services.Hash(user.Phone)

		if err:= UserDB.Create(&EncryptedUser).Error; err!=nil {
			HandleError(err, "Failed to add user to DB", res)
			return
		}

		// Save the hash secret to the database
		EncryptionKey := models.EncryptionKey{
			UserID: EncryptedUser.ID,
			Key: hashSecret,
		}

		if err:= UserDB.Create(&EncryptionKey).Error; err!=nil {
			HandleError(err, "Failed to add Encryption keys to DB", res)
			return
		}

		// Send Otp code to cliend via sms
		if err:= services.NewTwilioService().SendMessage(user.Phone, user.TotpCode); err!=nil {
			HandleError(err, "Failed to send message", res)
			return
		}

		res.WriteHeader(http.StatusOK)
		json.NewEncoder(res).Encode(user)

	}
}
