package pagelist

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PageList struct {
	gorm.Model
	Key uuid.UUID `gorm:"type:uuid;uniqueIndex"`
}

func New() *PageList {
	return &PageList{}
}

func (l *PageList) BeforeCreate(tx *gorm.DB) error {
	if l.Key == uuid.Nil {
		l.Key = uuid.New()
	}
	return nil
}
