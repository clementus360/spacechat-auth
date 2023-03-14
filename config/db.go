package config

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/clementus360/spacechat-auth/models"
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

func DeleteInactiveUsers(UserDB *gorm.DB) {
	// Calculate the time from which we will select users
	// cutoff := time.Now().Add(-24 * time.Hour)
	cutoff := time.Now().Add(-5 * time.Minute)

	fmt.Print("test")


	// Get all users that older than 24hours
	var inactiveUsers []models.User
	err:=UserDB.Where("Activated = ? AND created_at < ?", false, cutoff).Find(&inactiveUsers).Error
	if err!=nil {
		log.Fatal(err)
	}

	fmt.Println(inactiveUsers)
	fmt.Println(cutoff)


	// Delete the inactive users
	for _,user := range inactiveUsers {
		err = UserDB.Delete(&user).Error
		if err!=nil {
			log.Println("Error deleting user", err)
		}
	}
}
