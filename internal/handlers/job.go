package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type JobHandler struct {
	jobService services.JobService
}

func NewJobHandler() *JobHandler {
	return &JobHandler{
		jobService: services.NewJobService(),
	}
}

// GetAvailableJobs gets available jobs for drivers
func (h *JobHandler) GetAvailableJobs(c *gin.Context) {
	jobs, err := h.jobService.GetAvailableJobs()
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, jobs, "Available jobs retrieved successfully")
}

// AcceptJob accepts a job
func (h *JobHandler) AcceptJob(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid job ID", nil)
		return
	}

	job, err := h.jobService.AcceptJob(jobID, driverID)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, job, "Job accepted successfully")
}

// RejectJob rejects a job
func (h *JobHandler) RejectJob(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid job ID", nil)
		return
	}

	if err := h.jobService.RejectJob(jobID, driverID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Job rejected successfully")
}

// GetActiveJob gets the active job for the driver
func (h *JobHandler) GetActiveJob(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	job, err := h.jobService.GetActiveJob(driverID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, job, "Active job retrieved successfully")
}

// GetJobHistory gets job history for the driver
func (h *JobHandler) GetJobHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	jobs, err := h.jobService.GetJobHistory(driverID, limit, offset)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, jobs, "Job history retrieved successfully")
}

// UpdateJobStatus updates job status
func (h *JobHandler) UpdateJobStatus(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	jobID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid job ID", nil)
		return
	}

	var req dto.UpdateJobStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	job, err := h.jobService.UpdateJobStatus(jobID, driverID, req.Status)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, job, "Job status updated successfully")
}

