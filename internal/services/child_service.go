package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
	"github.com/telemoz/backend/internal/utils"
)

type ChildService interface {
	CreateChild(parentID uuid.UUID, name, schoolName string, busID *uuid.UUID) (*models.Child, error)
	GetChildByID(childID uuid.UUID) (*models.Child, error)
	ListChildren(parentID uuid.UUID) ([]models.Child, error)
	UpdateChild(childID, parentID uuid.UUID, name, schoolName *string, busID *uuid.UUID) (*models.Child, error)
	DeleteChild(childID, parentID uuid.UUID) error
}

type childService struct {
	childRepo repositories.ChildRepository
}

func NewChildService() ChildService {
	return &childService{
		childRepo: repositories.NewChildRepository(),
	}
}

func (s *childService) CreateChild(parentID uuid.UUID, name, schoolName string, busID *uuid.UUID) (*models.Child, error) {
	child := &models.Child{
		ParentID: parentID,
		Name:     utils.SanitizeString(name),
		BusID:    busID,
	}

	if schoolName != "" {
		school := utils.SanitizeString(schoolName)
		child.SchoolName = &school
	}

	if err := s.childRepo.Create(child); err != nil {
		return nil, errors.New("failed to create child")
	}

	return child, nil
}

func (s *childService) GetChildByID(childID uuid.UUID) (*models.Child, error) {
	child, err := s.childRepo.FindByID(childID)
	if err != nil {
		return nil, errors.New("child not found")
	}
	return child, nil
}

func (s *childService) ListChildren(parentID uuid.UUID) ([]models.Child, error) {
	children, err := s.childRepo.FindByParentID(parentID)
	if err != nil {
		return nil, errors.New("failed to fetch children")
	}
	return children, nil
}

func (s *childService) UpdateChild(childID, parentID uuid.UUID, name, schoolName *string, busID *uuid.UUID) (*models.Child, error) {
	child, err := s.childRepo.FindByID(childID)
	if err != nil {
		return nil, errors.New("child not found")
	}

	if child.ParentID != parentID {
		return nil, errors.New("unauthorized to update this child")
	}

	if name != nil {
		child.Name = utils.SanitizeString(*name)
	}
	if schoolName != nil {
		school := utils.SanitizeString(*schoolName)
		child.SchoolName = &school
	}
	if busID != nil {
		child.BusID = busID
	}

	if err := s.childRepo.Update(child); err != nil {
		return nil, errors.New("failed to update child")
	}

	return child, nil
}

func (s *childService) DeleteChild(childID, parentID uuid.UUID) error {
	child, err := s.childRepo.FindByID(childID)
	if err != nil {
		return errors.New("child not found")
	}

	if child.ParentID != parentID {
		return errors.New("unauthorized to delete this child")
	}

	return s.childRepo.Delete(childID)
}

