package main

import (
	"flag"
	"log"

	config "dutrozkladbot/config"
	"dutrozkladbot/controllers"
	"dutrozkladbot/helpers"
	"dutrozkladbot/model"
	"dutrozkladbot/util"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	config.Init()

	log.Println("Config init..")
	util.LoadConfig(config.Config)

	log.Println("init program..")

	debug := flag.Bool("debug", false, "Enables debug mode, shows more verbosity")
	flag.Parse()

	if *debug {
		config.Token = config.Config.Token.Development
		config.APIAddress = config.Config.APIAddress.Development
	} else {
		config.Token = config.Config.Token.Production
		config.APIAddress = config.Config.APIAddress.Production
	}

	// util.LoadUsersFromJSON()

	// fmt.Println("Saving users..")
	// util.SaveUsersToDatabase()
	// return

	var count int64
	config.DB.Model(model.User{}).Count(&count)
	log.Printf("Loaded %v users from database.\n", count)

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Println(err)
	}

	bot.Debug = *debug
	config.Bot = bot

	log.Printf("Authorized on account: %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Println(err)
	}

	for update := range updates {
		if update.Message != nil {

			if cmd := update.Message.Command(); cmd != "" {
				go controllers.ProcessCommand(&update, cmd)

			} else {
				config.MutexStats.Lock()
				config.Stats.Messages++
				config.MutexStats.Unlock()

				go helpers.SendDefaultMessage(update.Message.Chat.ID)
			}
		}

		if update.CallbackQuery != nil {
			go controllers.ProcessQueryCallback(&update)
		}
	}
}
