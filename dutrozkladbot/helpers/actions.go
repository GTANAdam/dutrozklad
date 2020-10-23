// Package helpers ..
package helpers

import (
	"dutrozkladbot/config"
	"dutrozkladbot/keyboard"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// SendDefaultMessage ..
func SendDefaultMessage(chatID int64) {
	msg := tgbotapi.NewMessage(chatID, config.DefaultMessage())
	msg.ReplyMarkup = keyboard.MainMenu

	config.Bot.Send(msg)
}
