package page

import "gorm.io/gorm"

type Page struct {
	gorm.Model `json:"-"`
	Title      string `json:"title" gorm:"not null"`
	Content    string `json:"content" gorm:"not null"`
	Slug       string `json:"slug" gorm:"uniqueIndex;not null"`
	Published  bool   `json:"-" gorm:"default:false"`
}
