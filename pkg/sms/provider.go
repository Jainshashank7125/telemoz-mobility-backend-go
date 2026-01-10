package sms

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/telemoz/backend/internal/config"
)

type Provider interface {
	SendSMS(to, message string) error
}

type TwilioProvider struct {
	accountSID string
	authToken  string
	fromNumber string
	httpClient *http.Client
}

func NewProvider() Provider {
	cfg := config.AppConfig.SMS
	
	if cfg.Provider == "twilio" {
		return &TwilioProvider{
			accountSID: cfg.APIKey,
			authToken:  cfg.APISecret,
			fromNumber: cfg.FromNumber,
			httpClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	}

	// Default: no-op provider for development
	return &NoOpProvider{}
}

func (p *TwilioProvider) SendSMS(to, message string) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Messages.json", p.accountSID)

	data := url.Values{}
	data.Set("From", p.fromNumber)
	data.Set("To", to)
	data.Set("Body", message)

	req, err := http.NewRequest("POST", apiURL, strings.NewReader(data.Encode()))
	if err != nil {
		return err
	}

	req.SetBasicAuth(p.accountSID, p.authToken)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := p.httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated && resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("failed to send SMS: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

type NoOpProvider struct{}

func (p *NoOpProvider) SendSMS(to, message string) error {
	// No-op implementation for development/testing
	return nil
}

