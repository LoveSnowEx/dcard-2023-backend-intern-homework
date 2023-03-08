package db

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UUID uuid.UUID

type UUIDModel struct {
	ID        UUID `gorm:"type:uuid;primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (u *UUID) BeforeCreate(tx *gorm.DB) (err error) {
	*u = UUID(uuid.New())
	return nil
}

func (u UUID) String() string {
	return uuid.UUID(u).String()
}
