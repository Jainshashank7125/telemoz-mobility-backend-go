package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
	"github.com/telemoz/backend/internal/utils"
)

type TripService interface {
	CreateTrip(customerID uuid.UUID, req dto.CreateTripRequest) (*dto.TripResponse, error)
	EstimateFare(req dto.EstimateFareRequest) (*dto.EstimateFareResponse, error)
	GetActiveTrip(customerID uuid.UUID) (*dto.TripResponse, error)
	GetTripHistory(customerID uuid.UUID, limit, offset int) ([]dto.TripResponse, error)
	GetTripByID(tripID uuid.UUID) (*dto.TripResponse, error)
	UpdateTrip(tripID uuid.UUID, req dto.UpdateTripRequest) (*dto.TripResponse, error)
	CancelTrip(tripID uuid.UUID, customerID uuid.UUID) error
	AcceptTrip(tripID uuid.UUID, driverID uuid.UUID) error
	GetTripStatus(tripID uuid.UUID) (string, error)
	ExpireSearchingTrips() error
}

type tripService struct {
	tripRepo       repositories.TripRepository
	jobRepo        repositories.JobRepository
	userRepo       repositories.UserRepository
	pricingService PricingService
}

func NewTripService() TripService {
	return &tripService{
		tripRepo:       repositories.NewTripRepository(),
		jobRepo:        repositories.NewJobRepository(),
		userRepo:       repositories.NewUserRepository(),
		pricingService: NewPricingService(),
	}
}

func (s *tripService) CreateTrip(customerID uuid.UUID, req dto.CreateTripRequest) (*dto.TripResponse, error) {
	// Validate coordinates
	if !utils.ValidateCoordinates(req.PickupLocation.Latitude, req.PickupLocation.Longitude) {
		return nil, errors.New("invalid pickup coordinates")
	}
	if !utils.ValidateCoordinates(req.DropoffLocation.Latitude, req.DropoffLocation.Longitude) {
		return nil, errors.New("invalid dropoff coordinates")
	}

	// Calculate fare using pricing service
	distance, fare, duration := s.pricingService.EstimateFare(
		req.PickupLocation.Latitude,
		req.PickupLocation.Longitude,
		req.DropoffLocation.Latitude,
		req.DropoffLocation.Longitude,
		req.ServiceType,
	)

	// Set search started time
	now := time.Now()

	// Create trip with searching status
	trip := &models.Trip{
		CustomerID:        customerID,
		ServiceType:       models.ServiceType(req.ServiceType),
		Status:            models.TripStatusSearching, // Changed from TripStatusPending
		PickupLatitude:    req.PickupLocation.Latitude,
		PickupLongitude:   req.PickupLocation.Longitude,
		DropoffLatitude:   req.DropoffLocation.Latitude,
		DropoffLongitude:  req.DropoffLocation.Longitude,
		EstimatedDistance: &distance,
		EstimatedDuration: utils.Float64ToIntPointer(duration),
		FareAmount:        &fare,
		PaymentMethod:     "cash",
		SearchStartedAt:   &now,
	}

	// Set addresses if provided
	if req.PickupLocation.Address != "" {
		trip.PickupAddress = &req.PickupLocation.Address
	}
	if req.DropoffLocation.Address != "" {
		trip.DropoffAddress = &req.DropoffLocation.Address
	}

	if req.PaymentMethod != "" {
		trip.PaymentMethod = req.PaymentMethod
	}

	if err := s.tripRepo.Create(trip); err != nil {
		return nil, errors.New("failed to create trip")
	}

	// Create a job for drivers to accept
	job := &models.Job{
		TripID: trip.ID,
		Status: models.JobStatusPending,
	}
	if err := s.jobRepo.Create(job); err != nil {
		// Log error but don't fail trip creation
		// In production, you might want to handle this differently
	}

	return s.tripToDTO(trip), nil
}

// EstimateFare estimates the fare for a trip without creating it
func (s *tripService) EstimateFare(req dto.EstimateFareRequest) (*dto.EstimateFareResponse, error) {
	// Validate coordinates
	if !utils.ValidateCoordinates(req.PickupLocation.Latitude, req.PickupLocation.Longitude) {
		return nil, errors.New("invalid pickup coordinates")
	}
	if !utils.ValidateCoordinates(req.DropoffLocation.Latitude, req.DropoffLocation.Longitude) {
		return nil, errors.New("invalid dropoff coordinates")
	}

	// Calculate fare
	distance, fare, duration := s.pricingService.EstimateFare(
		req.PickupLocation.Latitude,
		req.PickupLocation.Longitude,
		req.DropoffLocation.Latitude,
		req.DropoffLocation.Longitude,
		req.ServiceType,
	)

	return &dto.EstimateFareResponse{
		Distance:          distance,
		EstimatedDuration: duration,
		EstimatedFare:     fare,
		ServiceType:       req.ServiceType,
	}, nil
}

func (s *tripService) GetActiveTrip(customerID uuid.UUID) (*dto.TripResponse, error) {
	trip, err := s.tripRepo.FindActiveByCustomerID(customerID)
	if err != nil {
		return nil, errors.New("no active trip found")
	}
	return s.tripToDTO(trip), nil
}

