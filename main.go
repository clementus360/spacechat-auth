package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clementus360/spacechat-auth/config"
	"github.com/clementus360/spacechat-auth/controllers"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	router := mux.NewRouter()

	router.HandleFunc("/api/login", controllers.LoginHandler).Methods("GET")
	router.HandleFunc("/api/verify", controllers.VerifyHandler).Methods("POST")


	err := http.ListenAndServe(":3000", router)
	if err!=nil {
		fmt.Println("Failed to start server")
		log.Fatal(err)
	} else {
		fmt.Println("Server runnning on port 5000")
	}
}
