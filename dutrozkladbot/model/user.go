// Package model ..
package model

// User ..
type User struct {
	ID              int    `gorm:"primary_key:true" json:"-"`
	UserID          int    `gorm:"column:UserID" json:"UserID"`
	Faculty         int    `gorm:"column:Faculty" json:"Faculty"`
	Course          int    `gorm:"column:Course" json:"Course"`
	Group           int    `gorm:"column:Group" json:"Group"`
	GroupName       string `gorm:"column:GroupName" json:"GroupName"`
	LastQueriedDate string `gorm:"column:LastQueriedDate" json:"Date"`
}
