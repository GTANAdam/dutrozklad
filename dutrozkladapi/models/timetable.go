package models

type TimeTable struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Professor string `json:"professor"`
	Cabinet   string `json:"cabinet"`
	Date      string `json:"date"`
	Start     string `json:"start"`
	End       string `json:"end"`
	Misc      string `json:"misc"`
}
