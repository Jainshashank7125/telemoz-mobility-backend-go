package repositories

import (
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type NotificationRepository interface {
	Create(notification *models.Notification) error
	FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error)
	FindUnreadByUserID(userID uuid.UUID) ([]models.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
	Delete(notificationID uuid.UUID) error
}

type NotificationSettingsRepository interface {
	Create(settings *models.NotificationSettings) error
	FindByUserID(userID uuid.UUID) (*models.NotificationSettings, error)
	Update(settings *models.NotificationSettings) error
}

type notificationRepository struct {
	db *gorm.DB
}

func NewNotificationRepository() NotificationRepository {
	return &notificationRepository{
		db: database.DB,
	}
}

func (r *notificationRepository) Create(notification *models.Notification) error {
	return r.db.Create(notification).Error
}

func (r *notificationRepository) FindByUserID(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ?", userID).
		Order("created_at DESC").
		Limit(limit).Offset(offset).
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) FindUnreadByUserID(userID uuid.UUID) ([]models.Notification, error) {
	var notifications []models.Notification
	err := r.db.Where("user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Find(&notifications).Error
	return notifications, err
}

func (r *notificationRepository) MarkAsRead(notificationID uuid.UUID) error {
	return r.db.Model(&models.Notification{}).
		Where("id = ?", notificationID).
		Update("is_read", true).Error
}

func (r *notificationRepository) Delete(notificationID uuid.UUID) error {
	return r.db.Delete(&models.Notification{}, notificationID).Error
}

type notificationSettingsRepository struct {
	db *gorm.DB
}

func NewNotificationSettingsRepository() NotificationSettingsRepository {
	return &notificationSettingsRepository{
		db: database.DB,
	}
}

func (r *notificationSettingsRepository) Create(settings *models.NotificationSettings) error {
	return r.db.Create(settings).Error
}

func (r *notificationSettingsRepository) FindByUserID(userID uuid.UUID) (*models.NotificationSettings, error) {
	var settings models.NotificationSettings
	err := r.db.Where("user_id = ?", userID).First(&settings).Error
	if err != nil {
		return nil, err
	}
	return &settings, nil
}

func (r *notificationSettingsRepository) Update(settings *models.NotificationSettings) error {
	return r.db.Save(settings).Error
}

