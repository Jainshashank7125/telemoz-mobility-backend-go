package maps

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/telemoz/backend/internal/config"
)

type Client struct {
	apiKey     string
	httpClient *http.Client
}

type GeocodeResponse struct {
	Results []struct {
		FormattedAddress string `json:"formatted_address"`
		Geometry         struct {
			Location struct {
				Lat float64 `json:"lat"`
				Lng float64 `json:"lng"`
			} `json:"location"`
		} `json:"geometry"`
	} `json:"results"`
	Status string `json:"status"`
}

type DistanceMatrixResponse struct {
	Rows []struct {
		Elements []struct {
			Distance struct {
				Value int    `json:"value"` // in meters
				Text  string `json:"text"`
			} `json:"distance"`
			Duration struct {
				Value int    `json:"value"` // in seconds
				Text  string `json:"text"`
			} `json:"duration"`
			Status string `json:"status"`
		} `json:"elements"`
	} `json:"rows"`
	Status string `json:"status"`
}

type RouteInfo struct {
	Distance    float64 // in km
	Duration    int     // in minutes
	DistanceText string
	DurationText string
}

func NewClient() *Client {
	return &Client{
		apiKey: config.AppConfig.Maps.APIKey,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

func (c *Client) Geocode(address string) (float64, float64, error) {
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("address", address)
	params.Add("key", c.apiKey)

	resp, err := c.httpClient.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return 0, 0, err
	}

	var geocodeResp GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		return 0, 0, err
	}

	if geocodeResp.Status != "OK" || len(geocodeResp.Results) == 0 {
		return 0, 0, fmt.Errorf("geocoding failed: %s", geocodeResp.Status)
	}

	location := geocodeResp.Results[0].Geometry.Location
	return location.Lat, location.Lng, nil
}

func (c *Client) ReverseGeocode(lat, lng float64) (string, error) {
	baseURL := "https://maps.googleapis.com/maps/api/geocode/json"
	params := url.Values{}
	params.Add("latlng", fmt.Sprintf("%f,%f", lat, lng))
	params.Add("key", c.apiKey)

	resp, err := c.httpClient.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var geocodeResp GeocodeResponse
	if err := json.Unmarshal(body, &geocodeResp); err != nil {
		return "", err
	}

	if geocodeResp.Status != "OK" || len(geocodeResp.Results) == 0 {
		return "", fmt.Errorf("reverse geocoding failed: %s", geocodeResp.Status)
	}

	return geocodeResp.Results[0].FormattedAddress, nil
}

func (c *Client) GetDistanceAndDuration(originLat, originLng, destLat, destLng float64) (*RouteInfo, error) {
	baseURL := "https://maps.googleapis.com/maps/api/distancematrix/json"
	params := url.Values{}
	params.Add("origins", fmt.Sprintf("%f,%f", originLat, originLng))
	params.Add("destinations", fmt.Sprintf("%f,%f", destLat, destLng))
	params.Add("key", c.apiKey)
	params.Add("units", "metric")

	resp, err := c.httpClient.Get(fmt.Sprintf("%s?%s", baseURL, params.Encode()))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var matrixResp DistanceMatrixResponse
	if err := json.Unmarshal(body, &matrixResp); err != nil {
		return nil, err
	}

	if matrixResp.Status != "OK" || len(matrixResp.Rows) == 0 || len(matrixResp.Rows[0].Elements) == 0 {
		return nil, fmt.Errorf("distance matrix failed: %s", matrixResp.Status)
	}

	element := matrixResp.Rows[0].Elements[0]
	if element.Status != "OK" {
		return nil, fmt.Errorf("route calculation failed: %s", element.Status)
	}

	return &RouteInfo{
		Distance:     float64(element.Distance.Value) / 1000.0, // convert to km
		Duration:     element.Duration.Value / 60,              // convert to minutes
		DistanceText: element.Distance.Text,
		DurationText: element.Duration.Text,
	}, nil
}

