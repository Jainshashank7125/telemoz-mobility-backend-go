package services

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
)

type JobService interface {
	GetAvailableJobs() ([]dto.JobResponse, error)
	AcceptJob(jobID, driverID uuid.UUID) (*dto.JobResponse, error)
	RejectJob(jobID, driverID uuid.UUID) error
	GetActiveJob(driverID uuid.UUID) (*dto.JobResponse, error)
	GetJobHistory(driverID uuid.UUID, limit, offset int) ([]dto.JobResponse, error)
	UpdateJobStatus(jobID, driverID uuid.UUID, status string) (*dto.JobResponse, error)
}

type jobService struct {
	jobRepo  repositories.JobRepository
	tripRepo repositories.TripRepository
}

func NewJobService() JobService {
	return &jobService{
		jobRepo:  repositories.NewJobRepository(),
		tripRepo: repositories.NewTripRepository(),
	}
}

func (s *jobService) GetAvailableJobs() ([]dto.JobResponse, error) {
	jobs, err := s.jobRepo.FindAvailable()
	if err != nil {
		return nil, errors.New("failed to fetch available jobs")
	}

	responses := make([]dto.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = *s.jobToDTO(&job)
	}

	return responses, nil
}

func (s *jobService) AcceptJob(jobID, driverID uuid.UUID) (*dto.JobResponse, error) {
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if job.Status != models.JobStatusPending {
		return nil, errors.New("job is not available")
	}

	now := time.Now()
	job.DriverID = &driverID
	job.Status = models.JobStatusAccepted
	job.AcceptedAt = &now

	// Update trip
	trip, err := s.tripRepo.FindByID(job.TripID)
	if err != nil {
		return nil, errors.New("trip not found")
	}

	trip.DriverID = &driverID
	trip.Status = models.TripStatusAccepted

	if err := s.tripRepo.Update(trip); err != nil {
		return nil, errors.New("failed to update trip")
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, errors.New("failed to accept job")
	}

	return s.jobToDTO(job), nil
}

func (s *jobService) RejectJob(jobID, driverID uuid.UUID) error {
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		return errors.New("job not found")
	}

	if job.Status != models.JobStatusPending {
		return errors.New("job cannot be rejected")
	}

	job.Status = models.JobStatusRejected
	return s.jobRepo.Update(job)
}

func (s *jobService) GetActiveJob(driverID uuid.UUID) (*dto.JobResponse, error) {
	job, err := s.jobRepo.FindActiveByDriverID(driverID)
	if err != nil {
		return nil, errors.New("no active job found")
	}
	return s.jobToDTO(job), nil
}

func (s *jobService) GetJobHistory(driverID uuid.UUID, limit, offset int) ([]dto.JobResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}

	jobs, err := s.jobRepo.FindHistoryByDriverID(driverID, limit, offset)
	if err != nil {
		return nil, errors.New("failed to fetch job history")
	}

	responses := make([]dto.JobResponse, len(jobs))
	for i, job := range jobs {
		responses[i] = *s.jobToDTO(&job)
	}

	return responses, nil
}

func (s *jobService) UpdateJobStatus(jobID, driverID uuid.UUID, status string) (*dto.JobResponse, error) {
	job, err := s.jobRepo.FindByID(jobID)
	if err != nil {
		return nil, errors.New("job not found")
	}

	if job.DriverID == nil || *job.DriverID != driverID {
		return nil, errors.New("unauthorized to update this job")
	}

	job.Status = models.JobStatus(status)
	if status == "completed" {
		now := time.Now()
		job.CompletedAt = &now

		// Update trip status
		trip, err := s.tripRepo.FindByID(job.TripID)
		if err == nil {
			trip.Status = models.TripStatusCompleted
			s.tripRepo.Update(trip)
		}
	}

	if err := s.jobRepo.Update(job); err != nil {
		return nil, errors.New("failed to update job status")
	}

	return s.jobToDTO(job), nil
}

func (s *jobService) jobToDTO(job *models.Job) *dto.JobResponse {
	response := &dto.JobResponse{
		ID:              job.ID.String(),
		TripID:          job.TripID.String(),
		Status:          string(job.Status),
		EstimatedEarnings: job.EstimatedEarnings,
		Distance:        job.Distance,
		CreatedAt:       job.CreatedAt.Format(time.RFC3339),
	}

	if job.DriverID != nil {
		driverID := job.DriverID.String()
		response.DriverID = &driverID
	}

	// Include trip details
	if job.Trip.ID != uuid.Nil {
		response.ServiceType = string(job.Trip.ServiceType)
		response.CustomerID = job.Trip.CustomerID.String()
		response.PickupLocation = dto.Location{
			Latitude:  job.Trip.PickupLatitude,
			Longitude: job.Trip.PickupLongitude,
		}
		response.DropoffLocation = dto.Location{
			Latitude:  job.Trip.DropoffLatitude,
			Longitude: job.Trip.DropoffLongitude,
		}
	}

	return response
}

