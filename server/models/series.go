package models

import "gorm.io/gorm"

type Series struct {
	gorm.Model
	Title       string
	CoverUrl    string
	Description string
	Seasons     int32
}
