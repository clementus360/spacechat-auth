package controllers

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
	"github.com/gorilla/mux"
)

func AuthorizeClient(res http.ResponseWriter, req *http.Request) {
	tokenString := req.Header.Get("Authorization")
	vars := mux.Vars(req)
	phoneNumber := vars["userId"]

	token := strings.TrimPrefix(tokenString, "Bearer ")

	var DB_URI = os.Getenv("DB_URI")

	encryption, err := GetUserData(DB_URI, phoneNumber)
	if err != nil {
		HandleError(err, "Failed to get user and encryption data", res)
		return
	}

	userId, err := services.ValidateJWTToken(token, encryption.Key)
	if err != nil {
		HandleError(err, "Failed to validate token", res)
		return
	}

	if userId != services.Hash(phoneNumber) {
		HandleError(fmt.Errorf("unauthorized token"), "Failed to match token to user", res)
		return
	}

	sessionId, err := services.GenerateSessionToken()
	if err != nil {
		HandleError(err, "Failed to generate session token", res)
		return
	}

	fmt.Println(sessionId)

	err = services.StoreSession(userId, sessionId)
	if err != nil {
		HandleError(err, "failed to store session", res)
		return
	}

	res.WriteHeader(200)
	res.Write([]byte(sessionId))
}

func AuthorizeServer(res http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	userId := vars["userId"]

	body, err := io.ReadAll(req.Body)
	if err != nil {
		HandleError(err, "Failed to read request body", res)
		return
	}

	ticket := string(body)

	sessionId, err := services.GetSession(userId)
	if err != nil {
		HandleError(err, "Failed to get session token", res)
		return
	}

	fmt.Println(sessionId, "=?", ticket)

	if sessionId != ticket {
		HandleError(fmt.Errorf("ticket does not match"), "Ticket does not match", res)
		return
	}

	fmt.Println("Tokens match")
	res.Write([]byte("Tokens match"))

}

func GetUserData(DB_URI string, userId string) (models.EncryptionKey, error) {

	var user models.User
	var encryption models.EncryptionKey

	resp, err := http.Get(fmt.Sprintf("%v/user/%v", DB_URI, services.Hash(userId)))
	if err != nil {
		return models.EncryptionKey{}, err
	}

	if resp.StatusCode == 500 {
		return models.EncryptionKey{}, fmt.Errorf("failed to find user")
	} else {
		err := json.NewDecoder(resp.Body).Decode(&user)
		if err != nil {
			return models.EncryptionKey{}, err
		}

		resp, err = http.Get(fmt.Sprintf("%v/encryption/%d", DB_URI, user.ID))
		if err != nil {
			return models.EncryptionKey{}, err
		}

		err = json.NewDecoder(resp.Body).Decode(&encryption)
		if err != nil {
			return models.EncryptionKey{}, err
		}
	}

	return encryption, nil
}
