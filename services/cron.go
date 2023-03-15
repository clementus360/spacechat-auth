package services

import (
	"fmt"
	"log"

	"github.com/clementus360/spacechat-auth/config"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func DeleteUsersHandler(UserDB *gorm.DB) {
	fmt.Println("Test cron")
	c := cron.New()
	c.AddFunc("@hourly", func() {
		RowsAffected, err := config.DeleteInactiveUsers(UserDB)

		if err!=nil {
			log.Println(err)
		} else {
			log.Printf("Deleted %v users", RowsAffected)
		}
	})

	c.Start()

	select {}
}
