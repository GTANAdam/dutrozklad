package util

import (
	"bufio"
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/models"

	"github.com/gin-gonic/gin"
	"golang.org/x/net/html/charset"
)

func TimeTrack(start time.Time, name string) {
	elapsed := time.Since(start)
	log.Printf("[INFO] %s took %s", name, elapsed)
}

func DetectContentCharset(body io.Reader) string {
	r := bufio.NewReader(body)
	if data, err := r.Peek(1024); err == nil {
		if _, name, ok := charset.DetermineEncoding(data, ""); ok {
			return name
		}
	}
	return "utf-8"
}

func SaveToJSON(data interface{}, filename string) {
	file, _ := json.MarshalIndent(data, "", " ")
	if err := ioutil.WriteFile(filepath.FromSlash(filename), file, 0644); err != nil {
		log.Println("Error saving to json file.", err)
		return
	}
}

func RespondJSON(w *gin.Context, code int, message interface{}, payload interface{}) {
	w.JSON(200, &models.Response{Code: code, Error: message, Data: payload})
}

func LoadFaculties() {
	data, err := ioutil.ReadFile(filepath.FromSlash("data/faculties.json"))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &header.Faculties); err != nil {
		panic(err)
	}

	log.Println("Faculties loaded.")
}

func LoadTeachers() {
	data, err := ioutil.ReadFile(filepath.FromSlash("data/teachers.json"))
	if err != nil {
		panic(err)
	}

	if err := json.Unmarshal(data, &header.Kafedras); err != nil {
		panic(err)
	}

	log.Println("Teachers loaded.")
}

func AllowOrigin(cors string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Add("Access-Control-Allow-Origin", cors)
		c.Next()
	}
}

func TimeTableUnmarshal(content string) ([][]string, error) {
	// Find our regex match: this should return 2 groups
	if match := header.Regex.FindStringSubmatch(content); len(match) > 1 {
		// defer TimeTrack(time.Now(), "Unmarshal (TimeTable)")
		var data [][]string

		// TODO: Proper sanitation
		str := strings.NewReplacer(`\u000A`, "", `\n`, "", `''`, `""`, `['`, `["`, `','`, `","`, `']`, `"]`).Replace(match[1])

		// Unmarshal match
		if err := json.Unmarshal([]byte(str), &data); err != nil {
			if e, ok := err.(*json.SyntaxError); ok {
				log.Printf("syntax error at byte offset %d", e.Offset)
			}
			// Return error
			return nil, err
		}
		// Return found data
		return data, nil
	}
	// Return void
	return nil, nil
}

func CountTotalGroups() int {
	groupsCount := 0

	for _, faculty := range header.Faculties {
		for _, course := range faculty.Courses {
			groupsCount += len(course.Groups)
		}
	}

	return groupsCount
}

func CountTotalTeachers() int {
	teachersCount := 0

	for _, kafedra := range header.Kafedras {
		teachersCount += len(kafedra.Teachers)
	}

	return teachersCount
}
