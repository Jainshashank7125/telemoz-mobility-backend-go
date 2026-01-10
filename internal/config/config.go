package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Traccar  TraccarConfig
	Maps     MapsConfig
	SMS      SMSConfig
	Voice    VoiceConfig
	Firebase FirebaseConfig
	CORS     CORSConfig
}

type ServerConfig struct {
	Port    string
	Env     string
	GinMode string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type JWTConfig struct {
	Secret        string
	AccessExpiry  time.Duration
	RefreshExpiry time.Duration
}

type TraccarConfig struct {
	URL      string
	Username string
	Password string
}

type MapsConfig struct {
	APIKey string
}

type SMSConfig struct {
	Provider  string
	APIKey    string
	APISecret string
	FromNumber string
}

type VoiceConfig struct {
	Provider  string
	APIKey    string
	APISecret string
}

type FirebaseConfig struct {
	ProjectID      string
	CredentialsPath string
}

type CORSConfig struct {
	AllowedOrigins []string
}

var AppConfig *Config

func Load() error {
	// Load .env file if it exists (ignore error if file doesn't exist)
	_ = godotenv.Load()

	accessExpiry, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "15m"))
	refreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))

	AppConfig = &Config{
		Server: ServerConfig{
			Port:    getEnv("PORT", "8080"),
			Env:     getEnv("ENV", "development"),
			GinMode: getEnv("GIN_MODE", "debug"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "telemoz"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "telemoz_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			Secret:        getEnv("JWT_SECRET", "change-me-in-production"),
			AccessExpiry:  accessExpiry,
			RefreshExpiry: refreshExpiry,
		},
		Traccar: TraccarConfig{
			URL:      getEnv("TRACCAR_URL", "http://localhost:8082"),
			Username: getEnv("TRACCAR_USERNAME", "admin"),
			Password: getEnv("TRACCAR_PASSWORD", "admin"),
		},
		Maps: MapsConfig{
			APIKey: getEnv("GOOGLE_MAPS_API_KEY", ""),
		},
		SMS: SMSConfig{
			Provider:   getEnv("SMS_PROVIDER", "twilio"),
			APIKey:     getEnv("SMS_API_KEY", ""),
			APISecret:   getEnv("SMS_API_SECRET", ""),
			FromNumber: getEnv("SMS_FROM_NUMBER", ""),
		},
		Voice: VoiceConfig{
			Provider:  getEnv("VOICE_PROVIDER", "twilio"),
			APIKey:    getEnv("VOICE_API_KEY", ""),
			APISecret: getEnv("VOICE_API_SECRET", ""),
		},
		Firebase: FirebaseConfig{
			ProjectID:       getEnv("FIREBASE_PROJECT_ID", ""),
			CredentialsPath: getEnv("FIREBASE_CREDENTIALS_PATH", "./firebase-credentials.json"),
		},
		CORS: CORSConfig{
			AllowedOrigins: parseStringSlice(getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:19006")),
		},
	}

	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func parseStringSlice(value string) []string {
	if value == "" {
		return []string{}
	}
	result := []string{}
	for _, item := range splitString(value, ",") {
		if trimmed := trimString(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

func splitString(s, sep string) []string {
	result := []string{}
	current := ""
	for _, char := range s {
		if string(char) == sep {
			if current != "" {
				result = append(result, current)
				current = ""
			}
		} else {
			current += string(char)
		}
	}
	if current != "" {
		result = append(result, current)
	}
	return result
}

func trimString(s string) string {
	start := 0
	end := len(s)
	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n') {
		start++
	}
	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n') {
		end--
	}
	return s[start:end]
}

