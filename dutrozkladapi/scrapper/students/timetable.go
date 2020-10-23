package students

import (
	"fmt"
	"io"
	"log"
	"net/url"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/models"
	"dutrozkladapi/util"

	"github.com/PuerkitoBio/goquery"
)

// GetTimeTable ..
func GetTimeTable(faculty, course, group, date1, date2 string) ([]*models.TimeTable, error) {
	defer util.TimeTrack(time.Now(), "TimeTable")

	// Return from local database
	if exists, fac, cou := recordExists(group, date1, date2); exists {
		return returnOnlyRequestedDates(header.Faculties[fac].Courses[cou].Groups[group].TimeTable, date1, date2), nil
	}

	// Get from server
	data, err := getAllTimeTables(faculty, course, group, time.Now().Format("02.01.2006"), time.Now().AddDate(1, 0, 0).Format("02.01.2006"))
	if err != nil {
		return nil, err
	}

	return returnOnlyRequestedDates(data, date1, date2), nil
}

func getAllTimeTables(faculty, course, group, date1, date2 string) ([]*models.TimeTable, error) {
	tries := 0

START:
	tries++
	fmt.Println("--> Querying", faculty, course, group)

	formData := url.Values{
		"TimeTableForm[faculty]": {faculty},
		"TimeTableForm[course]":  {course},
		"TimeTableForm[group]":   {group},
		"TimeTableForm[date1]":   {date1},
		"TimeTableForm[date2]":   {date2},
		"TimeTableForm[r11]":     {"5"},
		"timeTable":              {"1"},
	}

	result, err := util.Post(header.URL, formData)
	if err != nil {
		if tries > 2 {
			return nil, err
		}

		log.Println(err, ", retrying in 1 second,", tries, "tries already..")
		time.Sleep(time.Second * 2)
		goto START
	}

	defer result.Close()

	data, err := extractTimeTableFromDoc(&result)
	if err != nil {
		return nil, err
	}

	var resultArr []*models.TimeTable
	for _, value := range data {
		timetable := &models.TimeTable{
			Type:      value[header.Type],
			Name:      value[header.Name],
			Professor: value[header.Prof],
			Cabinet:   value[header.Cabinet],
			Date:      value[header.Date],
			Start:     value[header.Start],
			End:       value[header.End],
			Misc:      value[header.Misc],
		}

		resultArr = append(resultArr, timetable)
	}

	return resultArr, nil
}

func returnOnlyRequestedDates(arr []*models.TimeTable, date1, date2 string) []*models.TimeTable {
	result := make([]*models.TimeTable, 0)

	for i := range arr {
		// API inconsistency: Datetime returned as 2020-10-19 whereas the request is submitted as 19.10.2020
		if dWR, _ := dateWithinRange(arr[i].Date, "2006-01-02", date1, date2, "02.01.2006"); dWR {
			result = append(result, arr[i])
		}
	}

	return result
}

func extractTimeTableFromDoc(result *io.ReadCloser) ([][]string, error) {
	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(*result)
	if err != nil {
		return nil, err
	}

	var data [][]string
	// Fast but may be unreliable, therefor we'll check if it exists just to make sure
	if scriptNode := doc.Find("script"); scriptNode.Length() > header.TimeTableNodeID {
		log.Println("[INFO] TimeTable using fast method")

		// Create a new document from node and unmarshal
		node := goquery.NewDocumentFromNode(scriptNode.Get(header.TimeTableNodeID))
		if data, err = util.TimeTableUnmarshal(node.Text()); err != nil {
			// log.Println(err)
			return nil, err
		}
	} else {
		// Slow but reliable, if nodes length is less than TimeTableNodeID, iterate over all nodes and select regex matched one
		log.Println("[INFO] TimeTable using slow method")

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
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("timetable found no data")
	}

	return data, nil
}

func recordExists(group, date1, date2 string) (bool, string, string) {
	for fI, f := range header.Faculties {
		for cI, c := range f.Courses {
			for id, g := range c.Groups {
				if id == group {
					if g.TimeTable == nil {
						goto END
					}

					for _, table := range g.TimeTable {
						// API inconsistency: Datetime returned as 2020-10-19 whereas the request is submitted as 19.10.2020
						if dWR, _ := dateWithinRange(table.Date, "2006-01-02", date1, date2, "02.01.2006"); dWR {
							return true, fI, cI
						}
					}
				}
			}
		}
	}

END:
	// TODO: Check if no such date
	return false, "", ""
}

func dateWithinRange(dateStr, dateFormat, date1Str, date2Str, date12Format string) (bool, error) {
	dateStamp, _ := time.Parse(dateFormat, dateStr)
	date1Stamp, _ := time.Parse(date12Format, date1Str)
	date2Stamp, _ := time.Parse(date12Format, date2Str)

	if dateStamp.Equal(date1Stamp) || dateStamp.After(date1Stamp) && dateStamp.Before(date2Stamp) || dateStamp.Equal(date2Stamp) {
		return true, nil
	}

	return false, nil
}
