package db

import (
	"gorm.io/gorm"
)

type Page struct {
	gorm.Model
	Title     string `gorm:"not null"`
	Content   string `gorm:"type:text;not null"`
	Slug      string `gorm:"uniqueIndex;not null"`
	Published bool   `gorm:"default:false"`
}

// Create a new page
func (db *DB) CreatePage(page *Page) error {
	return db.DB.Create(&page).Error
}

// Retrieve a page by ID
func (db *DB) GetPageByID(id uint) (*Page, error) {
	var page Page
	err := db.DB.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Retrieve a page by slug
func (db *DB) GetPageBySlug(slug string) (*Page, error) {
	var page Page
	err := db.DB.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Retrieve all pages
func (db *DB) GetPages() ([]Page, error) {
	var pages []Page
	err := db.DB.Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// Update a page
func (db *DB) UpdatePage(page *Page) error {
	return db.DB.Save(&page).Error
}

// Delete a page
func (db *DB) DeletePage(id uint) error {
	return db.DB.Delete(&Page{}, id).Error
}
