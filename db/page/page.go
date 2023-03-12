package page

import "gorm.io/gorm"

type Page struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	Published bool   `gorm:"default:false"`
}
