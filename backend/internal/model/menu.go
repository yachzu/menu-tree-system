package model

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Menu struct {
	ID         uuid.UUID  `gorm:"type:uuid;primaryKey;default:gen_random_uuid()" json:"id"`
	Name       string     `gorm:"type:varchar(255);not null" json:"name" binding:"required"`
	ParentID   *uuid.UUID `gorm:"type:uuid;index" json:"parent_id"`
	Parent     *Menu      `gorm:"foreignKey:ParentID;constraint:OnDelete:CASCADE" json:"-"`
	Depth      int        `gorm:"not null;default:0" json:"depth"`
	OrderIndex int        `gorm:"not null;default:0;index:idx_parent_order" json:"order_index"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	Children   []*Menu    `gorm:"foreignKey:ParentID" json:"children,omitempty"`
}

func (m *Menu) BeforeCreate(tx *gorm.DB) error {
	if m.ID == uuid.Nil {
		m.ID = uuid.New()
	}
	return nil
}
