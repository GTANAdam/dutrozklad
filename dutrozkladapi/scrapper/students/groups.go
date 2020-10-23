package students

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/models"
	"dutrozkladapi/util"

	"github.com/PuerkitoBio/goquery"
)

func GetGroups(facultyID string, course string) (map[string]*models.Group, error) {
	defer util.TimeTrack(time.Now(), "GetGroups")

	formData := url.Values{
		"TimeTableForm[faculty]": {facultyID},
		"TimeTableForm[course]":  {course},
	}

	response, err := util.Post(header.URL, formData)
	if err != nil {
		return nil, err
	}
	defer response.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response)
	if err != nil {
		log.Println(err)
	}

	// Look for course data
	node := doc.Find("#TimeTableForm_group")
	if node.Length() == 0 {
		return nil, fmt.Errorf("getGroups, empty data")
	}

	// Extract and populate Data map with group information
	result := make(map[string]*models.Group)
	sel := node.Eq(0).Children()
	if sel.Length() == 1 && sel.Nodes[0].Attr[0].Val == "" || sel.Length() == 0 {
		log.Println()
		return nil, fmt.Errorf("there are no groups")
	}

	for i, obj := range sel.Nodes {
		val := obj.Attr[0].Val
		text := strings.TrimSpace(sel.Eq(i).Text())

		if val != "" {
			// Add group array to global faculties
			result[val] = &models.Group{Name: text}
		}
	}

	log.Println("[INFO] Found " + fmt.Sprint(len(result)) + " groups for course " + course + ", faculty: " + facultyID)

	return result, nil
}
