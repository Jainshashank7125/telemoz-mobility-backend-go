package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Bus struct {
	ID              uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	Name            string    `gorm:"not null" json:"name"`
	DriverID        uuid.UUID `gorm:"type:uuid;not null;index" json:"driver_id"`
	TraccarDeviceID *string   `gorm:"type:varchar(255);uniqueIndex" json:"traccar_device_id,omitempty"`
	RouteName       *string   `gorm:"type:varchar(255)" json:"route_name,omitempty"`
	IsActive        bool      `gorm:"default:true" json:"is_active"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`

	// Relations
	Driver   User           `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
	Children []Child        `gorm:"foreignKey:BusID" json:"children,omitempty"`
	Locations []BusLocation `gorm:"foreignKey:BusID" json:"locations,omitempty"`
}

func (b *Bus) BeforeCreate(tx *gorm.DB) error {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return nil
}

type BusLocation struct {
	ID        uuid.UUID `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	BusID     uuid.UUID `gorm:"type:uuid;not null;index" json:"bus_id"`
	Latitude  float64   `gorm:"type:decimal(10,8);not null" json:"latitude"`
	Longitude float64   `gorm:"type:decimal(11,8);not null" json:"longitude"`
	Accuracy  *float64  `gorm:"type:decimal(10,2)" json:"accuracy,omitempty"`
	Speed     *float64  `gorm:"type:decimal(10,2)" json:"speed,omitempty"`
	Heading   *float64  `gorm:"type:decimal(5,2)" json:"heading,omitempty"`
	Timestamp time.Time `gorm:"default:CURRENT_TIMESTAMP;index" json:"timestamp"`

	// Relations
	Bus Bus `gorm:"foreignKey:BusID" json:"bus,omitempty"`
}

func (bl *BusLocation) BeforeCreate(tx *gorm.DB) error {
	if bl.ID == uuid.Nil {
		bl.ID = uuid.New()
	}
	return nil
}

