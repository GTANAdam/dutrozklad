package models

type Stats struct {
	Cached   int `json:"cached"`
	Uncached int `json:"uncached"`
}
