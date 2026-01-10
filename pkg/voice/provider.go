package voice

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
	MakeCall(to, message string) error
}

type TwilioVoiceProvider struct {
	accountSID string
	authToken  string
	httpClient *http.Client
}

func NewProvider() Provider {
	cfg := config.AppConfig.Voice
	
	if cfg.Provider == "twilio" {
		return &TwilioVoiceProvider{
			accountSID: cfg.APIKey,
			authToken:  cfg.APISecret,
			httpClient: &http.Client{
				Timeout: 30 * time.Second,
			},
		}
	}

	// Default: no-op provider for development
	return &NoOpProvider{}
}

func (p *TwilioVoiceProvider) MakeCall(to, message string) error {
	apiURL := fmt.Sprintf("https://api.twilio.com/2010-04-01/Accounts/%s/Calls.json", p.accountSID)

	// For Twilio, you typically need a TwiML URL or a callback URL
	// This is a simplified implementation
	data := url.Values{}
	data.Set("To", to)
	data.Set("From", "+1234567890") // This should be your Twilio phone number
	data.Set("Url", "http://your-server.com/twiml") // TwiML URL that contains the message

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
		return fmt.Errorf("failed to make call: status %d, body: %s", resp.StatusCode, string(body))
	}

	return nil
}

type NoOpProvider struct{}

func (p *NoOpProvider) MakeCall(to, message string) error {
	// No-op implementation for development/testing
	return nil
}

