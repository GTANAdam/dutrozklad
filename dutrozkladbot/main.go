package main

import (
	"flag"
	"log"
	"runtime"

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

	if runtime.GOOS == "windows" {
		config.Token = config.Config.Token.Development
		config.APIAddress = config.Config.Token.Development
	} else {
		config.Token = config.Config.Token.Production
		config.APIAddress = config.Config.Token.Production
	}

	log.Println("DUTRozkladBOT init..")

	debug := flag.Bool("debug", false, "Enables debug mode, shows more verbosity")
	flag.Parse()

	// util.LoadUsersFromJSON()

	// fmt.Println("Saving users..")
	// util.SaveUsersToDatabase()
	// return

	var count int64
	config.DB.Model(model.User{}).Count(&count)
	log.Printf("Loaded %v users from database.\n", count)

	bot, err := tgbotapi.NewBotAPI(config.Token)
	if err != nil {
		log.Panic(err)
	}

	bot.Debug = *debug
	config.Bot = bot

	log.Printf("Authorized on account: %s\n", bot.Self.UserName)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)
	if err != nil {
		log.Fatalln(err)
	}

	for update := range updates {
		if update.Message != nil {
			go func(*tgbotapi.Update) {
				// if update.Message.Text[0] == '/' {
				// 	controllers.ProcessCommand(update)
				// 	return
				// }
			}(&update)

			go helpers.SendDefaultMessage(update.Message.Chat.ID)
		}

		if update.CallbackQuery != nil {
			go controllers.ProcessQueryCallback(update)
		}
	}
}
