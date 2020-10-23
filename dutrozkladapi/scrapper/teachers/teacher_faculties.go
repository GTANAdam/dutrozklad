package teachers

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

func GetTeacherKafedra() {
	defer util.TimeTrack(time.Now(), "GetTeacherKafedra")

	response := util.Get(header.TeacherURL)
	defer response.Close()

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(response)
	if err != nil {
		log.Println(err)
	}

	// Look for Faculty data
	node := doc.Find("#TimeTableForm_chair")
	if node.Length() == 0 {
		log.Println("[ERROR] GetTeacherKafedra, empty data.")
		return
	}

	// Extract and populate Data map with faculty information
	sel := node.Eq(0).Children()
	if sel.Length() == 1 && sel.Nodes[0].Attr[0].Val == "" || sel.Length() == 0 {
		log.Println("[INFO] There are no Kafedras.")
		return
	}

	for i, obj := range sel.Nodes {
		val := obj.Attr[0].Val
		text := strings.TrimSpace(sel.Eq(i).Text())

		if val != "" {
			// Append to Kafedras a kafedra
			header.Kafedras[val] = &models.Kafedra{Name: text}
		}
	}

	log.Println("[INFO] Found " + fmt.Sprint(len(header.Kafedras)) + " teacher kafedras.")
}
