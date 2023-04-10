package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/clementus360/spacechat-auth/config"
	"github.com/clementus360/spacechat-auth/controllers"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

func main() {
	config.LoadEnv()

	router := mux.NewRouter()

	router.HandleFunc("/api/login", controllers.LoginHandler).Methods("GET")
	router.HandleFunc("/api/verify", controllers.VerifyHandler).Methods("POST")
	router.HandleFunc("/api/authorize/{userId}", controllers.AuthorizeClient).Methods("GET")

	// Create CORS handler with allowed headers and origins
	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	originsOk := handlers.AllowedOrigins([]string{"*"})
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST"})

	// Wrap the router with CORS middleware
	corsHandler := handlers.CORS(headersOk, originsOk, methodsOk)(router)

	err := http.ListenAndServe(":3000", corsHandler)
	if err != nil {
		fmt.Println("Failed to start server")
		log.Fatal(err)
	} else {
		fmt.Println("Server runnning on port 5000")
	}
}
