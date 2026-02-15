package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
)

type DriverAvailabilityService interface {
	UpdateAvailability(driverID uuid.UUID, isAvailable bool, serviceTypes []string) error
	GetAvailability(driverID uuid.UUID) (*models.DriverAvailability, error)
	GetAvailableDrivers(serviceType string) ([]uuid.UUID, error)
}

type driverAvailabilityService struct {
	availabilityRepo repositories.DriverAvailabilityRepository
}

func NewDriverAvailabilityService(repo repositories.DriverAvailabilityRepository) DriverAvailabilityService {
	return &driverAvailabilityService{
		availabilityRepo: repo,
	}
}

func (s *driverAvailabilityService) UpdateAvailability(
	driverID uuid.UUID,
	isAvailable bool,
	serviceTypes []string,
) error {
	availability, err := s.availabilityRepo.FindByDriverID(driverID)

	now := time.Now()

	if err != nil {
		// Create new record
		availability = &models.DriverAvailability{
			DriverID:     driverID,
			IsAvailable:  isAvailable,
			ServiceTypes: serviceTypes,
			LastActiveAt: now,
		}
		return s.availabilityRepo.Create(availability)
	}

	// Update existing
	availability.IsAvailable = isAvailable
	availability.ServiceTypes = serviceTypes
	availability.LastActiveAt = now

	return s.availabilityRepo.Update(availability)
}

func (s *driverAvailabilityService) GetAvailability(driverID uuid.UUID) (*models.DriverAvailability, error) {
	availability, err := s.availabilityRepo.FindByDriverID(driverID)
	if err != nil {
		return nil, errors.New("availability not found")
	}
	return availability, nil
}

func (s *driverAvailabilityService) GetAvailableDrivers(serviceType string) ([]uuid.UUID, error) {
	return s.availabilityRepo.FindAvailableByServiceType(serviceType)
}
