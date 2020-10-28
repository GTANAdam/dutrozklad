package model

// StatsResponse ..
type StatsResponse struct {
	Cached   int `json:"cached"`
	Uncached int `json:"uncached"`
}
