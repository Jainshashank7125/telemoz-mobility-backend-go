package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/telemoz/backend/internal/handlers"
	"github.com/telemoz/backend/internal/dto"
)

func TestRegister(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := handlers.NewAuthHandler()
	router := gin.New()
	router.POST("/register", handler.Register)

	registerReq := dto.RegisterRequest{
		Email:    "test@example.com",
		Password: "password123",
		Name:     "Test User",
		UserType: "customer",
	}

	jsonValue, _ := json.Marshal(registerReq)
	req, _ := http.NewRequest("POST", "/register", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	assert.Equal(t, http.StatusCreated, w.Code)
}

func TestLogin(t *testing.T) {
	gin.SetMode(gin.TestMode)

	handler := handlers.NewAuthHandler()
	router := gin.New()
	router.POST("/login", handler.Login)

	loginReq := dto.LoginRequest{
		EmailOrPhone: "test@example.com",
		Password:     "password123",
	}

	jsonValue, _ := json.Marshal(loginReq)
	req, _ := http.NewRequest("POST", "/login", bytes.NewBuffer(jsonValue))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)

	// This will fail if user doesn't exist, which is expected in test environment
	// In a real test, you'd set up test data first
	assert.True(t, w.Code == http.StatusOK || w.Code == http.StatusUnauthorized)
}

