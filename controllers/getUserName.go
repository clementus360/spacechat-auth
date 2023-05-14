package controllers

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
	"github.com/gorilla/mux"
)

type UserData struct {
	Username string `json:"username"`
}

func GetUserName(res http.ResponseWriter, req *http.Request) {
	var user models.User
	var encryption models.EncryptionKey

	var DB_URI = os.Getenv("DB_URI")
	vars := mux.Vars(req)

	resp, err := http.Get(fmt.Sprintf("%v/user/%v", DB_URI, services.Hash(vars["id"])))
	if err != nil {
		HandleError(err, "Failed to make request to DB", res)
		return
	}

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

	err = json.NewDecoder(resp.Body).Decode(&encryption)
	if err != nil {
		HandleError(err, "Failed to decode encryption db response", res)
		return
	}

	name, err := services.Decrypt(user.Name, encryption.Key)
	if err != nil {
		HandleError(err, "Failed to decrypt user name", res)
		return
	}

	userNameResponse := UserData{
		Username: name,
	}
	json.NewEncoder(res).Encode(userNameResponse)
}
