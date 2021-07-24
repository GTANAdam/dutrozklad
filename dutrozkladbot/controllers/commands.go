// Package controllers ..
package controllers

import (
	"dutrozkladbot/config"
	"dutrozkladbot/helpers"
	"dutrozkladbot/model"
	"dutrozkladbot/query"
	"dutrozkladbot/util"
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ProcessCommand ..
func ProcessCommand(update *tgbotapi.Update, cmd string) {
	switch cmd {
	case "status":
		{
			var count int64
			config.DB.Model(model.User{}).Count(&count)

			config.MutexStats.RLock()
			callbacks := config.Stats.Callbacks
			messages := config.Stats.Messages
			queries := config.Stats.Queries
			config.MutexStats.RUnlock()

			stats, err := query.GetStats()
			if err != nil {
				log.Println(err)
				return
			}

			msg := tgbotapi.MessageConfig{}
			msg.ChatID = update.Message.Chat.ID
			msg.Text = fmt.Sprintf("*Total users*: %v\n\n*Total callbacks received*: %v\n*Total messages received*: %v\n\n*Total clt queries invoked*: %v\n*Total srv cached queries*: %v\n*Total srv uncached queries*: %v\n\n*Uptime*: %s", count, callbacks, messages, queries, stats.Cached, stats.Uncached, util.FormatSince(config.StartedAt))
			msg.ParseMode = tgbotapi.ModeMarkdown

			config.Bot.Send(msg)
		}
	default:
		helpers.SendDefaultMessage(update.Message.Chat.ID)
	}
}
