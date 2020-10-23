package teachers

import (
	"log"
	"net/url"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/models"
	"dutrozkladapi/util"

	"github.com/PuerkitoBio/goquery"
)

func GetTeacherTimeTable(kafedra string, teacher string, date1 string, date2 string) []models.TimeTable {
	defer util.TimeTrack(time.Now(), "TeacherTimeTable")

	formData := url.Values{
		"TimeTableForm[kafedra]": {kafedra},
		"TimeTableForm[teacher]": {teacher},
		"TimeTableForm[date1]":   {date1},
		"TimeTableForm[date2]":   {date2},
		"TimeTableForm[r11]":     {"5"},
		"timeTable":              {"1"},
	}

	result, _ := util.Post(header.TeacherURL, formData)
	if result == nil {
		log.Println("[ERROR] Connection timeout.")
		return nil
	}

	defer result.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(result)
	if err != nil {
		log.Println(err)
	}

	var data [][]string
	// Iterate through found script nodes and break when data is extracted & unmarshal
	doc.Find("script").EachWithBreak(func(i int, s *goquery.Selection) bool {
		if data, err = util.TimeTableUnmarshal(s.Text()); data != nil {
			if err != nil { // In case of error
				log.Println(err)
			}
			return false // break
		}
		return true // continue
	})

	if len(data) == 0 {
		log.Println("[INFO] TeacherTimeTable found no data.")
		return nil
	}

	var toReturn []models.TimeTable
	for _, value := range data {
		timetable := models.TimeTable{
			Type:      value[header.Type],
			Name:      value[header.Name],
			Professor: value[header.Prof],
			Cabinet:   value[header.Cabinet],
			Date:      value[header.Date],
			Start:     value[header.Start],
			End:       value[header.End],
			Misc:      value[header.Misc],
		}

		toReturn = append(toReturn, timetable)
	}

	return toReturn
}
