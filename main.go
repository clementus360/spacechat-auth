package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clementus360/spacechat-auth/config"
	"github.com/clementus360/spacechat-auth/controllers"
	"github.com/clementus360/spacechat-auth/models"
	"github.com/clementus360/spacechat-auth/services"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	UserDB := config.ConnectDatabase()


	config.AutoMigrate(UserDB,&models.User{})
	config.AutoMigrate(UserDB, &models.EncryptionKey{})

	router := mux.NewRouter()

	router.HandleFunc("/api/login", controllers.LoginHandler(UserDB)).Methods("GET")
	router.HandleFunc("/api/verify", controllers.VerifyHandler(UserDB)).Methods("POST")

	go func ()  {
		services.DeleteUsersHandler(UserDB)
	}()


	err := http.ListenAndServe(":3000", router)
	if err!=nil {
		fmt.Println("Failed to start server")
		log.Fatal(err)
	} else {
		fmt.Println("Server runnning on port 5000")
	}
}
