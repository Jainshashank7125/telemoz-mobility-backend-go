package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type BusRepository interface {
	Create(bus *models.Bus) error
	FindByID(id uuid.UUID) (*models.Bus, error)
	FindByChildID(childID uuid.UUID) (*models.Bus, error)
	FindByTraccarDeviceID(deviceID string) (*models.Bus, error)
	Update(bus *models.Bus) error
	Delete(id uuid.UUID) error
}

type BusLocationRepository interface {
	Create(location *models.BusLocation) error
	FindLatestByBusID(busID uuid.UUID) (*models.BusLocation, error)
	FindByBusIDAndTimeRange(busID uuid.UUID, start, end time.Time) ([]models.BusLocation, error)
}

type busRepository struct {
	db *gorm.DB
}

func NewBusRepository() BusRepository {
	return &busRepository{
		db: database.DB,
	}
}

func (r *busRepository) Create(bus *models.Bus) error {
	return r.db.Create(bus).Error
}

func (r *busRepository) FindByID(id uuid.UUID) (*models.Bus, error) {
	var bus models.Bus
	err := r.db.Preload("Driver").Preload("Children").Where("id = ?", id).First(&bus).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *busRepository) FindByChildID(childID uuid.UUID) (*models.Bus, error) {
	var child models.Child
	err := r.db.Where("id = ?", childID).First(&child).Error
	if err != nil {
		return nil, err
	}
	if child.BusID == nil {
		return nil, gorm.ErrRecordNotFound
	}
	return r.FindByID(*child.BusID)
}

func (r *busRepository) FindByTraccarDeviceID(deviceID string) (*models.Bus, error) {
	var bus models.Bus
	err := r.db.Where("traccar_device_id = ?", deviceID).
		Preload("Driver").
		First(&bus).Error
	if err != nil {
		return nil, err
	}
	return &bus, nil
}

func (r *busRepository) Update(bus *models.Bus) error {
	return r.db.Save(bus).Error
}

func (r *busRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Bus{}, id).Error
}

type busLocationRepository struct {
	db *gorm.DB
}

func NewBusLocationRepository() BusLocationRepository {
	return &busLocationRepository{
		db: database.DB,
	}
}

func (r *busLocationRepository) Create(location *models.BusLocation) error {
	return r.db.Create(location).Error
}

func (r *busLocationRepository) FindLatestByBusID(busID uuid.UUID) (*models.BusLocation, error) {
	var location models.BusLocation
	err := r.db.Where("bus_id = ?", busID).
		Order("timestamp DESC").
		First(&location).Error
	if err != nil {
		return nil, err
	}
	return &location, nil
}

func (r *busLocationRepository) FindByBusIDAndTimeRange(busID uuid.UUID, start, end time.Time) ([]models.BusLocation, error) {
	var locations []models.BusLocation
	err := r.db.Where("bus_id = ? AND timestamp BETWEEN ? AND ?", busID, start, end).
		Order("timestamp ASC").
		Find(&locations).Error
	return locations, err
}

