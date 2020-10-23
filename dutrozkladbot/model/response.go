// Package model ..
package model

import "encoding/json"

// Response ..
type Response struct {
	Code  int             `json:"code"`
	Error interface{}     `json:"error"`
	Data  json.RawMessage `json:"content"`
}
