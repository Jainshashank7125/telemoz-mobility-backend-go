package traccar

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/telemoz/backend/internal/config"
)

type Client struct {
	baseURL    string
	username   string
	password   string
	httpClient *http.Client
}

type Device struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	UniqueID string `json:"uniqueId"`
	Status   string `json:"status"`
}

type Position struct {
	ID        int     `json:"id"`
	DeviceID  int     `json:"deviceId"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Speed     float64 `json:"speed"`
	Course    float64 `json:"course"`
	Accuracy  float64 `json:"accuracy"`
	FixTime   string  `json:"fixTime"`
}

func NewClient() *Client {
	cfg := config.AppConfig.Traccar
	return &Client{
		baseURL:  cfg.URL,
		username: cfg.Username,
		password: cfg.Password,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) createRequest(method, endpoint string, body interface{}) (*http.Request, error) {
	url := fmt.Sprintf("%s/api/%s", c.baseURL, endpoint)
	
	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, err
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.username, c.password)
	req.Header.Set("Content-Type", "application/json")

	return req, nil
}

func (c *Client) CreateDevice(name, uniqueID string) (*Device, error) {
	device := map[string]interface{}{
		"name":     name,
		"uniqueId": uniqueID,
	}

	req, err := c.createRequest("POST", "devices", device)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("failed to create device: %s", string(body))
	}

	var createdDevice Device
	if err := json.NewDecoder(resp.Body).Decode(&createdDevice); err != nil {
		return nil, err
	}

	return &createdDevice, nil
}

func (c *Client) GetDevice(deviceID int) (*Device, error) {
	req, err := c.createRequest("GET", fmt.Sprintf("devices/%d", deviceID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get device: status %d", resp.StatusCode)
	}

	var device Device
	if err := json.NewDecoder(resp.Body).Decode(&device); err != nil {
		return nil, err
	}

	return &device, nil
}

func (c *Client) GetLatestPosition(deviceID int) (*Position, error) {
	req, err := c.createRequest("GET", fmt.Sprintf("positions?deviceId=%d", deviceID), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to get position: status %d", resp.StatusCode)
	}

	var positions []Position
	if err := json.NewDecoder(resp.Body).Decode(&positions); err != nil {
		return nil, err
	}

	if len(positions) == 0 {
		return nil, fmt.Errorf("no positions found for device %d", deviceID)
	}

	return &positions[0], nil
}

