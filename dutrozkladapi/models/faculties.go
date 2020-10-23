package models

type Faculty struct {
	// ID      string             `json:"id"`
	Name    string             `json:"name"`
	Courses map[string]*Course `json:"courses"`
}

type Course struct {
	// ID     string            `json:"id"`
	Groups map[string]*Group `json:"groups"`
}

type Group struct {
	// ID        string       `json:"id"`
	Name      string       `json:"name"`
	TimeTable []*TimeTable `json:"timetable"`
}
