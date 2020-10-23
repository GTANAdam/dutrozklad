package teachers

import (
	"fmt"
	"log"
	"net/url"
	"strings"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/util"

	"github.com/PuerkitoBio/goquery"
)

func GetTeachers(KafedraID string) {
	defer util.TimeTrack(time.Now(), "GetTeachers")

	formData := url.Values{
		"TimeTableForm[chair]": {KafedraID},
	}

	response, err := util.Post(header.TeacherURL, formData)
	if err != nil {
		log.Println(err)
		return
	}

	defer response.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response)
	if err != nil {
		log.Println(err)
	}

	// Look for course data
	node := doc.Find("#TimeTableForm_teacher")
	if node.Length() == 0 {
		log.Println("[ERROR] GetTeachers, empty data.")
		return
	}

	// Extract and populate Data map with group information
	sel := node.Eq(0).Children()
	if sel.Length() == 1 && sel.Nodes[0].Attr[0].Val == "" || sel.Length() == 0 {
		log.Println("[INFO] Kafedra " + header.Kafedras[KafedraID].Name + " is empty.")
		return
	}

	for i, obj := range sel.Nodes {
		val := obj.Attr[0].Val
		text := strings.TrimSpace(sel.Eq(i).Text())

		if val != "" {
			if header.Kafedras[KafedraID].Teachers == nil {
				// get local faculties
				fac := header.Kafedras[KafedraID]

				// Create local courses map
				fac.Teachers = make(map[string]string)

				// overwrite global faculty with local
				header.Kafedras[KafedraID] = fac
			}

			// append to local result array
			header.Kafedras[KafedraID].Teachers[val] = text
		}
	}

	log.Println("[INFO] Found " + fmt.Sprint(len(header.Kafedras[KafedraID].Teachers)) + " teachers for kafedra: " + header.Kafedras[KafedraID].Name)
}
