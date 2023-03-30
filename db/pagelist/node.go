package pagelist

import (
	"fmt"

	"github.com/google/uuid"

	"gorm.io/gorm"
)

type PageNode struct {
	gorm.Model
	End     bool
	Key     uuid.UUID `gorm:"type:uuid;uniqueIndex"`
	PrevKey uuid.UUID `gorm:"type:uuid"`
	NextKey uuid.UUID `gorm:"type:uuid"`
	// Page   *page.Page `gorm:"foreignkey:PageID;references:ID"`
	PageID  uint
	List    PageList  `gorm:"foreignkey:ListKey;references:Key"`
	ListKey uuid.UUID `gorm:"type:uuid"`
}

func (n *PageNode) BeforeCreate(tx *gorm.DB) error {
	if n.Key == uuid.Nil {
		n.Key = uuid.New()
	}
	if n.PrevKey == uuid.Nil {
		n.PrevKey = n.Key
	} else if tx.First(&PageNode{Key: n.PrevKey}).Error != nil {
		return fmt.Errorf("before create page node: prev key not found")
	}
	if n.NextKey == uuid.Nil {
		n.NextKey = n.Key
	} else if tx.First(&PageNode{Key: n.NextKey}).Error != nil {
		return fmt.Errorf("before create page node: next key not found")
	}
	return nil
}
