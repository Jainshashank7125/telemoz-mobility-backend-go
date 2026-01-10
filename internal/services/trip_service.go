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
	GetActiveTrip(customerID uuid.UUID) (*dto.TripResponse, error)
	GetTripHistory(customerID uuid.UUID, limit, offset int) ([]dto.TripResponse, error)
	GetTripByID(tripID uuid.UUID) (*dto.TripResponse, error)
	UpdateTrip(tripID uuid.UUID, req dto.UpdateTripRequest) (*dto.TripResponse, error)
	CancelTrip(tripID uuid.UUID, customerID uuid.UUID) error
}

type tripService struct {
	tripRepo repositories.TripRepository
	jobRepo  repositories.JobRepository
	userRepo repositories.UserRepository
}

func NewTripService() TripService {
	return &tripService{
		tripRepo: repositories.NewTripRepository(),
		jobRepo:  repositories.NewJobRepository(),
		userRepo: repositories.NewUserRepository(),
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

	// Create trip
	trip := &models.Trip{
		CustomerID:      customerID,
		ServiceType:     models.ServiceType(req.ServiceType),
		Status:          models.TripStatusPending,
		PickupLatitude:   req.PickupLocation.Latitude,
		PickupLongitude:  req.PickupLocation.Longitude,
		DropoffLatitude:  req.DropoffLocation.Latitude,
		DropoffLongitude: req.DropoffLocation.Longitude,
		PaymentMethod:    "cash",
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
		ID:              trip.ID.String(),
		CustomerID:      trip.CustomerID.String(),
		ServiceType:     string(trip.ServiceType),
		Status:          string(trip.Status),
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

