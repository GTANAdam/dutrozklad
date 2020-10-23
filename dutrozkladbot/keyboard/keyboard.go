// Package keyboard ..
package keyboard

import (
	"dutrozkladbot/query"
	"dutrozkladbot/util"
	"dutrozkladbot/util/orderedmap"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// MainMenu ..
var MainMenu = tgbotapi.NewInlineKeyboardMarkup(
	tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData("Групи", "faculties"),
		tgbotapi.NewInlineKeyboardButtonData("Сьогодні", "today"),
		tgbotapi.NewInlineKeyboardButtonData("Завтра", "tomorrow"),
	),
)

// TimeTableMenu ..
func TimeTableMenu(key string) *tgbotapi.InlineKeyboardMarkup {
	if key == "" {
		key = "> Завтра"
	}

	result := tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(key, "nextday"),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData("Перевірити", "recheck"),
			tgbotapi.NewInlineKeyboardButtonData("Меню", "menu"),
		),
	)

	return &result
}

// GetFacultiesMenu ..
func GetFacultiesMenu() *tgbotapi.InlineKeyboardMarkup {
	response, err := query.GetFaculties()
	if err != nil {
		return nil
	}

	// Deserialize
	var data orderedmap.OrderedMap
	if err := json.Unmarshal(response.Data, &data); err != nil {
		log.Println(err)
	}

	// Filter out the additional text
	for i, k := range data.Map {
		str := strings.Replace(k.(string), "Навчально-науковий ", "", 1)
		data.Map[i] = util.UpcaseInitial(str)
	}

	return createKeyboard(&data, false)
}

// GetCoursesMenu ..
func GetCoursesMenu(faculty int) *tgbotapi.InlineKeyboardMarkup {
	response, err := query.GetCourses(faculty)
	if err != nil {
		return nil
	}

	var data []string
	if err := json.Unmarshal(response.Data, &data); err != nil {
		log.Println(err)
	}

	return createKeyboard(data, false)
}

// GetGroupsMenu ..
func GetGroupsMenu(faculty, course int) *tgbotapi.InlineKeyboardMarkup {
	response, err := query.GetGroups(faculty, course)
	if err != nil {
		return nil
	}

	var data orderedmap.OrderedMap
	if err := json.Unmarshal(response.Data, &data); err != nil {
		log.Println(err)
	}

	return createKeyboard(&data, true)
}

func createKeyboard(inter interface{}, groups bool) *tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	if arr, ok := inter.([]string); ok {
		// Courses
		for i := 0; i <= len(arr); i += 3 {
			row := make([]tgbotapi.InlineKeyboardButton, 0, 3)

			pickPos := i + 3
			if pickPos > len(arr) {
				pickPos = len(arr)
			}

			for _, d := range arr[i:pickPos] {
				callback := fmt.Sprintf("course/%s", d)
				button := tgbotapi.InlineKeyboardButton{Text: d, CallbackData: &callback}
				row = append(row, button)
			}

			rows = append(rows, row)
		}
	} else if data, ok := inter.(*orderedmap.OrderedMap); ok {
		if groups {
			// Groups
			for i := 0; i <= len(data.Order); i += 3 {
				row := make([]tgbotapi.InlineKeyboardButton, 0, 3)

				pickPos := i + 3
				if pickPos > len(data.Map) {
					pickPos = len(data.Map)
				}

				for _, d := range data.Order[i:pickPos] {
					txt := fmt.Sprint(data.Map[d])
					callback := fmt.Sprintf("group/%s/%s", util.Mapkey(data.Map, txt), txt)

					button := tgbotapi.InlineKeyboardButton{Text: txt, CallbackData: &callback}
					row = append(row, button)
				}

				rows = append(rows, row)
			}
		} else {
			// Faculties
			for _, key := range data.Order {
				button := tgbotapi.InlineKeyboardButton{}

				keyStr := fmt.Sprintf("faculty/%v", key)
				button.Text = data.Map[key].(string)
				button.CallbackData = &keyStr

				rows = append(rows, tgbotapi.NewInlineKeyboardRow(button))
			}
		}
	} else {
		panic("type not supported.")
	}

	result := tgbotapi.NewInlineKeyboardMarkup(rows...)
	return &result
}
