package router

import (
	"dutrozkladapi/header"
	"dutrozkladapi/scrapper/teachers"
	"dutrozkladapi/util"

	"github.com/gin-gonic/gin"
)

func GetKafedras(c *gin.Context) {
	result := make(map[string]string)
	for k, v := range header.Kafedras {
		result[k] = v.Name
	}

	util.RespondJSON(c, 200, nil, result)
}

func GetTeachers(c *gin.Context) {
	kafedra := c.Query("kafedra")
	if kafedra == "" {
		util.RespondJSON(c, 400, "missing or invalid kafedra parameter.", nil)
		return
	}

	util.RespondJSON(c, 200, nil, header.Kafedras[kafedra].Teachers)
}

func GetTeacherTimeTable(c *gin.Context) {
	kafedra := c.PostForm("kafedra")
	if kafedra == "" {
		util.RespondJSON(c, 400, "missing or invalid kafedra parameter.", nil)
		return
	}

	teacher := c.PostForm("teacher")
	if teacher == "" {
		util.RespondJSON(c, 400, "missing or invalid teacher parameter.", nil)
		return
	}

	startdate := c.PostForm("startdate")
	if startdate == "" {
		util.RespondJSON(c, 400, "missing or invalid startdate parameter.", nil)
		return
	}

	enddate := c.PostForm("enddate")
	if enddate == "" {
		util.RespondJSON(c, 400, "missing or invalid enddate parameter.", nil)
		return
	}

	result := teachers.GetTeacherTimeTable(kafedra, teacher, startdate, enddate)
	if result == nil {
		util.RespondJSON(c, 404, "empty timetable.", nil)
		return
	}

	util.RespondJSON(c, 200, nil, result)
}
