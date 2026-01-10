package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
)

type NotificationService interface {
	CreateNotification(userID uuid.UUID, notificationType, title, message string, data map[string]interface{}) error
	ListNotifications(userID uuid.UUID, limit, offset int) ([]models.Notification, error)
	MarkAsRead(notificationID uuid.UUID) error
	GetSettings(userID uuid.UUID) (*models.NotificationSettings, error)
	UpdateSettings(userID uuid.UUID, settings *models.NotificationSettings) error
}

type notificationService struct {
	notificationRepo         repositories.NotificationRepository
	notificationSettingsRepo repositories.NotificationSettingsRepository
}

func NewNotificationService() NotificationService {
	return &notificationService{
		notificationRepo:         repositories.NewNotificationRepository(),
		notificationSettingsRepo: repositories.NewNotificationSettingsRepository(),
	}
}

func (s *notificationService) CreateNotification(userID uuid.UUID, notificationType, title, message string, data map[string]interface{}) error {
	notification := &models.Notification{
		UserID:  userID,
		Type:    notificationType,
		Title:   title,
		Message: message,
		IsRead:  false,
	}

	if data != nil {
		notification.Data = models.JSONB(data)
	}

	return s.notificationRepo.Create(notification)
}

func (s *notificationService) ListNotifications(userID uuid.UUID, limit, offset int) ([]models.Notification, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	notifications, err := s.notificationRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch notifications")
	}
	return notifications, nil
}

func (s *notificationService) MarkAsRead(notificationID uuid.UUID) error {
	return s.notificationRepo.MarkAsRead(notificationID)
}

func (s *notificationService) GetSettings(userID uuid.UUID) (*models.NotificationSettings, error) {
	settings, err := s.notificationSettingsRepo.FindByUserID(userID)
	if err != nil {
		// Create default settings if not found
		defaultSettings := &models.NotificationSettings{
			UserID:         userID,
			BusNearbyAlert: true,
			BusArrived:     true,
			BusDeparted:    false,
			RouteChange:    true,
			SMSEnabled:     false,
			CallEnabled:    false,
		}
		if err := s.notificationSettingsRepo.Create(defaultSettings); err != nil {
			return nil, errors.New("failed to create default settings")
		}
		return defaultSettings, nil
	}
	return settings, nil
}

func (s *notificationService) UpdateSettings(userID uuid.UUID, settings *models.NotificationSettings) error {
	existing, err := s.notificationSettingsRepo.FindByUserID(userID)
	if err != nil {
		// Create if doesn't exist
		settings.UserID = userID
		return s.notificationSettingsRepo.Create(settings)
	}

	existing.BusNearbyAlert = settings.BusNearbyAlert
	existing.BusArrived = settings.BusArrived
	existing.BusDeparted = settings.BusDeparted
	existing.RouteChange = settings.RouteChange
	existing.SMSEnabled = settings.SMSEnabled
	existing.CallEnabled = settings.CallEnabled

	return s.notificationSettingsRepo.Update(existing)
}

