package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type EarningsHandler struct {
	earningsService services.EarningsService
}

func NewEarningsHandler() *EarningsHandler {
	return &EarningsHandler{
		earningsService: services.NewEarningsService(),
	}
}

// GetSummary gets earnings summary
func (h *EarningsHandler) GetSummary(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	summary, err := h.earningsService.GetSummary(driverID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, summary, "Earnings summary retrieved successfully")
}

// GetHistory gets earnings history
func (h *EarningsHandler) GetHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	driverID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	earnings, err := h.earningsService.GetHistory(driverID, limit, offset)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, earnings, "Earnings history retrieved successfully")
}

