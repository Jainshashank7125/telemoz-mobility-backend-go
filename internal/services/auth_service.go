package services

import (
	"errors"

	"github.com/google/uuid"
	"github.com/telemoz/backend/internal/dto"
	"github.com/telemoz/backend/internal/models"
	"github.com/telemoz/backend/internal/repositories"
	"github.com/telemoz/backend/internal/utils"
)

type AuthService interface {
	Register(req dto.RegisterRequest) (*dto.LoginResponse, error)
	Login(req dto.LoginRequest) (*dto.LoginResponse, error)
	RefreshToken(refreshToken string) (*dto.LoginResponse, error)
	Logout(userID uuid.UUID, refreshToken string) error
}

type authService struct {
	userRepo         repositories.UserRepository
	refreshTokenRepo repositories.RefreshTokenRepository
}

func NewAuthService() AuthService {
	return &authService{
		userRepo:         repositories.NewUserRepository(),
		refreshTokenRepo: repositories.NewRefreshTokenRepository(),
	}
}

func (s *authService) Register(req dto.RegisterRequest) (*dto.LoginResponse, error) {
	// Validate email
	if !utils.ValidateEmail(req.Email) {
		return nil, errors.New("invalid email format")
	}

	// Validate phone if provided
	if req.Phone != "" && !utils.ValidatePhone(req.Phone) {
		return nil, errors.New("invalid phone format")
	}

	// Validate password
	if valid, msg := utils.ValidatePassword(req.Password); !valid {
		return nil, errors.New(msg)
	}

	// Check if user already exists
	_, err := s.userRepo.FindByEmail(req.Email)
	if err == nil {
		return nil, errors.New("user with this email already exists")
	}

	// Check phone if provided
	if req.Phone != "" {
		_, err = s.userRepo.FindByPhone(req.Phone)
		if err == nil {
			return nil, errors.New("user with this phone already exists")
		}
	}

	// Hash password
	passwordHash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, errors.New("failed to hash password")
	}

	// Create user
	user := &models.User{
		Email:        req.Email,
		PasswordHash: passwordHash,
		Name:         utils.SanitizeString(req.Name),
		UserType:     models.UserType(req.UserType),
		IsActive:     true,
	}

	if req.Phone != "" {
		phone := utils.SanitizeString(req.Phone)
		user.Phone = &phone
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, errors.New("failed to create user")
	}

	// Generate tokens
	return s.generateTokens(user)
}

func (s *authService) Login(req dto.LoginRequest) (*dto.LoginResponse, error) {
	// Find user by email or phone
	user, err := s.userRepo.FindByEmailOrPhone(req.EmailOrPhone)
	if err != nil {
		return nil, errors.New("invalid credentials")
	}

	// Check password
	if !utils.CheckPassword(req.Password, user.PasswordHash) {
		return nil, errors.New("invalid credentials")
	}

	// Check if user is active
	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate tokens
	return s.generateTokens(user)
}

func (s *authService) RefreshToken(refreshToken string) (*dto.LoginResponse, error) {
	// Find refresh token in database
	token, err := s.refreshTokenRepo.FindByToken(refreshToken)
	if err != nil {
		return nil, errors.New("invalid refresh token")
	}

	// Get user
	user, err := s.userRepo.FindByID(token.UserID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if !user.IsActive {
		return nil, errors.New("account is deactivated")
	}

	// Generate new tokens
	return s.generateTokens(user)
}

func (s *authService) Logout(userID uuid.UUID, refreshToken string) error {
	if refreshToken != "" {
		return s.refreshTokenRepo.DeleteByToken(refreshToken)
	}
	// If no refresh token provided, delete all tokens for user
	return s.refreshTokenRepo.DeleteByUserID(userID)
}

func (s *authService) generateTokens(user *models.User) (*dto.LoginResponse, error) {
	// Generate access token
	accessToken, err := utils.GenerateAccessToken(user.ID, string(user.UserType), user.Email)
	if err != nil {
		return nil, errors.New("failed to generate access token")
	}

	// Generate refresh token
	refreshToken, expiresAt, err := utils.GenerateRefreshToken(user.ID)
	if err != nil {
		return nil, errors.New("failed to generate refresh token")
	}

	// Save refresh token to database
	refreshTokenModel := &models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: expiresAt,
	}
	if err := s.refreshTokenRepo.Create(refreshTokenModel); err != nil {
		return nil, errors.New("failed to save refresh token")
	}

	// Build user response
	userResponse := dto.UserResponse{
		ID:        user.ID.String(),
		Email:     user.Email,
		Name:      user.Name,
		UserType:  string(user.UserType),
		AvatarURL: user.AvatarURL,
	}
	if user.Phone != nil {
		userResponse.Phone = user.Phone
	}

	return &dto.LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		User:         userResponse,
	}, nil
}

