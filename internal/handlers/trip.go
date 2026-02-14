package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type TripHandler struct {
	tripService services.TripService
}

func NewTripHandler() *TripHandler {
	return &TripHandler{
		tripService: services.NewTripService(),
	}
}

// EstimateFare estimates the fare for a trip
func (h *TripHandler) EstimateFare(c *gin.Context) {
	var req dto.EstimateFareRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", nil)
		return
	}

	response, err := h.tripService.EstimateFare(req)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Fare estimated successfully")
}

// CreateTrip creates a new trip
func (h *TripHandler) CreateTrip(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	customerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	var req dto.CreateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", nil)
		return
	}

	trip, err := h.tripService.CreateTrip(customerID, req)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, trip, "Trip created successfully")
}

// GetActiveTrip gets the active trip for a customer
func (h *TripHandler) GetActiveTrip(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	customerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	trip, err := h.tripService.GetActiveTrip(customerID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, trip, "Active trip retrieved successfully")
}

// GetTripHistory gets the trip history for a customer
func (h *TripHandler) GetTripHistory(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	customerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	limit := 20
	offset := 0

	trips, err := h.tripService.GetTripHistory(customerID, limit, offset)
	if err != nil {
		utils.InternalError(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, trips, "Trip history retrieved successfully")
}

// GetTripByID gets a trip by ID
func (h *TripHandler) GetTripByID(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid trip ID", nil)
		return
	}

	trip, err := h.tripService.GetTripByID(tripID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, trip, "Trip retrieved successfully")
}

// UpdateTrip updates a trip
func (h *TripHandler) UpdateTrip(c *gin.Context) {
	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid trip ID", nil)
		return
	}

	var req dto.UpdateTripRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", nil)
		return
	}

	trip, err := h.tripService.UpdateTrip(tripID, req)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, trip, "Trip updated successfully")
}

// CancelTrip cancels a trip
func (h *TripHandler) CancelTrip(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	customerID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	tripID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid trip ID", nil)
		return
	}

	if err := h.tripService.CancelTrip(tripID, customerID); err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Trip cancelled successfully")
}
