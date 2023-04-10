package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
)

type LoginResponse struct {
	message string
}

type CreateUserRequestData struct {
	User models.User `json:"user"`
	Key  string      `json:"key"`
}

// Handles Error logs in the console and the response
func HandleError(err error, message string, res http.ResponseWriter) {
	fmt.Println(message)
	fmt.Println(err)
	http.Error(res, fmt.Sprintf("%v: %v", message, err), http.StatusInternalServerError)
}

// Encrypts user data
func EncryptUserData(user *models.User, hashSecret string) (models.User, error) {

	var encryptedUser models.User

	encryptedName, err := services.Encrypt(user.Name, hashSecret)
	if err != nil {
		return encryptedUser, fmt.Errorf("failed to encrypt user's name")
	}

	encryptedEmail, err := services.Encrypt(user.Email, hashSecret)
	if err != nil {
		return encryptedUser, fmt.Errorf("failed to encrypt user's Email")
	}

	encryptedPhone, err := services.Encrypt(user.Phone, hashSecret)
	if err != nil {
		return encryptedUser, fmt.Errorf("failed to encrypt user's Number")
	}

	encryptedTotpSecret, err := services.Encrypt(user.TotpSecret, hashSecret)
	if err != nil {
		return encryptedUser, fmt.Errorf("failed to encrypt user's TotpSecret")
	}

	encryptedUser.Name = encryptedName
	encryptedUser.Email = encryptedEmail
	encryptedUser.Phone = encryptedPhone
	encryptedUser.TotpSecret = encryptedTotpSecret

	return encryptedUser, nil
}

// Handles requests to /login route
func LoginHandler(res http.ResponseWriter, req *http.Request) {
	var user models.User
	var (
		TotpCode    string
		PhoneNumber string
	)

	var DB_URI = os.Getenv("DB_URI")

	err := json.NewDecoder(req.Body).Decode(&user)
	if err != nil {
		HandleError(err, "Failed to decode request body", res)
		return
	}

	resp, err := http.Get(fmt.Sprintf("%v/user/%v", DB_URI, services.Hash(user.Phone)))
	if err != nil {
		HandleError(err, "Failed to make request to DB", res)
		return
	}

	if resp.StatusCode == 500 {
		TotpCode, PhoneNumber, err = CreateUser(&user, DB_URI, res)
		if err != nil {
			HandleError(err, "Failed to create user", res)
			return
		}
	} else {
		err = json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			HandleError(err, "Failed to decode db response body", res)
			return
		}

		resp, err = http.Get(fmt.Sprintf("%v/encryption/%d", DB_URI, user.ID))
		if err != nil {
			HandleError(err, "Failed to get encryption data from DB service", res)
			return
		}

		var encryption models.EncryptionKey

		err = json.NewDecoder(resp.Body).Decode(&encryption)
		if err != nil {
			HandleError(err, "Failed to decode encryption db response", res)
			return
		}

		otpSecret, err := services.Decrypt(user.TotpSecret, encryption.Key)
		if err != nil {
			HandleError(err, "Failed to decrypt totp secret", res)
			return
		}

		TotpCode, err = services.GenerateTotpCode(otpSecret)
		if err != nil {
			HandleError(err, "Failed to generate totp code", res)
			return
		}

		PhoneNumber, err = services.Decrypt(user.Phone, encryption.Key)
		if err != nil {
			HandleError(err, "Failed to decrypt phone number", res)
			return
		}
	}

	// Send Otp code to cliend via sms
	if err := services.NewTwilioService().SendMessage(PhoneNumber, TotpCode); err != nil {
		HandleError(err, "Failed to send message", res)
		return
	}

	response := LoginResponse{
		message: "Login successful",
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(response.message))

}

func CreateUser(user *models.User, DB_URI string, res http.ResponseWriter) (string, string, error) {
	// Generate a random secret for TOTP and encryption
	secret, err := services.GenerateSecret()
	if err != nil {
		return "", "", err
	}

	// Generate a TOTP key and secret using the user's phone number and the random secret
	OtpSecret, err := services.NewOtpService().GenerateTotp(secret, user.Phone)
	if err != nil {
		return "", "", err
	}

	// Generate a random secret for encrypting the user data
	hashSecret, err := services.GenerateSecret()
	if err != nil {
		return "", "", err
	}

	// Update the user struct with the TOTP key and secret
	totpCode, err := services.GenerateTotpCode(OtpSecret)
	if err != nil {
		return "", "", err
	}

	// Add OtpSecret to user model
	user.TotpSecret = OtpSecret

	// Encrypt the user data using the hash secret
	EncryptedUser, err := EncryptUserData(user, hashSecret)
	if err != nil {
		return "", "", err
	}

	EncryptedUser.PhoneHash = services.Hash(user.Phone)

	// Defining the user data to be sent to the database
	userData := CreateUserRequestData{
		User: EncryptedUser,
		Key:  hashSecret,
	}

	// Transforming the data into json
	requestBody, err := json.Marshal(userData)
	if err != nil {
		return "", "", err
	}

	// Creating a buffer to hold the user data
	requestBodyBuffer := bytes.NewBuffer(requestBody)

	// Defining the content type header
	contentType := "application/json"

	// Defining the request client
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Making a post request to the db service to store the user data
	resp, err := client.Post(fmt.Sprintf("%v/user", DB_URI), contentType, requestBodyBuffer)
	if err != nil {
		return "", "", err
	}
	defer resp.Body.Close()

	// Handling the responce from the DB service
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		HandleError(err, "Failed to decode create user response", res)
		return "", "", err
	}

	fmt.Print("PHone number: ", user.Phone)
	fmt.Println("Response body:", string(responseBody))

	return totpCode, user.Phone, nil
}
