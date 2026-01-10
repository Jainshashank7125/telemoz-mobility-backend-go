package dto

type CreateTripRequest struct {
	ServiceType     string    `json:"service_type" binding:"required,oneof=delivery taxi school_bus"`
	PickupLocation  Location  `json:"pickup_location" binding:"required"`
	DropoffLocation Location  `json:"dropoff_location" binding:"required"`
	PaymentMethod   string    `json:"payment_method,omitempty"`
}

type TripResponse struct {
	ID                string    `json:"id"`
	CustomerID        string    `json:"customer_id"`
	DriverID          *string   `json:"driver_id,omitempty"`
	ServiceType       string    `json:"service_type"`
	Status            string    `json:"status"`
	PickupLocation    Location  `json:"pickup_location"`
	DropoffLocation   Location  `json:"dropoff_location"`
	EstimatedDistance *float64  `json:"estimated_distance,omitempty"`
	EstimatedDuration *int      `json:"estimated_duration,omitempty"`
	EstimatedArrival  *int64    `json:"estimated_arrival,omitempty"`
	FareAmount        *float64  `json:"fare_amount,omitempty"`
	PaymentMethod     string    `json:"payment_method"`
	CreatedAt         string    `json:"created_at"`
	UpdatedAt         string    `json:"updated_at"`
}

type UpdateTripRequest struct {
	Status            *string   `json:"status,omitempty"`
	EstimatedArrival  *int64    `json:"estimated_arrival,omitempty"`
	PickupLocation    *Location `json:"pickup_location,omitempty"`
	DropoffLocation   *Location `json:"dropoff_location,omitempty"`
}

