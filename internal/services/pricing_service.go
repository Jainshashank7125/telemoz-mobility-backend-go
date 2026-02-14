package services

import (
	"math"

	"github.com/telemoz/backend/internal/config"
)

type PricingService interface {
	CalculateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64, serviceType string) (distance float64, fare float64)
	EstimateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64, serviceType string) (distance float64, fare float64, duration float64)
}

type pricingService struct{}

func NewPricingService() PricingService {
	return &pricingService{}
}

// CalculateFare calculates fare based on distance using Haversine formula
func (s *pricingService) CalculateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64, serviceType string) (float64, float64) {
	// Calculate distance in kilometers
	distance := calculateHaversineDistance(pickupLat, pickupLng, dropoffLat, dropoffLng)

	// Get pricing config for service type
	pricing := config.GetPricingForService(serviceType)

	// Calculate fare: base fare + (distance * per km rate) * surge multiplier
	fare := (pricing.BaseFare + (distance * pricing.PerKmRate)) * pricing.SurgeMultiplier

	// Apply minimum fare
	if fare < pricing.MinimumFare {
		fare = pricing.MinimumFare
	}

	// Round to 2 decimal places
	distance = math.Round(distance*100) / 100
	fare = math.Round(fare*100) / 100

	return distance, fare
}

// EstimateFare estimates fare and adds estimated duration
func (s *pricingService) EstimateFare(pickupLat, pickupLng, dropoffLat, dropoffLng float64, serviceType string) (float64, float64, float64) {
	distance, fare := s.CalculateFare(pickupLat, pickupLng, dropoffLat, dropoffLng, serviceType)

	// Estimate duration: assume average speed of 40 km/h in city
	// Duration in minutes
	duration := (distance / 40.0) * 60.0
	duration = math.Round(duration*100) / 100

	return distance, fare, duration
}

// calculateHaversineDistance calculates the distance between two points on Earth
// using the Haversine formula. Returns distance in kilometers.
func calculateHaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const earthRadiusKm = 6371.0

	// Convert degrees to radians
	lat1Rad := lat1 * math.Pi / 180
	lat2Rad := lat2 * math.Pi / 180
	deltaLat := (lat2 - lat1) * math.Pi / 180
	deltaLon := (lon2 - lon1) * math.Pi / 180

	// Haversine formula
	a := math.Sin(deltaLat/2)*math.Sin(deltaLat/2) +
		math.Cos(lat1Rad)*math.Cos(lat2Rad)*
			math.Sin(deltaLon/2)*math.Sin(deltaLon/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	return earthRadiusKm * c
}
