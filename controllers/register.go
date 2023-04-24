package controllers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
)

func RegisterHandler(res http.ResponseWriter, req *http.Request) {
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

	TotpCode, PhoneNumber, err = CreateUser(&user, DB_URI, res)
	if err != nil {
		HandleError(err, "Failed to create user", res)
		return
	}

	// Send Otp code to cliend via sms
	if err := services.NewTwilioService().SendMessage(PhoneNumber, TotpCode); err != nil {
		HandleError(err, "Failed to send message", res)
		return
	}

	response := LoginResponse{
		message: "Register successful",
	}

	res.WriteHeader(http.StatusOK)
	res.Write([]byte(response.message))
}
