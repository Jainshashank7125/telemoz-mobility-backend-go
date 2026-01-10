package utils

import (
	"regexp"
	"strings"
)

func ValidateEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func ValidatePhone(phone string) bool {
	// Remove spaces and special characters
	cleaned := regexp.MustCompile(`[^\d+]`).ReplaceAllString(phone, "")
	// Check if it's a valid phone number (at least 10 digits)
	return len(cleaned) >= 10
}

func ValidatePassword(password string) (bool, string) {
	if len(password) < 8 {
		return false, "Password must be at least 8 characters long"
	}
	if len(password) > 128 {
		return false, "Password must be less than 128 characters"
	}
	// Check for at least one letter and one number
	hasLetter := regexp.MustCompile(`[a-zA-Z]`).MatchString(password)
	hasNumber := regexp.MustCompile(`[0-9]`).MatchString(password)
	
	if !hasLetter {
		return false, "Password must contain at least one letter"
	}
	if !hasNumber {
		return false, "Password must contain at least one number"
	}
	return true, ""
}

func ValidateCoordinates(lat, lng float64) bool {
	return lat >= -90 && lat <= 90 && lng >= -180 && lng <= 180
}

func SanitizeString(s string) string {
	s = strings.TrimSpace(s)
	return s
}

