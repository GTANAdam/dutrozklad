package router

import (
	"dutrozkladapi/header"
	"dutrozkladapi/scrapper/students"
	"dutrozkladapi/util"
	"log"
	"sort"

	"github.com/gin-gonic/gin"
)

func GetFaculties(c *gin.Context) {
	result := make(map[string]string)
	for k, v := range header.Faculties {
		result[k] = v.Name
	}

	util.RespondJSON(c, 200, nil, result)
}

func GetCourses(c *gin.Context) {
	faculty := c.Param("faculty")
	if faculty == "" {
		util.RespondJSON(c, 400, "missing or invalid faculty parameter", nil)
		return
	}

	// if len(header.Faculties[faculty].Courses) == 0 {
	// 	util.RespondJSON(c, 400, "Empty course list.", nil)
	// 	return
	// }

	if header.Faculties[faculty] == nil {
		util.RespondJSON(c, 400, "empty set.", nil)
		return
	}

	result := make([]string, 0, len(header.Faculties[faculty].Courses))
	for k := range header.Faculties[faculty].Courses {
		result = append(result, k)
	}

	sort.Slice(result, func(i, j int) bool {
		return result[i] < result[j]
	})

	util.RespondJSON(c, 200, nil, result)
}

func GetGroups(c *gin.Context) {
	faculty := c.Param("faculty")
	if faculty == "" {
		util.RespondJSON(c, 400, "missing or invalid faculty parameter.", nil)
		return
	}

	course := c.Param("course")
	if course == "" {
		util.RespondJSON(c, 400, "missing or invalid course parameter.", nil)
		return
	}

	result := make(map[string]string)
	for k, v := range header.Faculties[faculty].Courses[course].Groups {
		result[k] = v.Name
	}

	util.RespondJSON(c, 200, nil, result)
}

// TODO: param binding
func GetTimeTable(c *gin.Context) {
	// faculty := c.PostForm("faculty")
	// if faculty == "" {
	// 	util.RespondJSON(c, 400, "missing or invalid faculty parameter.", nil)
	// 	return
	// }

	// course := c.PostForm("course")
	// if course == "" {
	// 	util.RespondJSON(c, 400, "missing or invalid course parameter.", nil)
	// 	return
	// }

	group := c.Param("group")
	if group == "" {
		util.RespondJSON(c, 400, "missing or invalid group parameter.", nil)
		return
	}

	startdate := c.Param("startdate")
	if startdate == "" {
		util.RespondJSON(c, 400, "missing or invalid startdate parameter.", nil)
		return
	}

	enddate := c.Param("enddate")
	if enddate == "" {
		enddate = startdate
	}

	result, err := students.GetTimeTable("0", "0", group, startdate, enddate)
	if err != nil {
		log.Println(err)
		util.RespondJSON(c, 404, err, nil)
		return
	}

	if result == nil {
		util.RespondJSON(c, 404, "empty timetable.", nil)
		return
	}

	util.RespondJSON(c, 200, nil, result)
}
