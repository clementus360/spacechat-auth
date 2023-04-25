package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
)

type AuthData struct {
	PhoneNumber string `json:"phoneNumber"`
	OtpCode     string `json:"otpCode"`
}

type VerifyResponse struct {
	Token string `json:"token"`
}

func VerifyHandler(res http.ResponseWriter, req *http.Request) {
	var AuthData AuthData
	err := json.NewDecoder(req.Body).Decode(&AuthData)
	if err != nil {
		HandleError(err, "Failed to decode request body", res)
		return
	}

	var DB_URI = os.Getenv("DB_URI")

	var user models.User
	var encryption models.EncryptionKey

	resp, err := http.Get(fmt.Sprintf("%v/user/%v", DB_URI, services.Hash(AuthData.PhoneNumber)))
	if err != nil {
		HandleError(err, "Failed to get user from DB", res)
		return
	}

	if resp.StatusCode == 500 {
		http.Error(res, "User does not exist", http.StatusBadRequest)
		return
	} else {
		err := json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			HandleError(err, "Failed to decode user", res)
			return
		}

		resp, err = http.Get(fmt.Sprintf("%v/encryption/%d", DB_URI, user.ID))
		if err != nil {
			HandleError(err, "Failed to get encryption data from DB service", res)
			return
		}

		err = json.NewDecoder(resp.Body).Decode(&encryption)
		if err != nil {
			HandleError(err, "Failed to decode user", res)
			return
		}
	}

	TotpSecret, err := services.Decrypt(user.TotpSecret, encryption.Key)
	if err != nil {
		HandleError(err, "Failed to decrypt totp secret", res)
		return
	}

	valid := services.VerifyTotpCode(AuthData.OtpCode, TotpSecret)

	if !valid {
		http.Error(res, "TotpCode is invalid", http.StatusUnauthorized)
		return
	}

	jwtToken, err := services.GenerateJWTToken(user.PhoneHash, encryption.Key)
	if err != nil {
		HandleError(err, "Failed to generate token", res)
		return
	}

	ResponseData := VerifyResponse{
		Token: jwtToken,
	}

	request, err := http.NewRequest("PUT", fmt.Sprintf("%v/user/%v", DB_URI, user.PhoneHash), nil)
	if err != nil {
		HandleError(err, "Failed to create update request", res)
		return
	}

	client := &http.Client{}

	resp, err = client.Do(request)
	if err != nil {
		HandleError(err, "Failed to activate user", res)
		return
	}

	if resp.StatusCode != http.StatusOK {
		HandleError(fmt.Errorf("failed to activate user"), "Failed to activate user", res)
		return
	}

	json.NewEncoder(res).Encode(ResponseData)

}
