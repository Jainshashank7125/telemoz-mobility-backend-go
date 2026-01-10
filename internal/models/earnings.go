package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type DriverEarning struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	DriverID  uuid.UUID `gorm:"type:uuid;not null;index" json:"driver_id"`
	JobID     uuid.UUID `gorm:"type:uuid;not null" json:"job_id"`
	Amount    float64   `gorm:"type:decimal(10,2);not null" json:"amount"`
	Date      time.Time `gorm:"type:date;not null;index" json:"date"`
	CreatedAt time.Time `json:"created_at"`

	// Relations
	Driver User `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
	Job    Job  `gorm:"foreignKey:JobID" json:"job,omitempty"`
}

func (de *DriverEarning) BeforeCreate(tx *gorm.DB) error {
	if de.ID == uuid.Nil {
		de.ID = uuid.New()
	}
	return nil
}

