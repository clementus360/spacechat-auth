package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
	"gorm.io/gorm"
)

type LoginResponse struct {
	message string
}

// Handles Error logs in the console and the response
func HandleError(err error,message string, res http.ResponseWriter) {
		fmt.Println(message)
		fmt.Println(err)
		http.Error(res, message, http.StatusInternalServerError)
}

// Encrypts user data
func EncryptUserData(user *models.User, hashSecret string) (models.User, error) {

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


	encryptedTotpSecret,err := services.Encrypt(user.TotpSecret,hashSecret)
	if err!=nil {
		return encryptedUser,fmt.Errorf("failed to encrypt user's TotpSecret")
	}

	encryptedUser.Name = encryptedName
	encryptedUser.Email = encryptedEmail
	encryptedUser.Phone = encryptedPhone
	encryptedUser.TotpSecret = encryptedTotpSecret

	return encryptedUser, nil
}


// Handles requests to /login route
func LoginHandler(UserDB *gorm.DB) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request)  {
		var user models.User
		var (
			TotpCode string
			PhoneNumber string
		)
		err := json.NewDecoder(req.Body).Decode(&user)
		if err!=nil {
			HandleError(err, "Failed to decode request body", res)
			return
		}

		result := UserDB.Where("phone_hash = ?", services.Hash(user.Phone)).First(&user)
		if result.Error != nil {
			if errors.Is(result.Error, gorm.ErrRecordNotFound) {
				TotpCode,PhoneNumber,err = CreateUser(&user,UserDB,res)
				if err!=nil {
					HandleError(err, "Failed to create user", res)
					return
				}
			}
		} else {
			var encryption models.EncryptionKey
			if result := UserDB.Where("user_id = ?", user.ID).First(&encryption); result.Error!=nil {
				HandleError(result.Error,"Failed to get encryption key from DB", res)
				return
			} else {
				OtpSecret,err := services.Decrypt(user.TotpSecret, encryption.Key)
				if err!=nil {
					HandleError(err, "Failed to decrypt totp code", res)
					return
				}
				phoneNumber,err := services.Decrypt(user.Phone,encryption.Key)
				if err!=nil {
					HandleError(err, "Failed to decrypt phone number", res)
				}
				TotpCode,err = services.GenerateTotpCode(OtpSecret)
				if err!=nil {
					HandleError(err, "Failed to generate otp code", res)
				}

				PhoneNumber = phoneNumber
			}
		}

		// Send Otp code to cliend via sms
		if err:= services.NewTwilioService().SendMessage(PhoneNumber, TotpCode); err!=nil {
			HandleError(err, "Failed to send message", res)
			return
		}

		response := LoginResponse{
			message: "Login successful",
		}

		res.WriteHeader(http.StatusOK)
		res.Write([]byte(response.message))

	}
}

func CreateUser(user *models.User, UserDB *gorm.DB, res http.ResponseWriter) (string, string, error){
	// Generate a random secret for TOTP and encryption
	secret,err := services.GenerateSecret()
	if err!=nil {
		HandleError(err, "Failed to generate secret", res)
		return "","",err
	}

	// Generate a TOTP key and secret using the user's phone number and the random secret
	OtpSecret,err := services.NewOtpService().GenerateTotp(secret,user.Phone)
	if err!=nil {
		HandleError(err, "Failed to generate otp", res)
		return "","",err
	}

	// Generate a random secret for encrypting the user data
	hashSecret,err := services.GenerateSecret()
	if err!=nil {
		HandleError(err, "Failed to generate hash secret", res)
		return "","",err
	}

	// Update the user struct with the TOTP key and secret
	totpCode,err := services.GenerateTotpCode(OtpSecret)
	if err!=nil {
		HandleError(err, "Failed to generate otp code", res)
	}

	// Add OtpSecret to user model
	user.TotpSecret = OtpSecret

	// Encrypt the user data using the hash secret
	EncryptedUser,err := EncryptUserData(user,hashSecret)
	if err!=nil {
		HandleError(err, "Failed to encrypt user", res)
		return "","",err
	}

	EncryptedUser.PhoneHash = services.Hash(user.Phone)
	if err:= UserDB.Create(&EncryptedUser).Error; err!=nil {
		HandleError(err, "Failed to add user to DB", res)
		return "","",err
	}

	// Save the hash secret to the database
	EncryptionKey := models.EncryptionKey{
		UserID: EncryptedUser.ID,
		Key: hashSecret,
	}

	if err:= UserDB.Create(&EncryptionKey).Error; err!=nil {
		HandleError(err, "Failed to add Encryption keys to DB", res)
		return "","",err
	}

	return totpCode,user.Phone,nil
}
