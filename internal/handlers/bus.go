package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type BusHandler struct {
	busService services.BusService
}

func NewBusHandler() *BusHandler {
	return &BusHandler{
		busService: services.NewBusService(),
	}
}

// GetBusByChildID gets the bus for a specific child
func (h *BusHandler) GetBusByChildID(c *gin.Context) {
	childID, err := uuid.Parse(c.Param("childId"))
	if err != nil {
		utils.BadRequest(c, "Invalid child ID", nil)
		return
	}

	bus, err := h.busService.GetBusByChildID(childID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, bus, "Bus retrieved successfully")
}

// TrackBus gets bus location and tracking data
func (h *BusHandler) TrackBus(c *gin.Context) {
	busID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.BadRequest(c, "Invalid bus ID", nil)
		return
	}

	location, err := h.busService.GetBusLocation(busID)
	if err != nil {
		utils.NotFound(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, location, "Bus location retrieved successfully")
}

