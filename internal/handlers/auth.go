package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/services"
	"github.com/telemoz/backend/internal/utils"
)

type AuthHandler struct {
	authService services.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{
		authService: services.NewAuthService(),
	}
}

// Register handles user registration
// @Summary Register a new user
// @Description Register a new customer, driver, or parent user
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RegisterRequest true "Registration data"
// @Success 201 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Router /api/auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req dto.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	response, err := h.authService.Register(req)
	if err != nil {
		utils.BadRequest(c, err.Error(), nil)
		return
	}

	utils.SuccessResponse(c, http.StatusCreated, response, "Registration successful")
}

// Login handles user login
// @Summary Login user
// @Description Authenticate user and return JWT tokens
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req dto.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	response, err := h.authService.Login(req)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Login successful")
}

// RefreshToken handles token refresh
// @Summary Refresh access token
// @Description Get a new access token using refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Param request body dto.RefreshTokenRequest true "Refresh token"
// @Success 200 {object} dto.LoginResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req dto.RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.BadRequest(c, "Invalid request data", err.Error())
		return
	}

	response, err := h.authService.RefreshToken(req.RefreshToken)
	if err != nil {
		utils.Unauthorized(c, err.Error())
		return
	}

	utils.SuccessResponse(c, http.StatusOK, response, "Token refreshed successfully")
}

// Logout handles user logout
// @Summary Logout user
// @Description Invalidate refresh token
// @Tags auth
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.RefreshTokenRequest true "Refresh token to invalidate"
// @Success 200 {object} dto.SuccessResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /api/auth/logout [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	userIDStr, exists := c.Get("user_id")
	if !exists {
		utils.Unauthorized(c, "User ID not found")
		return
	}

	userID, err := uuid.Parse(userIDStr.(string))
	if err != nil {
		utils.Unauthorized(c, "Invalid user ID")
		return
	}

	var req dto.RefreshTokenRequest
	refreshToken := ""
	if err := c.ShouldBindJSON(&req); err == nil {
		refreshToken = req.RefreshToken
	}

	if err := h.authService.Logout(userID, refreshToken); err != nil {
		utils.InternalError(c, "Failed to logout")
		return
	}

	utils.SuccessResponse(c, http.StatusOK, nil, "Logout successful")
}

