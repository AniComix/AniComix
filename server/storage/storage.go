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
	dataDir  = "/"
	cacheDir = "/"
	db       *gorm.DB
)

func InitStorage() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	dataDir = filepath.Join(homeDir, ".AniComix")
	cacheDir, err = os.UserCacheDir()
	if err != nil {
		cacheDir = filepath.Join(homeDir, ".cache")
	} else {
		cacheDir = filepath.Join(cacheDir, ".AniComix")
	}
	if _, err := os.Stat(cacheDir); os.IsExist(err) {
		// clear cache
		err = os.RemoveAll(cacheDir)
		if err != nil {
			log.Fatal(err)
		}
		err = os.MkdirAll(cacheDir, 0755)
		if err != nil {
			log.Fatal(err)
		}
	}
	initDb()
}

func initDb() {
	os.Mkdir(dataDir, os.ModePerm)
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

func CacheDir() string {
	return cacheDir
}
