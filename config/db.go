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

func DeleteInactiveUsers(UserDB *gorm.DB) (int64,error) {
	cutoff := time.Now().Add(-1 * time.Hour)

	tx := UserDB.Begin()

	// Delete the users
	result := tx.Where("activated = ? AND created_at < ?", false, cutoff).Unscoped().Delete(&models.User{})
	if result.Error != nil {
		return 0,result.Error
	}

	err := tx.Commit().Error
	if err != nil {
		fmt.Println(err)
	    tx.Rollback()
	    return 0, err
	}

	return result.RowsAffected,nil
}
