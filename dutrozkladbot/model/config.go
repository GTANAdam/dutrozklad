// Package model ..
package model

// Config ..
type Config struct {
	Token      *Type `json:"token"`
	APIAddress *Type `json:"api_address"`
}

// Type ..
type Type struct {
	Production  string `json:"production"`
	Development string `json:"development"`
}
