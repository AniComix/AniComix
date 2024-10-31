package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	Username     string `gorm:"unique"`
	Nickname     string
	PasswordHash string
	AvatarPath   string // relative path to $dataDir/avatars
	IsAdmin      bool
	bio          string
}
