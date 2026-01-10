package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/repositories"
	"github.com/telemoz/backend/internal/utils"
)

type ProfileHandler struct {
	userRepo repositories.UserRepository
}

func NewProfileHandler() *ProfileHandler {
	return &ProfileHandler{
		userRepo: repositories.NewUserRepository(),
	}
}

// GetProfile gets the current user's profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	response := dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		UserType:  string(user.UserType),
		AvatarURL: user.AvatarURL,
	}
	if user.Phone != nil {
		response.Phone = user.Phone
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Profile retrieved successfully")
}

// UpdateProfile updates the current user's profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
	userIDStr, _ := c.Get("user_id")
	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	user, err := h.userRepo.FindByID(userID)
	if err != nil {
		utils.NotFound(c, "User not found")
		return
	}

	var req struct {
		Name     *string `json:"name,omitempty"`
		Phone    *string `json:"phone,omitempty"`
		AvatarURL *string `json:"avatar_url,omitempty"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	if req.Name != nil {
		user.Name = utils.SanitizeString(*req.Name)
	}
	if req.Phone != nil {
		phone := utils.SanitizeString(*req.Phone)
		if !utils.ValidatePhone(phone) {
			utils.BadRequest(c, "Invalid phone format", nil)
			return
		}
		user.Phone = &phone
	}
	if req.AvatarURL != nil {
		user.AvatarURL = req.AvatarURL
	}

	if err := h.userRepo.Update(user); err != nil {
		utils.InternalError(c, "Failed to update profile")
		return
	}

	response := dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		UserType:  string(user.UserType),
		AvatarURL: user.AvatarURL,
	}
	if user.Phone != nil {
		response.Phone = user.Phone
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Profile updated successfully")
}

