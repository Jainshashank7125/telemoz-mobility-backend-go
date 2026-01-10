package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type ChildHandler struct {
	childService services.ChildService
}

func NewChildHandler() *ChildHandler {
	return &ChildHandler{
		childService: services.NewChildService(),
	}
}

// ListChildren lists all children for the parent
func (h *ChildHandler) ListChildren(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	parentID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	children, err := h.childService.ListChildren(parentID)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, children, "Children retrieved successfully")
}

// CreateChild creates a new child
func (h *ChildHandler) CreateChild(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	parentID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	var req struct {
		Name       string    `json:"name" binding:"required"`
		SchoolName string    `json:"school_name,omitempty"`
		BusID      *uuid.UUID `json:"bus_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	child, err := h.childService.CreateChild(parentID, req.Name, req.SchoolName, req.BusID)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, child, "Child created successfully")
}

// GetChildByID gets a child by ID
func (h *ChildHandler) GetChildByID(c *gin.Context) {
	childID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid child ID", nil)
		return
	}

	child, err := h.childService.GetChildByID(childID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, child, "Child retrieved successfully")
}

// UpdateChild updates a child
func (h *ChildHandler) UpdateChild(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	parentID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	childID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid child ID", nil)
		return
	}

	var req struct {
		Name       *string    `json:"name,omitempty"`
		SchoolName *string    `json:"school_name,omitempty"`
		BusID      *uuid.UUID `json:"bus_id,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	child, err := h.childService.UpdateChild(childID, parentID, req.Name, req.SchoolName, req.BusID)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, child, "Child updated successfully")
}

// DeleteChild deletes a child
func (h *ChildHandler) DeleteChild(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	parentID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	childID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid child ID", nil)
		return
	}

	if err := h.childService.DeleteChild(childID, parentID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Child deleted successfully")
}

