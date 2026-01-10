package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
	"github.com/telemoz/backend/pkg/maps"
	"github.com/telemoz/backend/pkg/traccar"
)

type LocationService interface {
	UpdateBusLocation(busID uuid.UUID, lat, lng, accuracy, speed, heading float64) error
	CalculateETA(tripID uuid.UUID) (*time.Time, error)
	CalculateDistance(lat1, lng1, lat2, lng2 float64) (float64, error)
}

type locationService struct {
	busRepo         repositories.BusRepository
	busLocationRepo repositories.BusLocationRepository
	tripRepo        repositories.TripRepository
	traccarClient   *traccar.Client
	mapsClient      *maps.Client
}

func NewLocationService() LocationService {
	return &locationService{
		busRepo:         repositories.NewBusRepository(),
		busLocationRepo: repositories.NewBusLocationRepository(),
		tripRepo:        repositories.NewTripRepository(),
		traccarClient:   traccar.NewClient(),
		mapsClient:      maps.NewClient(),
	}
}

func (s *locationService) UpdateBusLocation(busID uuid.UUID, lat, lng, accuracy, speed, heading float64) error {
	location := &models.BusLocation{
		BusID:     busID,
		Latitude:  lat,
		Longitude: lng,
		Timestamp: time.Now(),
	}

	if accuracy > 0 {
		location.Accuracy = &accuracy
	}
	if speed >= 0 {
		location.Speed = &speed
	}
	if heading >= 0 {
		location.Heading = &heading
	}

	return s.busLocationRepo.Create(location)
}

func (s *locationService) CalculateETA(tripID uuid.UUID) (*time.Time, error) {
	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		return nil, errors.New("trip not found")
	}

	// Get current location of driver/bus
	var currentLat, currentLng float64
	if trip.DriverID != nil {
		// For now, use pickup location as current location
		// In production, get actual driver location from Traccar
		currentLat = trip.PickupLatitude
		currentLng = trip.PickupLongitude
	} else {
		currentLat = trip.PickupLatitude
		currentLng = trip.PickupLongitude
	}

	// Calculate distance and duration
	routeInfo, err := s.mapsClient.GetDistanceAndDuration(
		currentLat, currentLng,
		trip.DropoffLatitude, trip.DropoffLongitude,
	)
	if err != nil {
		return nil, errors.New("failed to calculate route")
	}

	eta := time.Now().Add(time.Duration(routeInfo.Duration) * time.Minute)
	return &eta, nil
}

func (s *locationService) CalculateDistance(lat1, lng1, lat2, lng2 float64) (float64, error) {
	routeInfo, err := s.mapsClient.GetDistanceAndDuration(lat1, lng1, lat2, lng2)
	if err != nil {
		return 0, err
	}
	return routeInfo.Distance, nil
}

