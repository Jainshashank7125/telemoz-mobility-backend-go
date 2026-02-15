package repositories

import (
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type DriverAvailabilityRepository interface {
	Create(availability *models.DriverAvailability) error
	Update(availability *models.DriverAvailability) error
	FindByDriverID(driverID uuid.UUID) (*models.DriverAvailability, error)
	FindAvailableByServiceType(serviceType string) ([]uuid.UUID, error)
}

type driverAvailabilityRepository struct {
	db *gorm.DB
}

func NewDriverAvailabilityRepository() DriverAvailabilityRepository {
	return &driverAvailabilityRepository{
		db: database.DB,
	}
}

func (r *driverAvailabilityRepository) Create(availability *models.DriverAvailability) error {
	return r.db.Create(availability).Error
}

func (r *driverAvailabilityRepository) Update(availability *models.DriverAvailability) error {
	return r.db.Save(availability).Error
}

func (r *driverAvailabilityRepository) FindByDriverID(driverID uuid.UUID) (*models.DriverAvailability, error) {
	var availability models.DriverAvailability
	err := r.db.Where("driver_id = ?", driverID).First(&availability).Error
	if err != nil {
		return nil, err
	}
	return &availability, nil
}

func (r *driverAvailabilityRepository) FindAvailableByServiceType(serviceType string) ([]uuid.UUID, error) {
	var driverIDs []uuid.UUID
	err := r.db.Model(&models.DriverAvailability{}).
		Where("is_available = ? AND ? = ANY(service_types)", true, serviceType).
		Pluck("driver_id", &driverIDs).Error
	return driverIDs, err
}
