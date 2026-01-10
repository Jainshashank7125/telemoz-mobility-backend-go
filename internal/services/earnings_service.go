package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/database"
	"github.com/telemoz/backend/internal/models"
	"gorm.io/gorm"
)

type EarningsService interface {
	GetSummary(driverID uuid.UUID) (map[string]float64, error)
	GetHistory(driverID uuid.UUID, limit, offset int) ([]models.DriverEarning, error)
}

type earningsService struct {
	db *gorm.DB
}

func NewEarningsService() EarningsService {
	return &earningsService{
		db: database.DB,
	}
}

func (s *earningsService) GetSummary(driverID uuid.UUID) (map[string]float64, error) {
	now := time.Now()
	todayStart := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
	weekStart := todayStart.AddDate(0, 0, -int(now.Weekday()))
	monthStart := time.Date(now.Year(), now.Month(), 1, 0, 0, 0, 0, now.Location())

	summary := make(map[string]float64)

	// Today's earnings
	var todayTotal float64
	err := s.db.Model(&models.DriverEarning{}).
		Where("driver_id = ? AND date = ?", driverID, todayStart.Format("2006-01-02")).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&todayTotal).Error
	if err != nil {
		return nil, errors.New("failed to calculate today's earnings")
	}
	summary["today"] = todayTotal

	// This week's earnings
	var weekTotal float64
	err = s.db.Model(&models.DriverEarning{}).
		Where("driver_id = ? AND date >= ?", driverID, weekStart.Format("2006-01-02")).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&weekTotal).Error
	if err != nil {
		return nil, errors.New("failed to calculate week's earnings")
	}
	summary["week"] = weekTotal

	// This month's earnings
	var monthTotal float64
	err = s.db.Model(&models.DriverEarning{}).
		Where("driver_id = ? AND date >= ?", driverID, monthStart.Format("2006-01-02")).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&monthTotal).Error
	if err != nil {
		return nil, errors.New("failed to calculate month's earnings")
	}
	summary["month"] = monthTotal

	// Total earnings
	var total float64
	err = s.db.Model(&models.DriverEarning{}).
		Where("driver_id = ?", driverID).
		Select("COALESCE(SUM(amount), 0)").
		Scan(&total).Error
	if err != nil {
		return nil, errors.New("failed to calculate total earnings")
	}
	summary["total"] = total

	return summary, nil
}

func (s *earningsService) GetHistory(driverID uuid.UUID, limit, offset int) ([]models.DriverEarning, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	var earnings []models.DriverEarning
	err := s.db.Where("driver_id = ?", driverID).
		Preload("Job").Preload("Job.Trip").
		Order("date DESC, created_at DESC").
		Limit(limit).Offset(offset).
		Find(&earnings).Error

	if err != nil {
		return nil, errors.New("failed to fetch earnings history")
	}

	return earnings, nil
}

