package students

import (
	"fmt"
	"log"
	"net/url"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/util"

	"github.com/PuerkitoBio/goquery"
)

func GetCourses(facultyID string) ([]string, error) {
	defer util.TimeTrack(time.Now(), "GetCourses")

	formData := url.Values{
		"TimeTableForm[faculty]": {facultyID},
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
	node := doc.Find("#TimeTableForm_course")
	if node.Length() == 0 {
		return nil, fmt.Errorf("getCourses, empty data")
	}

	// Extract and populate Data map with courses information
	var result []string
	sel := node.Eq(0).Children()
	if sel.Length() == 1 && sel.Nodes[0].Attr[0].Val == "" || sel.Length() == 0 {
		return nil, fmt.Errorf("there are no courses")
	}

	for _, obj := range sel.Nodes {
		val := obj.Attr[0].Val

		if val != "" {
			result = append(result, obj.Attr[0].Val)
		}
	}

	log.Println("[INFO] Found " + fmt.Sprint(len(result)) + " courses for faculty: " + facultyID)
	return result, nil
}
