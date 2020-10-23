package students

import (
	"fmt"
	"log"
	"strings"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/models"
	"dutrozkladapi/util"

	"github.com/PuerkitoBio/goquery"
)

func GetFaculties() (map[string]*models.Faculty, error) {
	defer util.TimeTrack(time.Now(), "GetFaculties")

	faculties := make(map[string]*models.Faculty)

	response := util.Get(header.URL)
	defer response.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response)
	if err != nil {
		log.Println(err)
	}

	// Look for Faculty data
	node := doc.Find("#TimeTableForm_faculty")
	if node.Length() == 0 {
		return nil, fmt.Errorf("getFaculties, empty data")
	}

	// Extract and populate Data map with faculty information
	sel := node.Eq(0).Children()
	if sel.Length() == 1 && sel.Nodes[0].Attr[0].Val == "" || sel.Length() == 0 {
		return nil, fmt.Errorf("there are no faculties")
	}

	for i, obj := range sel.Nodes {
		val := obj.Attr[0].Val
		text := strings.TrimSpace(sel.Eq(i).Text())

		if val != "" {
			// Append to Faculties a faculty
			faculties[val] = &models.Faculty{Name: text, Courses: make(map[string]*models.Course)}
		}
	}

	log.Println("[INFO] Found " + fmt.Sprint(len(faculties)) + " faculties.")

	return faculties, nil
}
