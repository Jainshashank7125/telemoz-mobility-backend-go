package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverAvailability struct {
	DriverID     uuid.UUID `gorm:"type:uuid;primary_key" json:"driver_id"`
	IsAvailable  bool      `gorm:"default:false;index" json:"is_available"`
	ServiceTypes []string  `gorm:"type:text[]" json:"service_types"`
	LastActiveAt time.Time `json:"last_active_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// Relations
	Driver User `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
}

func (da *DriverAvailability) BeforeCreate(tx *gorm.DB) error {
	if da.LastActiveAt.IsZero() {
		da.LastActiveAt = time.Now()
	}
	return nil
}
