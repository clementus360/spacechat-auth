package models

import "gorm.io/gorm"

type EncryptionKey struct {
	gorm.Model
    UserID    uint
    Key       string
    User    User    `gorm:"constraint:OnDelete:CASCADE"`
}
