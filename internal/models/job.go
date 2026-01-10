package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type JobStatus string

const (
	JobStatusPending    JobStatus = "pending"
	JobStatusAccepted   JobStatus = "accepted"
	JobStatusRejected   JobStatus = "rejected"
	JobStatusInProgress JobStatus = "in_progress"
	JobStatusCompleted  JobStatus = "completed"
)

type Job struct {
	ID              uuid.UUID  `gorm:"type:uuid;primary_key;default:gen_random_uuid()" json:"id"`
	TripID          uuid.UUID  `gorm:"type:uuid;not null;index" json:"trip_id"`
	DriverID        *uuid.UUID `gorm:"type:uuid;index" json:"driver_id,omitempty"`
	Status          JobStatus  `gorm:"type:varchar(20);not null;default:'pending';index" json:"status"`
	EstimatedEarnings *float64  `gorm:"type:decimal(10,2)" json:"estimated_earnings,omitempty"`
	ActualEarnings  *float64   `gorm:"type:decimal(10,2)" json:"actual_earnings,omitempty"`
	Distance        *float64   `gorm:"type:decimal(10,2)" json:"distance,omitempty"`
	AcceptedAt      *time.Time `json:"accepted_at,omitempty"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`

	// Relations
	Trip   Trip  `gorm:"foreignKey:TripID" json:"trip,omitempty"`
	Driver *User `gorm:"foreignKey:DriverID" json:"driver,omitempty"`
}

func (j *Job) BeforeCreate(tx *gorm.DB) error {
	if j.ID == uuid.Nil {
		j.ID = uuid.New()
	}
	return nil
}

