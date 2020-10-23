package controllers

import (
	"dutrozkladbot/config"
	"dutrozkladbot/helpers"
	"dutrozkladbot/keyboard"
	"dutrozkladbot/model"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// ProcessQueryCallback ..
func ProcessQueryCallback(update tgbotapi.Update) {

	// Show typing status
	chataction := tgbotapi.NewChatAction(update.CallbackQuery.Message.Chat.ID, tgbotapi.ChatTyping)
	config.Bot.Send(chataction)

	var arg int
	var err error
	var user model.User

	data := strings.Split(update.CallbackQuery.Data, "/")
	if len(data) > 1 {
		arg, err = strconv.Atoi(data[1])
	}

	userID := update.CallbackQuery.From.ID
	today := time.Now().Format("02.01.2006")

	switch data[0] {
	case "faculties":
		{
			msg := tgbotapi.MessageConfig{}
			msg.ChatID = update.CallbackQuery.Message.Chat.ID
			msg.Text = "факультет:"
			msg.ReplyMarkup = keyboard.GetFacultiesMenu()

			config.Bot.Send(msg)
			return
		}

	case "menu":
		{
			helpers.SendDefaultMessage(update.CallbackQuery.Message.Chat.ID)
			return
		}

	case "today":
		{
			if err := config.DB.Select("Group", "LastQueriedDate").Where("UserID = ?", userID).First(&user).Error; err != nil {
				log.Println(err)
			}

			if user.LastQueriedDate == "" {
				id := int64(userID)
				config.Bot.Send(tgbotapi.NewMessage(id, config.MissingData))
				return
			}

			msg := helpers.GetTimeTable(user.Group, today, userID)
			msg.ChatID = update.CallbackQuery.Message.Chat.ID

			config.Bot.Send(msg)

			if err := config.DB.Exec("UPDATE users SET LastQueriedDate = ? WHERE UserID = ?", today, userID).Error; err != nil {
				log.Println(err)
			}
		}
	case "tomorrow":
		{
			if err := config.DB.Select("Group", "LastQueriedDate").Where("UserID = ?", userID).First(&user).Error; err != nil {
				log.Println(err)
			}

			if user.LastQueriedDate == "" {
				id := int64(userID)
				config.Bot.Send(tgbotapi.NewMessage(id, config.MissingData))
				return
			}

			tomorrowStr := time.Now().AddDate(0, 0, 1).Format("02.01.2006")

			msg := helpers.GetTimeTable(user.Group, tomorrowStr, userID)
			msg.ChatID = update.CallbackQuery.Message.Chat.ID

			config.Bot.Send(msg)

			if err := config.DB.Exec("UPDATE users SET LastQueriedDate = ? WHERE UserID = ?", tomorrowStr, userID).Error; err != nil {
				log.Println(err)
			}
		}

	case "nextday":
		{
			if err := config.DB.Select("Group", "LastQueriedDate").Where("UserID = ?", userID).First(&user).Error; err != nil {
				log.Println(err)
			}

			if user.LastQueriedDate == "" {
				id := int64(userID)
				config.Bot.Send(tgbotapi.NewMessage(id, config.MissingData))
				return
			}

			nextDay, _ := time.Parse("02.01.2006", user.LastQueriedDate)
			nextDayStr := nextDay.AddDate(0, 0, 1).Format("02.01.2006")

			msg := helpers.GetTimeTable(user.Group, nextDayStr, userID)
			msg.ChatID = update.CallbackQuery.Message.Chat.ID

			config.Bot.Send(msg)

			if err := config.DB.Exec("UPDATE users SET LastQueriedDate = ? WHERE UserID = ?", nextDayStr, userID).Error; err != nil {
				log.Println(err)
			}
		}
	case "recheck":
		{
			if err := config.DB.Select("Group", "LastQueriedDate").Where("UserID = ?", userID).First(&user).Error; err != nil {
				log.Println(err)
			}

			if user.LastQueriedDate == "" {
				id := int64(userID)
				config.Bot.Send(tgbotapi.NewMessage(id, config.MissingData))
				return
			}

			msg := helpers.GetTimeTable(user.Group, user.LastQueriedDate, userID)
			msg.ChatID = update.CallbackQuery.Message.Chat.ID

			config.Bot.Send(msg)
		}
	case "faculty":
		{
			if err != nil {
				return
			}

			msg := tgbotapi.MessageConfig{}
			msg.ChatID = update.CallbackQuery.Message.Chat.ID
			msg.Text = "курс:"
			msg.ReplyMarkup = keyboard.GetCoursesMenu(arg)

			config.Bot.Send(msg)

			if err := config.DB.Where("UserID = ?", userID).First(&user).Error; err != nil {
				user.UserID = userID
				user.Faculty = arg

				if err := config.DB.Create(&user).Error; err != nil {
					log.Println(err)
				}
			} else {
				if err := config.DB.Exec("UPDATE users SET Faculty = ? WHERE UserID = ?", arg, userID).Error; err != nil {
					log.Println(err)
				}
			}
		}
	case "course":
		{
			if err != nil {
				return
			}

			if err := config.DB.Select("Faculty").Where("UserID = ?", userID).First(&user).Error; err != nil {
				log.Println(err)
			}

			msg := tgbotapi.MessageConfig{}
			msg.ChatID = update.CallbackQuery.Message.Chat.ID
			msg.Text = "Група:"
			msg.ReplyMarkup = keyboard.GetGroupsMenu(user.Faculty, arg)

			config.Bot.Send(msg)

			if err := config.DB.Exec("UPDATE users SET Course = ? WHERE UserID = ?", arg, userID).Error; err != nil {
				log.Println(err)
			}
		}
	case "group":
		if err != nil {
			return
		}

		msg := helpers.GetTimeTable(arg, today, userID)
		msg.ChatID = update.CallbackQuery.Message.Chat.ID

		config.Bot.Send(msg)

		if err := config.DB.Exec("UPDATE users SET `Group` = ?, `GroupName` = ?, `LastQueriedDate` = ? WHERE UserID = ?", arg, data[2], today, userID).Error; err != nil {
			log.Println(err)
		}
	}

	if config.Bot.Debug {
		fmt.Println("---> "+update.CallbackQuery.Data, "UserID:", update.CallbackQuery.From.ID)
	}
}
