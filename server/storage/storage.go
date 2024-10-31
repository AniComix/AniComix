package storage

import (
	"github.com/AniComix/server/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

var (
	dataDir = "/"
	db      *gorm.DB
)

func InitStorage() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir = filepath.Join(homeDir, ".AniComix")
	initDb()
}

func initDb() {
	db, err := gorm.Open(sqlite.Open(filepath.Join(dataDir, "main.db")), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	_ = db.AutoMigrate(&models.User{})
}

func DB() *gorm.DB {
	return db
}

func DataDir() string {
	return dataDir
}
