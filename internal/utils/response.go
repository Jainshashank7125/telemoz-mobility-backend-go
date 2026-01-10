package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Response struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
	Code    string      `json:"code,omitempty"`
	Details interface{} `json:"details,omitempty"`
}

func SuccessResponse(c *gin.Context, statusCode int, data interface{}, message string) {
	c.JSON(statusCode, Response{
		Success: true,
		Data:    data,
		Message: message,
	})
}

func ErrorResponse(c *gin.Context, statusCode int, errorMsg, code string, details interface{}) {
	c.JSON(statusCode, Response{
		Success: false,
		Error:   errorMsg,
		Code:    code,
		Details: details,
	})
}

func BadRequest(c *gin.Context, errorMsg string, details interface{}) {
	ErrorResponse(c, http.StatusBadRequest, errorMsg, "VALIDATION_ERROR", details)
}

func Unauthorized(c *gin.Context, errorMsg string) {
	ErrorResponse(c, http.StatusUnauthorized, errorMsg, "AUTH_REQUIRED", nil)
}

func Forbidden(c *gin.Context, errorMsg string) {
	ErrorResponse(c, http.StatusForbidden, errorMsg, "FORBIDDEN", nil)
}

func NotFound(c *gin.Context, errorMsg string) {
	ErrorResponse(c, http.StatusNotFound, errorMsg, "NOT_FOUND", nil)
}

func InternalError(c *gin.Context, errorMsg string) {
	ErrorResponse(c, http.StatusInternalServerError, errorMsg, "INTERNAL_ERROR", nil)
}

