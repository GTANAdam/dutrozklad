// Package query ..
package query

import (
	"dutrozkladbot/config"
	"dutrozkladbot/model"
	"encoding/json"
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

// GetStats ..
func GetStats() (*model.StatsResponse, error) {
	stuff, err := query(fmt.Sprintf("%s/stats", config.APIAddress))
	if err != nil {
		return nil, err
	}

	var result model.StatsResponse
	if err := json.Unmarshal(stuff.Data, &result); err != nil {
		return nil, err
	}

	return &result, nil
}
