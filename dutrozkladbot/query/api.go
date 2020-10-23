// Package query ..
package query

import (
	"dutrozkladbot/config"
	"dutrozkladbot/model"
	"fmt"
)

// GetFaculties ..
func GetFaculties() (*model.Response, error) {
	return query(fmt.Sprintf("%s/faculties", config.APIAddress))
}

// GetCourses ..
func GetCourses(faculty int) (*model.Response, error) {
	return query(fmt.Sprintf("%s/courses/%v", config.APIAddress, faculty))
}

// GetGroups ..
func GetGroups(faculty, course int) (*model.Response, error) {
	return query(fmt.Sprintf("%s/groups/%v/%v", config.APIAddress, faculty, course))
}

// GetTimeTable ..
func GetTimeTable(group int, date string) (*model.Response, error) {
	return query(fmt.Sprintf("%s/timetable/%v/%s/%s", config.APIAddress, group, date, date))
}
