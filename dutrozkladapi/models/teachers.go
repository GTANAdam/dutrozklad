package models

type Kafedra struct {
	Name     string            `json:"name"`
	Teachers map[string]string `json:"teachers"`
}
