package traccar

import (
	"fmt"
	"log"
	"time"

	"github.com/gorilla/websocket"
	"github.com/telemoz/backend/internal/config"
)

type WebSocketClient struct {
	conn     *websocket.Conn
	url      string
	username string
	password string
	handlers map[int]func(Position)
}

type WebSocketMessage struct {
	Devices []Device  `json:"devices,omitempty"`
	Events  []Event   `json:"events,omitempty"`
	Positions []Position `json:"positions,omitempty"`
}

type Event struct {
	Type     string `json:"type"`
	DeviceID int    `json:"deviceId"`
	Position Position `json:"position,omitempty"`
}

func NewWebSocketClient() *WebSocketClient {
	cfg := config.AppConfig.Traccar
	url := fmt.Sprintf("ws://%s/api/socket", cfg.URL)
	
	// Remove http:// or https:// prefix
	if len(url) > 7 && url[:7] == "http://" {
		url = url[7:]
	} else if len(url) > 8 && url[:8] == "https://" {
		url = url[8:]
		url = "wss://" + url
	} else {
		url = "ws://" + url
	}

	return &WebSocketClient{
		url:      url,
		username: cfg.Username,
		password: cfg.Password,
		handlers: make(map[int]func(Position)),
	}
}

func (c *WebSocketClient) Connect() error {
	dialer := websocket.Dialer{
		HandshakeTimeout: 10 * time.Second,
	}

	conn, _, err := dialer.Dial(c.url, nil)
	if err != nil {
		return fmt.Errorf("failed to connect to Traccar WebSocket: %w", err)
	}

	c.conn = conn

	// Send authentication
	authMsg := map[string]string{
		"type":     "authenticate",
		"username": c.username,
		"password": c.password,
	}

	if err := conn.WriteJSON(authMsg); err != nil {
		conn.Close()
		return fmt.Errorf("failed to authenticate: %w", err)
	}

	// Start reading messages
	go c.readMessages()

	return nil
}

func (c *WebSocketClient) readMessages() {
	for {
		var msg WebSocketMessage
		if err := c.conn.ReadJSON(&msg); err != nil {
			log.Printf("WebSocket read error: %v", err)
			break
		}

		// Handle positions
		for _, pos := range msg.Positions {
			if handler, exists := c.handlers[pos.DeviceID]; exists {
				handler(pos)
			}
		}

		// Handle events
		for _, event := range msg.Events {
			if event.Type == "position" && event.Position.DeviceID > 0 {
				if handler, exists := c.handlers[event.DeviceID]; exists {
					handler(event.Position)
				}
			}
		}
	}
}

func (c *WebSocketClient) SubscribeToDevice(deviceID int, handler func(Position)) {
	c.handlers[deviceID] = handler
}

func (c *WebSocketClient) UnsubscribeFromDevice(deviceID int) {
	delete(c.handlers, deviceID)
}

func (c *WebSocketClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

