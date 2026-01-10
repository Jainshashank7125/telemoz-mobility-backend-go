package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Child struct {
	ID         uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	ParentID   uuid.UUID `gorm:"type:uuid;not null;index" json:"parent_id"`
	Name       string    `gorm:"not null" json:"name"`
	SchoolName *string   `gorm:"type:varchar(255)" json:"school_name,omitempty"`
	BusID      *uuid.UUID `gorm:"type:uuid;index" json:"bus_id,omitempty"`
	AvatarURL  *string   `json:"avatar_url,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`

	// Relations
	Parent User  `gorm:"foreignKey:ParentID" json:"parent,omitempty"`
	Bus    *Bus  `gorm:"foreignKey:BusID" json:"bus,omitempty"`
}

func (c *Child) BeforeCreate(tx *gorm.DB) error {
	if c.ID == uuid.Nil {
		c.ID = uuid.New()
	}
	return nil
}

