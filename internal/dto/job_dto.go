package dto

type JobResponse struct {
	ID              string    `json:"id"`
	TripID          string    `json:"trip_id"`
	DriverID        *string   `json:"driver_id,omitempty"`
	ServiceType     string    `json:"service_type"`
	Status          string    `json:"status"`
	PickupLocation  Location  `json:"pickup_location"`
	DropoffLocation Location  `json:"dropoff_location"`
	EstimatedEarnings *float64 `json:"estimated_earnings,omitempty"`
	Distance        *float64  `json:"distance,omitempty"`
	CustomerID      string    `json:"customer_id"`
	CreatedAt       string    `json:"created_at"`
}

type UpdateJobStatusRequest struct {
	Status string `json:"status" binding:"required,oneof=accepted rejected in_progress completed"`
}

