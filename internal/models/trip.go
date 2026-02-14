package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ServiceType string

const (
	ServiceTypeDelivery  ServiceType = "delivery"
	ServiceTypeTaxi      ServiceType = "taxi"
	ServiceTypeSchoolBus ServiceType = "school_bus"
)

type TripStatus string

const (
	TripStatusPending             TripStatus = "pending"
	TripStatusSearching           TripStatus = "searching"
	TripStatusAccepted            TripStatus = "accepted"
	TripStatusInProgress          TripStatus = "in_progress"
	TripStatusCompleted           TripStatus = "completed"
	TripStatusCancelled           TripStatus = "cancelled"
	TripStatusExpired             TripStatus = "expired"
	TripStatusCancelledByCustomer TripStatus = "cancelled_by_customer"
)

type Trip struct {
	ID                uuid.UUID   `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	CustomerID        uuid.UUID   `gorm:"type:uuid;not null;index" json:"customer_id"`
	DriverID          *uuid.UUID  `gorm:"type:uuid;index" json:"driver_id,omitempty"`
	ServiceType       ServiceType `gorm:"type:varchar(20);not null;index" json:"service_type"`
	Status            TripStatus  `gorm:"type:varchar(30);not null;default:'pending';index" json:"status"`
	PickupLatitude    float64     `gorm:"type:decimal(10,8);not null" json:"pickup_latitude"`
	PickupLongitude   float64     `gorm:"type:decimal(11,8);not null" json:"pickup_longitude"`
	PickupAddress     *string     `gorm:"type:text" json:"pickup_address,omitempty"`
	DropoffLatitude   float64     `gorm:"type:decimal(10,8);not null" json:"dropoff_latitude"`
	DropoffLongitude  float64     `gorm:"type:decimal(11,8);not null" json:"dropoff_longitude"`
	DropoffAddress    *string     `gorm:"type:text" json:"dropoff_address,omitempty"`
	EstimatedDistance *float64    `gorm:"type:decimal(10,2)" json:"estimated_distance,omitempty"`
	EstimatedDuration *int        `gorm:"type:integer" json:"estimated_duration,omitempty"`
	EstimatedArrival  *time.Time  `json:"estimated_arrival,omitempty"`
	FareAmount        *float64    `gorm:"type:decimal(10,2)" json:"fare_amount,omitempty"`
	PaymentMethod     string      `gorm:"type:varchar(20);default:'cash'" json:"payment_method"`
	TraccarDeviceID   *string     `gorm:"type:varchar(255)" json:"traccar_device_id,omitempty"`

	// Search tracking
	SearchStartedAt    *time.Time `json:"search_started_at,omitempty"`
	SearchEndedAt      *time.Time `json:"search_ended_at,omitempty"`
	CancelledBy        *uuid.UUID `gorm:"type:uuid" json:"cancelled_by,omitempty"`
	CancellationReason *string    `gorm:"type:text" json:"cancellation_reason,omitempty"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relations
	Customer User  `gorm:"foreignKey:CustomerID" json:"customer,omitempty"`
	Driver   *User `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
}

func (t *Trip) BeforeCreate(tx *gorm.DB) error {
	if t.ID == uuid.Nil {
		t.ID = uuid.New()
	}
	return nil
}
