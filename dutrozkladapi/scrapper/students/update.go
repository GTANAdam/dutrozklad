package students

import (
	"fmt"
	"log"
	"sync"
	"time"

	"dutrozkladapi/header"
	"dutrozkladapi/models"
	"dutrozkladapi/util"
)

var mutex = &sync.Mutex{}
var wg = &sync.WaitGroup{}
var hResult = make(map[string]*models.Faculty)
var totalGroups int
var failedParse int

// UpdateFaculties ..
func UpdateFaculties() {
	// Get faculties listing
	faculties, err := GetFaculties()
	if err != nil {
		panic("Something went wrong with faculties.")
	}

	// Assign faculties to hResult
	hResult = faculties

	for facultyID := range faculties {
		for _, courseID := range getCourses(facultyID) {

			groups := getGroups(facultyID, courseID)
			wg.Add(len(groups))

			for groupID := range groups {
				go getTimeTables(facultyID, courseID, groupID)
			}
		}
	}

	// Wait for all timetables to finish fetching
	wg.Wait()

	// Save to global data
	header.Mutex.Lock()
	header.Faculties = hResult
	header.Mutex.Unlock()

	util.SaveToJSON(header.Faculties, "data/faculties.json")
	log.Println("> Faculties updated. Failed/Total:", failedParse, "/", totalGroups)
}

func getCourses(facultyID string) []string {
	// Get courses of this faculty
	courses, err := GetCourses(facultyID)
	if err != nil {
		log.Println(err)
		return nil
	}

	return courses
}

func getGroups(facultyID, courseID string) map[string]*models.Group {
	// Get groups of this course
	groups, err := GetGroups(facultyID, courseID)
	if err != nil {
		log.Println(err)
		return nil
	}

	// Assign groups to this course
	mutex.Lock()
	hResult[facultyID].Courses[courseID] = &models.Course{Groups: groups}
	mutex.Unlock()

	return groups
}

func getTimeTables(facultyID, courseID, groupID string) {
	totalGroups++

	// Set up date
	tYear, tMonth, tDay := time.Now().Date()
	nYear, nMonth, nDay := time.Now().AddDate(1, 0, 0).Date() // 1 year
	todayStr := fmt.Sprintf("%v.%v.%v", tDay, int(tMonth), tYear)
	nextStr := fmt.Sprintf("%v.%v.%v", nDay, int(nMonth), nYear)

	// Get timetables of this group
	result, err := GetTimeTable(facultyID, courseID, groupID, todayStr, nextStr)
	if err != nil {
		log.Println(err)
		failedParse++
		goto END
	}

	// Assign timtables to this group
	mutex.Lock()
	hResult[facultyID].Courses[courseID].Groups[groupID].TimeTable = result
	mutex.Unlock()

END:
	wg.Done()
}

// func UpdateTeachers() {
// 	GetTeacherKafedra()

// 	for kafedraID := range header.Kafedras {
// 		GetTeachers(kafedraID)
// 	}

// 	util.SaveToJSON(header.Kafedras, "data/teachers.json")
// 	log.Println("> Teachers updated.")
// }
