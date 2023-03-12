package services

import (
	"github.com/clementus360/spacechat-auth/config"
	"github.com/robfig/cron"
	"gorm.io/gorm"
)

func SetupDeleteInactiveUsers(UserDB *gorm.DB) {
	c := cron.New()
	c.AddFunc("@daily", func() {config.DeleteInactiveUsers(UserDB)})
    c.Start()
}
