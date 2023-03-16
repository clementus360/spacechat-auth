package config

import (
	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// Connect to postgres database

func ConnectDatabase() *gorm.DB {
	var (
		host = os.Getenv("HOST")
		user = os.Getenv("USER")
		dbname = os.Getenv("DBNAME")
		sslmode = os.Getenv("SSLMODE")
		password = os.Getenv("PASSWORD")
		dbport = os.Getenv("PORT")
	)

	connStr := fmt.Sprintf("host=%s  user=%s dbname=%s sslmode=%s password=%s port=%s", host, user, dbname, sslmode, password, dbport)

	db,err := gorm.Open(postgres.Open(connStr), &gorm.Config{})

	if err!=nil {
		fmt.Println("Failed to connect to DB")
		panic(err)
	}

	return db
}

func AutoMigrate(db *gorm.DB, model interface{}) {
	err := db.AutoMigrate(model)

	if err!=nil {
		fmt.Printf("Failed to migrate the model")
		log.Fatal(err)
	}
}
