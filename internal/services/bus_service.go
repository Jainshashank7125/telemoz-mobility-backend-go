package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
)

type BusService interface {
	GetBusByChildID(childID uuid.UUID) (*models.Bus, error)
	GetBusLocation(busID uuid.UUID) (*models.BusLocation, error)
	UpdateBusLocation(busID uuid.UUID, lat, lng, accuracy, speed, heading float64) error
}

type busService struct {
	busRepo         repositories.BusRepository
	busLocationRepo repositories.BusLocationRepository
	childRepo       repositories.ChildRepository
}

func NewBusService() BusService {
	return &busService{
		busRepo:         repositories.NewBusRepository(),
		busLocationRepo: repositories.NewBusLocationRepository(),
		childRepo:       repositories.NewChildRepository(),
	}
}

func (s *busService) GetBusByChildID(childID uuid.UUID) (*models.Bus, error) {
	bus, err := s.busRepo.FindByChildID(childID)
	if err != nil {
		return nil, errors.New("bus not found for child")
	}
	return bus, nil
}

func (s *busService) GetBusLocation(busID uuid.UUID) (*models.BusLocation, error) {
	location, err := s.busLocationRepo.FindLatestByBusID(busID)
	if err != nil {
		return nil, errors.New("bus location not found")
	}
	return location, nil
}

func (s *busService) UpdateBusLocation(busID uuid.UUID, lat, lng, accuracy, speed, heading float64) error {
	location := &models.BusLocation{
		BusID:     busID,
		Latitude:  lat,
		Longitude: lng,
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

