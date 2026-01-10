package handlers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type NotificationHandler struct {
	notificationService services.NotificationService
}

func NewNotificationHandler() *NotificationHandler {
	return &NotificationHandler{
		notificationService: services.NewNotificationService(),
	}
}

// ListNotifications lists notifications for the user
func (h *NotificationHandler) ListNotifications(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	notifications, err := h.notificationService.ListNotifications(userID, limit, offset)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, notifications, "Notifications retrieved successfully")
}

// MarkAsRead marks a notification as read
func (h *NotificationHandler) MarkAsRead(c *gin.Context) {
	notificationID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid notification ID", nil)
		return
	}

	if err := h.notificationService.MarkAsRead(notificationID); err != nil {
		utils.InternalError(c, "Failed to mark notification as read")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Notification marked as read")
}

// GetSettings gets notification settings
func (h *NotificationHandler) GetSettings(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	settings, err := h.notificationService.GetSettings(userID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, settings, "Notification settings retrieved successfully")
}

// UpdateSettings updates notification settings
func (h *NotificationHandler) UpdateSettings(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	var settings models.NotificationSettings
	if err := c.ShouldBindJSON(&settings); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	if err := h.notificationService.UpdateSettings(userID, &settings); err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, settings, "Notification settings updated successfully")
}

