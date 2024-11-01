package models

import "gorm.io/gorm"

type Episode struct {
	gorm.Model
	SeriesID    int32
	Season      int32
	Title       string
	Description string
	ContentUrl  string
	CoverUrl    string
}
