package db

import "github.com/LoveSnowEx/dcard-2023-backend-intern-homework/db/page"

// Create a new page
func (db *DB) CreatePage(page *page.Page) error {
	return db.DB.Create(&page).Error
}

// Retrieve a page by ID
func (db *DB) GetPageByID(id uint) (*page.Page, error) {
	var page page.Page
	err := db.DB.First(&page, id).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Retrieve a page by slug
func (db *DB) GetPageBySlug(slug string) (*page.Page, error) {
	var page page.Page
	err := db.DB.Where("slug = ?", slug).First(&page).Error
	if err != nil {
		return nil, err
	}
	return &page, nil
}

// Retrieve all pages
func (db *DB) GetPages() ([]page.Page, error) {
	var pages []page.Page
	err := db.DB.Find(&pages).Error
	if err != nil {
		return nil, err
	}
	return pages, nil
}

// Update a page
func (db *DB) UpdatePage(page *page.Page) error {
	return db.DB.Save(&page).Error
}

// Delete a page
func (db *DB) DeletePage(id uint) error {
	return db.DB.Delete(&page.Page{}, id).Error
}