func (s *tripService) GetTripHistory(customerID uuid.UUID, limit, offset int) ([]dto.TripResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	trips, err := s.tripRepo.FindHistoryByCustomerID(customerID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch trip history")
	}

	responses := make([]dto.TripResponse, len(trips))
	for i, trip := range trips {
		responses[i] = *s.tripToDTO(&trip)
	}

	return responses, nil
}

func (s *tripService) GetTripByID(tripID uuid.UUID) (*dto.TripResponse, error) {
	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		return nil, errors.New("trip not found")
	}
	return s.tripToDTO(trip), nil
}

func (s *tripService) UpdateTrip(tripID uuid.UUID, req dto.UpdateTripRequest) (*dto.TripResponse, error) {
	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		return nil, errors.New("trip not found")
	}

	if req.Status != nil {
		trip.Status = models.TripStatus(*req.Status)
	}
	if req.EstimatedArrival != nil {
		arrival := time.Unix(*req.EstimatedArrival, 0)
		trip.EstimatedArrival = &arrival
	}
	if req.PickupLocation != nil {
		trip.PickupLatitude = req.PickupLocation.Latitude
		trip.PickupLongitude = req.PickupLocation.Longitude
	}
	if req.DropoffLocation != nil {
		trip.DropoffLatitude = req.DropoffLocation.Latitude
		trip.DropoffLongitude = req.DropoffLocation.Longitude
	}

	if err := s.tripRepo.Update(trip); err != nil {
		return nil, errors.New("failed to update trip")
	}

	return s.tripToDTO(trip), nil
}

func (s *tripService) CancelTrip(tripID uuid.UUID, customerID uuid.UUID) error {
	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		return errors.New("trip not found")
	}

	if trip.CustomerID != customerID {
		return errors.New("unauthorized to cancel this trip")
	}

	if trip.Status == models.TripStatusCompleted || trip.Status == models.TripStatusCancelled {
		return errors.New("trip cannot be cancelled")
	}

	trip.Status = models.TripStatusCancelled
	return s.tripRepo.Update(trip)
}

func (s *tripService) tripToDTO(trip *models.Trip) *dto.TripResponse {
	response := &dto.TripResponse{
		ID:          trip.ID.String(),
		CustomerID:  trip.CustomerID.String(),
		ServiceType: string(trip.ServiceType),
		Status:      string(trip.Status),
		PickupLocation: dto.Location{
			Latitude:  trip.PickupLatitude,
			Longitude: trip.PickupLongitude,
		},
		DropoffLocation: dto.Location{
			Latitude:  trip.DropoffLatitude,
			Longitude: trip.DropoffLongitude,
		},
		EstimatedDistance: trip.EstimatedDistance,
		EstimatedDuration: trip.EstimatedDuration,
		FareAmount:        trip.FareAmount,
		PaymentMethod:     trip.PaymentMethod,
		CreatedAt:         trip.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         trip.UpdatedAt.Format(time.RFC3339),
	}

	if trip.DriverID != nil {
		driverID := trip.DriverID.String()
		response.DriverID = &driverID
	}

	if trip.EstimatedArrival != nil {
		timestamp := trip.EstimatedArrival.Unix()
		response.EstimatedArrival = &timestamp
	}

	return response
}
// AcceptTrip allows a driver to accept a trip
func (s *tripService) AcceptTrip(tripID uuid.UUID, driverID uuid.UUID) error {
	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		return errors.New("trip not found")
	}

	// Check if trip is still searching
	if trip.Status != models.TripStatusSearching {
		return errors.New("trip is no longer available")
	}

	// Check if trip already has a driver (race condition)
	if trip.DriverID != nil {
		return errors.New("trip already accepted by another driver")
	}

	// Update trip
	now := time.Now()
	trip.Status = models.TripStatusAccepted
	trip.DriverID = &driverID
	trip.SearchEndedAt = &now

	if err := s.tripRepo.Update(trip); err != nil {
		return errors.New("failed to accept trip")
	}

	return nil
}

// GetTripStatus returns the current status of a trip
func (s *tripService) GetTripStatus(tripID uuid.UUID) (string, error) {
	trip, err := s.tripRepo.FindByID(tripID)
	if err != nil {
		return "", errors.New("trip not found")
	}
	return string(trip.Status), nil
}

// ExpireSearchingTrips marks trips that have been searching for >3 minutes as expired
func (s *tripService) ExpireSearchingTrips() error {
	// Find all trips in searching status older than 3 minutes
	threeMinutesAgo := time.Now().Add(-3 * time.Minute)

	trips, err := s.tripRepo.FindSearchingBefore(threeMinutesAgo)
	if err != nil {
		return err
	}

	now := time.Now()
	for _, trip := range trips {
		trip.Status = models.TripStatusExpired
		trip.SearchEndedAt = &now
		if err := s.tripRepo.Update(&trip); err != nil {
			// Log error but continue with other trips
			continue
		}
	}

	return nil
}
