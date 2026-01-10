package repositories

import (
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type ChildRepository interface {
	Create(child *models.Child) error
	FindByID(id uuid.UUID) (*models.Child, error)
	FindByParentID(parentID uuid.UUID) ([]models.Child, error)
	Update(child *models.Child) error
	Delete(id uuid.UUID) error
}

type childRepository struct {
	db *gorm.DB
}

func NewChildRepository() ChildRepository {
	return &childRepository{
		db: database.DB,
	}
}

func (r *childRepository) Create(child *models.Child) error {
	return r.db.Create(child).Error
}

func (r *childRepository) FindByID(id uuid.UUID) (*models.Child, error) {
	var child models.Child
	err := r.db.Preload("Parent").Preload("Bus").Where("id = ?", id).First(&child).Error
	if err != nil {
		return nil, err
	}
	return &child, nil
}

func (r *childRepository) FindByParentID(parentID uuid.UUID) ([]models.Child, error) {
	var children []models.Child
	err := r.db.Where("parent_id = ?", parentID).
		Preload("Parent").Preload("Bus").
		Order("created_at DESC").
		Find(&children).Error
	return children, err
}

func (r *childRepository) Update(child *models.Child) error {
	return r.db.Save(child).Error
}

func (r *childRepository) Delete(id uuid.UUID) error {
	return r.db.Delete(&models.Child{}, id).Error
}

