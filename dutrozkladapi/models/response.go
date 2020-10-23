package models

type Response struct {
	Code  int         `json:"code"`
	Error interface{} `json:"error"`
	Data  interface{} `json:"content"`
}
