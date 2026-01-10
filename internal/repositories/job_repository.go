package repositories

import (
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type JobRepository interface {
	Create(job *models.Job) error
	FindByID(id uuid.UUID) (*models.Job, error)
	FindAvailable() ([]models.Job, error)
	FindActiveByDriverID(driverID uuid.UUID) (*models.Job, error)
	FindHistoryByDriverID(driverID uuid.UUID, limit, offset int) ([]models.Job, error)
	Update(job *models.Job) error
	Delete(id uuid.UUID) error
}

type jobRepository struct {
	db *gorm.DB
}

func NewJobRepository() JobRepository {
	return &jobRepository{
		db: database.DB,
	}
}

func (r *jobRepository) Create(job *models.Job) error {
	return r.db.Create(job).Error
}

func (r *jobRepository) FindByID(id uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := r.db.Preload("Trip").Preload("Trip.Customer").Preload("Driver").
		Where("id = ?", id).First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) FindAvailable() ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Where("status = ?", models.JobStatusPending).
		Preload("Trip").Preload("Trip.Customer").
		Order("created_at DESC").
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) FindActiveByDriverID(driverID uuid.UUID) (*models.Job, error) {
	var job models.Job
	err := r.db.Where("driver_id = ? AND status IN ?", driverID, []string{"accepted", "in_progress"}).
		Preload("Trip").Preload("Trip.Customer").Preload("Driver").
		Order("created_at DESC").First(&job).Error
	if err != nil {
		return nil, err
	}
	return &job, nil
}

func (r *jobRepository) FindHistoryByDriverID(driverID uuid.UUID, limit, offset int) ([]models.Job, error) {
	var jobs []models.Job
	err := r.db.Where("driver_id = ?", driverID).
		Preload("Trip").Preload("Trip.Customer").
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&jobs).Error
	return jobs, err
}

func (r *jobRepository) Update(job *models.Job) error {
	return r.db.Save(job).Error
}

func (r *jobRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Job{}, id).Error
}

