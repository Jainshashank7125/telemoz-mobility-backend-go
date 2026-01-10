package repositories

import (
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type RefreshTokenRepository interface {
	Create(token *models.RefreshToken) error
	FindByToken(token string) (*models.RefreshToken, error)
	DeleteByToken(token string) error
	DeleteByUserID(userID uuid.UUID) error
	DeleteExpired() error
}

type refreshTokenRepository struct {
	db *gorm.DB
}

func NewRefreshTokenRepository() RefreshTokenRepository {
	return &refreshTokenRepository{
		db: database.DB,
	}
}

func (r *refreshTokenRepository) Create(token *models.RefreshToken) error {
	return r.db.Create(token).Error
}

func (r *refreshTokenRepository) FindByToken(token string) (*models.RefreshToken, error) {
	var refreshToken models.RefreshToken
	err := r.db.Where("token = ? AND expires_at > ?", token, time.Now()).
		Preload("User").
		First(&refreshToken).Error
	if err != nil {
		return nil, err
	}
	return &refreshToken, nil
}

func (r *refreshTokenRepository) DeleteByToken(token string) error {
	return r.db.Where("token = ?", token).Delete(&models.RefreshToken{}).Error
}

func (r *refreshTokenRepository) DeleteByUserID(userID uuid.UUID) error {
	return r.db.Where("user_id = ?", userID).Delete(&models.RefreshToken{}).Error
}

func (r *refreshTokenRepository) DeleteExpired() error {
	return r.db.Where("expires_at < ?", time.Now()).Delete(&models.RefreshToken{}).Error
}

