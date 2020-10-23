package header

import (
	"regexp"
	"sync"

	"dutrozkladapi/models"
)

const (
	Name = iota
	Type
	Prof
	Cabinet
	Date
	Start
	End
	Misc
)

const (
	TimeTableNodeID = 29
	URL             = "http://e-rozklad.dut.edu.ua/timeTable/group"
	TeacherURL      = "http://e-rozklad.dut.edu.ua/timeTable/teacher"
)

var Faculties = make(map[string]*models.Faculty)
var Kafedras = make(map[string]*models.Kafedra)
var Regex, _ = regexp.Compile(`var arr=(.*);`)
var Mutex = &sync.Mutex{}
