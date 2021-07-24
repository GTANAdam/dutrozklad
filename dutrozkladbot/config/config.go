// Package config ..
package config

import (
	"dutrozkladbot/model"
	"log"
	"sync"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Bot ..
var Bot *tgbotapi.BotAPI

// Config ..
var Config = &model.Config{}

// Token ..
var Token string

// APIAddress ..
var APIAddress string

// DB ..
var DB *gorm.DB

// Stats ..
var Stats = &model.Stats{}

// StartedAt ..
var StartedAt = time.Now()

// MutexStats ..
var MutexStats = &sync.RWMutex{}

// Init ..
func Init() {
	log.Println("Database init..")

	var err error
	DB, err = gorm.Open(sqlite.Open("data/database.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info), // Info
	})

	if err != nil {
		panic("Failed to open the SQLite database.")
	}

	ddl := `PRAGMA synchronous = NORMAL;`
	if err := DB.Exec(ddl).Error; err != nil {
		log.Fatal(err)
	}

	// Migrate the schema
	DB.Debug().AutoMigrate(&model.User{})
}
