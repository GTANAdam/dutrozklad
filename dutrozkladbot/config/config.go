// Package config ..
package config

import (
	"dutrozkladbot/model"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
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

// Init ..
func Init() {
	log.Println("Database init..")

	var err error
	DB, err = gorm.Open(sqlite.Open("data/database.db"), &gorm.Config{
		//Logger: logger.Default.LogMode(logger.Info)
	})
	if err != nil {
		panic("Failed to open the SQLite database.")
	}

	// Migrate the schema
	DB.Debug().AutoMigrate(&model.User{})
}
