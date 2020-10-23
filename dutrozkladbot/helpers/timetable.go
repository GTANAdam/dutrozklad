// Package helpers ..
package helpers

import (
	"dutrozkladbot/config"
	"dutrozkladbot/keyboard"
	"dutrozkladbot/model"
	"dutrozkladbot/query"
	"dutrozkladbot/util/weekday"
	"encoding/json"
	"fmt"
	"log"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// GetTimeTable ..
func GetTimeTable(arg int, date string, userID int) *tgbotapi.MessageConfig {
	var err error
	result := config.NoSchedule

	dateQueried, _ := time.Parse("02.01.2006", date)
	if dateQueried.Weekday().String() != "Sunday" {
		result, err = QueryTimeTable(arg, date)
		if err != nil {
			if err.Error() != "empty timetable" {
				result = err.Error() + "\nPlease notify the maintainer about this error."
				goto RESUME
			}

			result = config.NoSchedule
		}
	}

RESUME:
	var callback string
	// Check if it isn't today's date
	if time.Now().Format("02.01.2006") != date {
		callback = "> Наступний день"
	}

	var user model.User
	if err := config.DB.Model(&user).Select("GroupName").Where("UserID = ?", userID).First(&user).Error; err != nil {
		log.Println(err)
	}

	msg := tgbotapi.MessageConfig{}
	msg.Text = fmt.Sprintf("<b>%s (%s) [%s]</b>:\n%s", weekday.Time(dateQueried).Weekday().String(), date, user.GroupName, result)
	msg.ParseMode = tgbotapi.ModeHTML
	msg.ReplyMarkup = keyboard.TimeTableMenu(callback)

	return &msg
}

// QueryTimeTable ..
func QueryTimeTable(group int, date string) (string, error) {
	data, err := query.GetTimeTable(group, date)
	if err != nil {
		return "", fmt.Errorf("error fetching data, %s", err)
	}

	var timetable []model.TimeTable
	if err := json.Unmarshal(data.Data, &timetable); err != nil {
		log.Println(err)
	}

	if len(timetable) == 0 {
		return "", fmt.Errorf("empty timetable")
	}

	var result string
	for _, t := range timetable {
		result += fmt.Sprintf("<b>%s-%s</b>: [%s] %s\n<b>aуд.</b>: %s\n<b>вик.</b>: %s\n\n", t.Start, t.End, t.Type, strings.ReplaceAll(t.Name, "\n", ""), t.Cabinet, t.Professor)
	}

	return result, nil
}
