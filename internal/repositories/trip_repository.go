package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type TripRepository interface {
	Create(trip *models.Trip) error
	FindByID(id uuid.UUID) (*models.Trip, error)
	FindActiveByCustomerID(customerID uuid.UUID) (*models.Trip, error)
	FindHistoryByCustomerID(customerID uuid.UUID, limit, offset int) ([]models.Trip, error)
	FindByDriverID(driverID uuid.UUID) ([]models.Trip, error)
	Update(trip *models.Trip) error
	Delete(id uuid.UUID) error
	FindPendingTrips() ([]models.Trip, error)
	FindSearchingBefore(time time.Time) ([]models.Trip, error)
}

type tripRepository struct {
	db *gorm.DB
}

func NewTripRepository() TripRepository {
	return &tripRepository{
		db: database.DB,
	}
}

func (r *tripRepository) Create(trip *models.Trip) error {
	return r.db.Create(trip).Error
}

func (r *tripRepository) FindByID(id uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	err := r.db.Preload("Customer").Preload("Driver").Where("id = ?", id).First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) FindActiveByCustomerID(customerID uuid.UUID) (*models.Trip, error) {
	var trip models.Trip
	err := r.db.Where("customer_id = ? AND status IN ?", customerID, []string{"pending", "accepted", "in_progress"}).
		Preload("Customer").Preload("Driver").
		Order("created_at DESC").First(&trip).Error
	if err != nil {
		return nil, err
	}
	return &trip, nil
}

func (r *tripRepository) FindHistoryByCustomerID(customerID uuid.UUID, limit, offset int) ([]models.Trip, error) {
	var trips []models.Trip
	err := r.db.Where("customer_id = ?", customerID).
		Preload("Customer").Preload("Driver").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&trips).Error
	return trips, err
}

func (r *tripRepository) FindByDriverID(driverID uuid.UUID) ([]models.Trip, error) {
	var trips []models.Trip
	err := r.db.Where("driver_id = ?", driverID).
		Preload("Customer").Preload("Driver").
		Order("created_at DESC").
		Find(&trips).Error
	return trips, err
}

func (r *tripRepository) Update(trip *models.Trip) error {
	return r.db.Save(trip).Error
}

func (r *tripRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Trip{}, id).Error
}

func (r *tripRepository) FindPendingTrips() ([]models.Trip, error) {
	var trips []models.Trip
	err := r.db.Where("status = ?", models.TripStatusPending).
		Preload("Customer").
		Order("created_at DESC").
		Find(&trips).Error
	return trips, err
}

// FindSearchingBefore finds all trips in searching status created before the given time
func (r *tripRepository) FindSearchingBefore(before time.Time) ([]models.Trip, error) {
	var trips []models.Trip
	err := r.db.Where("status = ? AND search_started_at < ?", models.TripStatusSearching, before).
		Find(&trips).Error
	return trips, err
}
