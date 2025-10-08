package provider

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// Global semaphore to ensure only one WebSocket operation at a time
var globalWSMutex sync.Mutex

// Client represents the Uptime Kuma API client
type Client struct {
	BaseURL           string
	Username          string
	Password          string
	HTTPClient        *http.Client
	wsConn            *websocket.Conn
	connected         bool
	mu                sync.RWMutex
	wsMu              sync.Mutex // Protects WebSocket writes from concurrent access
	eventID           int
	responses         map[int]chan SocketResponse
	respMu            sync.RWMutex
	userID            int
	token             string
	monitors          map[string]interface{} // Cache for monitors from monitorList event
	monitorsMu        sync.RWMutex
	notificationCache []Notification // Cache for notifications from notificationList event
	notificationsMu   sync.RWMutex
}

// SocketIOMessage represents a Socket.IO message
type SocketIOMessage struct {
	Type      int           `json:"0,omitempty"`
	Namespace string        `json:"1,omitempty"`
	Event     string        `json:"2,omitempty"`
	ID        int           `json:"3,omitempty"`
	Data      []interface{} `json:"4,omitempty"`
}

// SocketResponse represents a Socket.IO response
type SocketResponse struct {
	Event string        `json:"event"`
	Data  []interface{} `json:"data"`
	Error string        `json:"error,omitempty"`
}

// Monitor represents an Uptime Kuma monitor
type Monitor struct {
	ID                  int               `json:"id,omitempty"`
	Name                string            `json:"name"`
	Type                string            `json:"type"`
	URL                 string            `json:"url,omitempty"`
	Hostname            string            `json:"hostname,omitempty"`
	Port                int               `json:"port,omitempty"`
	Interval            int               `json:"interval"`
	Timeout             int               `json:"timeout"`
	RetryInterval       int               `json:"retryInterval,omitempty"`
	ResendInterval      int               `json:"resendInterval,omitempty"`
	MaxRetries          int               `json:"maxretries,omitempty"`
	UpsideDown          bool              `json:"upsideDown,omitempty"`
	MaxRedirects        int               `json:"maxredirects,omitempty"`
	AcceptedStatusCodes []string          `json:"accepted_statuscodes,omitempty"`
	FollowRedirect      bool              `json:"follow_redirect,omitempty"`
	Tags                []string          `json:"tags,omitempty"`
	NotificationIDList  []int             `json:"notificationIDList,omitempty"`
	Active              bool              `json:"active"`
	IgnoreTLS           bool              `json:"ignoreTls,omitempty"`
	HTTPMethod          string            `json:"method,omitempty"`
	Body                string            `json:"body,omitempty"`
	Headers             map[string]string `json:"headers,omitempty"`
	BasicAuthUser       string            `json:"basic_auth_user,omitempty"`
	BasicAuthPass       string            `json:"basic_auth_pass,omitempty"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token"`
	OK    bool   `json:"ok"`
	Msg   string `json:"msg"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// MonitorResponse represents the response when getting a monitor
type MonitorResponse struct {
	Monitor Monitor `json:"monitor"`
	OK      bool    `json:"ok"`
	Msg     string  `json:"msg"`
}

// MonitorsResponse represents the response when getting all monitors
type MonitorsResponse struct {
	Monitors []Monitor `json:"monitors"`
	OK       bool      `json:"ok"`
	Msg      string    `json:"msg"`
}

// Notification represents an Uptime Kuma notification
type Notification struct {
	ID            int                    `json:"id,omitempty"`
	Name          string                 `json:"name"`
	Type          string                 `json:"type"`
	IsDefault     bool                   `json:"isDefault,omitempty"`
	ApplyExisting bool                   `json:"applyExisting,omitempty"`
	Active        bool                   `json:"active,omitempty"`
	UserID        int                    `json:"userId,omitempty"`
	Config        map[string]interface{} `json:"config,omitempty"`
}

// NotificationResponse represents the response when getting a notification
type NotificationResponse struct {
	Notification Notification `json:"notification"`
	OK           bool         `json:"ok"`
	Msg          string       `json:"msg"`
}

// NotificationsResponse represents the response when getting all notifications
type NotificationsResponse struct {
	Notifications []Notification `json:"notifications"`
	OK            bool           `json:"ok"`
	Msg           string         `json:"msg"`
}

// NotificationCreateResponse represents the response when creating a notification
type NotificationCreateResponse struct {
	ID  int    `json:"id"`
	OK  bool   `json:"ok"`
	Msg string `json:"msg"`
}

// NewClient creates a new Uptime Kuma API client
func NewClient(baseURL, username, password string) (*Client, error) {
	client := &Client{
		BaseURL:    baseURL,
		Username:   username,
		Password:   password,
		HTTPClient: &http.Client{Timeout: 30 * time.Second},
		responses:  make(map[int]chan SocketResponse),
		monitors:   make(map[string]interface{}),
	}

	// Connect to Socket.IO endpoint
	err := client.connect()
	if err != nil {
		return nil, fmt.Errorf("failed to connect: %w", err)
	}

	// Authenticate
	err = client.login()
	if err != nil {
		client.disconnect()
		return nil, fmt.Errorf("authentication failed: %w", err)
	}

	// Wait a moment for the initial events to be received
	time.Sleep(1 * time.Second)

	return client, nil
}

// connect establishes WebSocket connection to Uptime Kuma Socket.IO endpoint
func (c *Client) connect() error {
	// Parse base URL
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid URL: %w", err)
	}

	// Convert to WebSocket URL
	scheme := "ws"
	if u.Scheme == "https" {
		scheme = "wss"
	}

	// Socket.IO WebSocket endpoint
	wsURL := fmt.Sprintf("%s://%s/socket.io/?EIO=4&transport=websocket", scheme, u.Host)

	// Connect to WebSocket
	c.wsConn, _, err = websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		return fmt.Errorf("websocket connection failed: %w", err)
	}

	c.mu.Lock()
	c.connected = true
	c.mu.Unlock()

	// Start message handler
	go c.handleMessages()

	// Send Socket.IO connect message
	globalWSMutex.Lock()
	err = c.wsConn.WriteMessage(websocket.TextMessage, []byte("40"))
	globalWSMutex.Unlock()
	if err != nil {
		return fmt.Errorf("failed to send connect message: %w", err)
	}

	// Wait a moment for connection to be established
	time.Sleep(100 * time.Millisecond)

	return nil
}

// disconnect closes the WebSocket connection
func (c *Client) disconnect() {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.wsConn != nil && c.connected {
		c.wsConn.Close()
		c.connected = false
	}
}

// handleMessages processes incoming WebSocket messages
func (c *Client) handleMessages() {
	for {
		c.mu.RLock()
		if !c.connected {
			c.mu.RUnlock()
			break
		}
		conn := c.wsConn
		c.mu.RUnlock()

		_, message, err := conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse Socket.IO message format
		c.parseMessage(string(message))
	}
}

// parseMessage parses Socket.IO protocol messages
func (c *Client) parseMessage(message string) {
	if len(message) < 2 {
		return
	}

	// Socket.IO message format: [type][optional-id][optional-json]
	msgType := message[:1]
	content := message[1:]

	switch msgType {
	case "4": // Message type
		if strings.HasPrefix(content, "3") {
			// Callback response: 43[ack_id][response_data]
			content = content[1:] // Remove "3"

			// Extract callback ID and response data
			var callbackID int
			var responseData string

			// Find where the JSON starts (after the callback ID)
			for i := 0; i < len(content); i++ {
				if content[i] == '[' || content[i] == '{' {
					if i > 0 {
						callbackID, _ = strconv.Atoi(content[:i])
					}
					responseData = content[i:]
					break
				}
			}

			// Parse response data - it should be a single JSON object, not an array
			var responseObj interface{}
			if err := json.Unmarshal([]byte(responseData), &responseObj); err == nil {
				response := SocketResponse{
					Data: []interface{}{responseObj}, // Wrap in array for consistent handling
				}

				// Send to waiting goroutine
				c.respMu.RLock()
				if ch, exists := c.responses[callbackID]; exists {
					select {
					case ch <- response:
					default:
					}
				}
				c.respMu.RUnlock()
			}
		} else if strings.HasPrefix(content, "2") {
			// Event message: 42["event", data...]
			content = content[1:] // Remove "2"
			var eventData []interface{}
			if err := json.Unmarshal([]byte(content), &eventData); err == nil {
				if len(eventData) >= 2 {
					event := eventData[0].(string)
					data := eventData[1:]

					// Handle specific events we care about
					if event == "monitorList" && len(data) > 0 {
						// Cache the monitor list data
						if monitorData, ok := data[0].(map[string]interface{}); ok {
							c.monitorsMu.Lock()
							c.monitors = monitorData
							c.monitorsMu.Unlock()
						}
					} else if event == "notificationList" && len(data) > 0 {
						// Cache the notification list data
						if notifList, ok := data[0].([]interface{}); ok {
							c.notificationsMu.Lock()
							c.notificationCache = c.notificationCache[:0] // Clear existing cache
							for _, item := range notifList {
								if notifMap, ok := item.(map[string]interface{}); ok {
									notification := parseNotificationMap(notifMap)
									c.notificationCache = append(c.notificationCache, notification)
								}
							}
							c.notificationsMu.Unlock()
						}
					}
				}
			}
		}
	case "0": // Connect
		// Connection successful
	case "3": // Heartbeat
		// Send pong
		globalWSMutex.Lock()
		c.wsConn.WriteMessage(websocket.TextMessage, []byte("3"))
		globalWSMutex.Unlock()
	}
}

// login authenticates with Uptime Kuma using Socket.IO
func (c *Client) login() error {
	// Send login event via Socket.IO
	loginData := map[string]interface{}{
		"username": c.Username,
		"password": c.Password,
		"token":    "",
	}

	response, err := c.call("login", loginData)
	if err != nil {
		return fmt.Errorf("login failed: %w", err)
	}

	// Check if login was successful
	if ok, exists := response["ok"]; exists {
		if okBool, isBool := ok.(bool); isBool && !okBool {
			msg := "login failed"
			if msgStr, exists := response["msg"]; exists {
				if msgString, isString := msgStr.(string); isString {
					msg = msgString
				}
			}
			return fmt.Errorf("authentication failed: %s", msg)
		}
	}

	// Store authentication information
	if token, exists := response["token"]; exists {
		if tokenStr, isString := token.(string); isString {
			c.token = tokenStr
		}
	}

	// Store user ID - in Uptime Kuma, user ID is typically 1 for admin
	// but we should check if it's provided in the response
	if userID, exists := response["userID"]; exists {
		if id, ok := userID.(float64); ok {
			c.userID = int(id)
		}
	} else {
		// Default to user ID 1 for admin if not provided
		c.userID = 1
	}

	return nil
}

// emit sends a Socket.IO event without waiting for response (fire-and-forget)
func (c *Client) emit(event string, data interface{}) error {
	c.mu.Lock()
	if !c.connected || c.wsConn == nil {
		c.mu.Unlock()
		return fmt.Errorf("not connected")
	}

	conn := c.wsConn
	c.mu.Unlock()

	// Create Socket.IO event message without acknowledgment: 42["event", data]
	eventData := []interface{}{event, data}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Socket.IO message format without acknowledgment: 42[json_array]
	message := fmt.Sprintf("42%s", string(eventJSON))

	// Use global mutex only for the WebSocket write operation
	globalWSMutex.Lock()
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	globalWSMutex.Unlock()

	if err != nil {
		return fmt.Errorf("failed to send message: %w", err)
	}

	return nil
}

// emit sends a Socket.IO event and waits for response
func (c *Client) call(event string, data interface{}) (map[string]interface{}, error) {
	c.mu.Lock()
	if !c.connected || c.wsConn == nil {
		c.mu.Unlock()
		return nil, fmt.Errorf("not connected")
	}

	// Get next event ID for callback
	c.eventID++
	eventID := c.eventID

	// Create response channel
	responseCh := make(chan SocketResponse, 1)
	c.respMu.Lock()
	c.responses[eventID] = responseCh
	c.respMu.Unlock()

	conn := c.wsConn
	c.mu.Unlock()

	// Clean up response channel on exit
	defer func() {
		c.respMu.Lock()
		delete(c.responses, eventID)
		c.respMu.Unlock()
		close(responseCh)
	}()

	// Create Socket.IO call message with callback: 42[ack_id]["event", data, callback]
	// The callback parameter is handled by the acknowledgment system
	eventData := []interface{}{event, data}
	eventJSON, err := json.Marshal(eventData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal event data: %w", err)
	}

	// Socket.IO message format for binary events with acknowledgments: 42[ack_id][json_array]
	message := fmt.Sprintf("42%d%s", eventID, string(eventJSON))

	// Use global mutex only for the WebSocket write operation
	globalWSMutex.Lock()
	err = conn.WriteMessage(websocket.TextMessage, []byte(message))
	globalWSMutex.Unlock()

	if err != nil {
		return nil, fmt.Errorf("failed to send message: %w", err)
	} // Wait for response with timeout
	select {
	case response := <-responseCh:
		if response.Error != "" {
			return nil, fmt.Errorf("server error: %s", response.Error)
		}

		// Parse response data
		if len(response.Data) > 0 {
			if result, ok := response.Data[0].(map[string]interface{}); ok {
				// Check for "ok" field in response
				if ok, exists := result["ok"]; exists {
					if okBool, isBool := ok.(bool); isBool && !okBool {
						msg := "unknown error"
						if msgStr, exists := result["msg"]; exists {
							if msgString, isString := msgStr.(string); isString {
								msg = msgString
							}
						}
						return nil, fmt.Errorf("API error: %s", msg)
					}
				}
				return result, nil
			}
		}

		return map[string]interface{}{}, nil

	case <-time.After(30 * time.Second):
		return nil, fmt.Errorf("timeout waiting for response")
	}
}

// makeRequest makes an authenticated request to the Uptime Kuma API
// GetMonitor retrieves a specific monitor by ID
func (c *Client) GetMonitor(id int) (*Monitor, error) {
	// Use cached monitor data from the monitorList event
	c.monitorsMu.RLock()
	monitorData := c.monitors
	c.monitorsMu.RUnlock()

	// Look for the specific monitor ID in the cached data
	for _, monitorDataInterface := range monitorData {
		if monitorMap, ok := monitorDataInterface.(map[string]interface{}); ok {
			// Parse ID
			if monitorID, ok := monitorMap["id"].(float64); ok && int(monitorID) == id {
				monitor := Monitor{}
				monitor.ID = int(monitorID)

				// Parse basic fields
				if name, ok := monitorMap["name"].(string); ok {
					monitor.Name = name
				}
				if monitorType, ok := monitorMap["type"].(string); ok {
					monitor.Type = monitorType
				}
				if url, ok := monitorMap["url"].(string); ok {
					monitor.URL = url
				}
				if hostname, ok := monitorMap["hostname"].(string); ok {
					monitor.Hostname = hostname
				}
				if port, ok := monitorMap["port"].(float64); ok {
					monitor.Port = int(port)
				}
				if interval, ok := monitorMap["interval"].(float64); ok {
					monitor.Interval = int(interval)
				}
				if timeout, ok := monitorMap["timeout"].(float64); ok {
					monitor.Timeout = int(timeout)
				}
				if retryInterval, ok := monitorMap["retryInterval"].(float64); ok {
					monitor.RetryInterval = int(retryInterval)
				}
				if resendInterval, ok := monitorMap["resendInterval"].(float64); ok {
					monitor.ResendInterval = int(resendInterval)
				}
				if maxretries, ok := monitorMap["maxretries"].(float64); ok {
					monitor.MaxRetries = int(maxretries)
				}
				if maxredirects, ok := monitorMap["maxredirects"].(float64); ok {
					monitor.MaxRedirects = int(maxredirects)
				}

				// Parse boolean fields - try both bool and float64 (0/1)
				if active, ok := monitorMap["active"].(bool); ok {
					monitor.Active = active
				} else if active, ok := monitorMap["active"].(float64); ok {
					monitor.Active = active == 1
				}
				if ignoreTls, ok := monitorMap["ignoreTls"].(bool); ok {
					monitor.IgnoreTLS = ignoreTls
				}
				if upsideDown, ok := monitorMap["upsideDown"].(bool); ok {
					monitor.UpsideDown = upsideDown
				}
				if followRedirect, ok := monitorMap["follow_redirect"].(bool); ok {
					monitor.FollowRedirect = followRedirect
				}

				if method, ok := monitorMap["method"].(string); ok {
					monitor.HTTPMethod = method
				}
				if body, ok := monitorMap["body"].(string); ok {
					monitor.Body = body
				}
				if basicAuthUser, ok := monitorMap["basic_auth_user"].(string); ok {
					monitor.BasicAuthUser = basicAuthUser
				}
				if basicAuthPass, ok := monitorMap["basic_auth_pass"].(string); ok {
					monitor.BasicAuthPass = basicAuthPass
				}

				// Parse accepted_statuscodes
				if statusCodes, ok := monitorMap["accepted_statuscodes"].([]interface{}); ok {
					var codes []string
					for _, codeInterface := range statusCodes {
						if codeStr, ok := codeInterface.(string); ok {
							codes = append(codes, codeStr)
						}
					}
					monitor.AcceptedStatusCodes = codes
				}

				// Parse notification_id_list
				if notificationIDs, ok := monitorMap["notification_id_list"].([]interface{}); ok {
					var ids []int
					for _, idInterface := range notificationIDs {
						if idStr, ok := idInterface.(string); ok {
							if id, err := strconv.Atoi(idStr); err == nil {
								ids = append(ids, id)
							}
						} else if idFloat, ok := idInterface.(float64); ok {
							ids = append(ids, int(idFloat))
						}
					}
					monitor.NotificationIDList = ids
				}

				return &monitor, nil
			}
		}
	}

	// Monitor not found in cache
	return nil, fmt.Errorf("monitor with ID %d not found", id)
}

// RefreshMonitors requests fresh monitor list from the server
func (c *Client) RefreshMonitors() error {
	// Request monitor list via Socket.IO
	err := c.emit("getMonitorList", nil)
	if err != nil {
		return fmt.Errorf("failed to request monitor list: %w", err)
	}

	// Wait for the monitorList event to update the cache
	time.Sleep(1 * time.Second)

	return nil
}

// GetMonitors retrieves all monitors
func (c *Client) GetMonitors() ([]Monitor, error) {
	// Force refresh to get latest data
	err := c.RefreshMonitors()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh monitors: %w", err)
	}

	// Use cached monitor data from the monitorList event
	c.monitorsMu.RLock()
	monitorData := c.monitors
	c.monitorsMu.RUnlock()

	// Parse the cached monitors into our Monitor struct
	var monitors []Monitor

	for _, monitorDataInterface := range monitorData {
		if monitorMap, ok := monitorDataInterface.(map[string]interface{}); ok {
			monitor := Monitor{}

			// Parse ID
			if id, ok := monitorMap["id"].(float64); ok {
				monitor.ID = int(id)
			}

			// Parse basic fields
			if name, ok := monitorMap["name"].(string); ok {
				monitor.Name = name
			}
			if monitorType, ok := monitorMap["type"].(string); ok {
				monitor.Type = monitorType
			}
			if url, ok := monitorMap["url"].(string); ok {
				monitor.URL = url
			}
			if hostname, ok := monitorMap["hostname"].(string); ok {
				monitor.Hostname = hostname
			}
			if port, ok := monitorMap["port"].(float64); ok {
				monitor.Port = int(port)
			}
			if interval, ok := monitorMap["interval"].(float64); ok {
				monitor.Interval = int(interval)
			}
			if timeout, ok := monitorMap["timeout"].(float64); ok {
				monitor.Timeout = int(timeout)
			}
			if active, ok := monitorMap["active"].(bool); ok {
				monitor.Active = active
			} else if active, ok := monitorMap["active"].(float64); ok {
				monitor.Active = active == 1
			}
			if method, ok := monitorMap["method"].(string); ok {
				monitor.HTTPMethod = method
			}
			if retryInterval, ok := monitorMap["retryInterval"].(float64); ok {
				monitor.RetryInterval = int(retryInterval)
			}
			if resendInterval, ok := monitorMap["resendInterval"].(float64); ok {
				monitor.ResendInterval = int(resendInterval)
			}
			if maxretries, ok := monitorMap["maxretries"].(float64); ok {
				monitor.MaxRetries = int(maxretries)
			}
			if maxredirects, ok := monitorMap["maxredirects"].(float64); ok {
				monitor.MaxRedirects = int(maxredirects)
			}
			if ignoreTls, ok := monitorMap["ignoreTls"].(bool); ok {
				monitor.IgnoreTLS = ignoreTls
			}
			if upsideDown, ok := monitorMap["upsideDown"].(bool); ok {
				monitor.UpsideDown = upsideDown
			}
			if followRedirect, ok := monitorMap["follow_redirect"].(bool); ok {
				monitor.FollowRedirect = followRedirect
			}
			if body, ok := monitorMap["body"].(string); ok {
				monitor.Body = body
			}
			if basicAuthUser, ok := monitorMap["basic_auth_user"].(string); ok {
				monitor.BasicAuthUser = basicAuthUser
			}
			if basicAuthPass, ok := monitorMap["basic_auth_pass"].(string); ok {
				monitor.BasicAuthPass = basicAuthPass
			}

			// Parse accepted_statuscodes
			if statusCodes, ok := monitorMap["accepted_statuscodes"].([]interface{}); ok {
				var codes []string
				for _, codeInterface := range statusCodes {
					if codeStr, ok := codeInterface.(string); ok {
						codes = append(codes, codeStr)
					}
				}
				monitor.AcceptedStatusCodes = codes
			}

			monitors = append(monitors, monitor)
		}
	}

	return monitors, nil
}

// CreateMonitor creates a new monitor using Socket.IO
func (c *Client) CreateMonitor(monitor *Monitor) (*Monitor, error) {
	// Set default accepted status codes if not provided
	acceptedStatusCodes := monitor.AcceptedStatusCodes
	if len(acceptedStatusCodes) == 0 {
		acceptedStatusCodes = []string{"200-299"}
	}

	// Build monitor data in the format expected by Uptime Kuma
	monitorData := map[string]interface{}{
		"type":                 monitor.Type,
		"name":                 monitor.Name,
		"url":                  monitor.URL,
		"hostname":             monitor.Hostname,
		"port":                 monitor.Port,
		"interval":             monitor.Interval,
		"timeout":              monitor.Timeout,
		"retryInterval":        monitor.RetryInterval,
		"resendInterval":       monitor.ResendInterval,
		"maxretries":           monitor.MaxRetries,
		"upsideDown":           monitor.UpsideDown,
		"maxredirects":         monitor.MaxRedirects,
		"accepted_statuscodes": acceptedStatusCodes,
		"method":               monitor.HTTPMethod,
		"body":                 monitor.Body,
		"headers":              "",
		"authMethod":           "",
		"basic_auth_user":      monitor.BasicAuthUser,
		"basic_auth_pass":      monitor.BasicAuthPass,
		"ignoreTls":            monitor.IgnoreTLS,
		"active":               monitor.Active,
		"notificationIDList":   map[string]interface{}{}, // Use empty object like existing monitors
		"httpBodyEncoding":     "json",
		"expiryNotification":   false,
		"dns_resolve_server":   "1.1.1.1",
		"dns_resolve_type":     "A",
		"proxyId":              nil,
		"mqttUsername":         "",
		"mqttPassword":         "",
		"mqttTopic":            "",
		"mqttSuccessMessage":   "",
		"keyword":              "",
		"invertKeyword":        false,
		"packetSize":           56,
	}

	// Add notification IDs if any are specified
	if len(monitor.NotificationIDList) > 0 {
		notificationIDListMap := make(map[string]interface{})
		for _, id := range monitor.NotificationIDList {
			notificationIDListMap[fmt.Sprintf("%d", id)] = true
		}
		monitorData["notificationIDList"] = notificationIDListMap
	} else {
		monitorData["notificationIDList"] = map[string]interface{}{}
	}

	// Remove empty fields (but keep important fields)
	for key, value := range monitorData {
		switch v := value.(type) {
		case string:
			if v == "" {
				delete(monitorData, key)
			}
		case int:
			if v == 0 && key != "port" && key != "interval" && key != "timeout" {
				delete(monitorData, key)
			}
		}
	}

	// Call the "add" API endpoint and wait for response
	response, err := c.call("add", monitorData)
	if err != nil {
		return nil, fmt.Errorf("failed to create monitor: %w", err)
	}

	// Extract monitor ID from response
	if monitorID, ok := response["monitorID"].(float64); ok {
		monitor.ID = int(monitorID)
	} else if msgData, ok := response["msg"].(map[string]interface{}); ok {
		// Some responses nest the ID in a msg object
		if monitorID, ok := msgData["monitorID"].(float64); ok {
			monitor.ID = int(monitorID)
		}
	}

	// If we still don't have an ID, try to find it from the cache as fallback
	if monitor.ID == 0 {
		time.Sleep(500 * time.Millisecond)
		monitors, err := c.GetMonitors()
		if err == nil {
			maxID := 0
			for _, m := range monitors {
				if m.Name == monitor.Name && m.URL == monitor.URL && m.ID > maxID {
					maxID = m.ID
				}
			}
			if maxID > 0 {
				monitor.ID = maxID
			}
		}
	}

	return monitor, nil
}

// UpdateMonitor updates an existing monitor
func (c *Client) UpdateMonitor(monitor *Monitor) (*Monitor, error) {
	// Set default accepted status codes if not provided
	acceptedStatusCodes := monitor.AcceptedStatusCodes
	if len(acceptedStatusCodes) == 0 {
		acceptedStatusCodes = []string{"200-299"}
	}

	// Build monitor data in the format expected by Uptime Kuma (same as create)
	monitorData := map[string]interface{}{
		"id":                   monitor.ID,
		"type":                 monitor.Type,
		"name":                 monitor.Name,
		"url":                  monitor.URL,
		"hostname":             monitor.Hostname,
		"port":                 monitor.Port,
		"interval":             monitor.Interval,
		"timeout":              monitor.Timeout,
		"retryInterval":        monitor.RetryInterval,
		"resendInterval":       monitor.ResendInterval,
		"maxretries":           monitor.MaxRetries,
		"upsideDown":           monitor.UpsideDown,
		"maxredirects":         monitor.MaxRedirects,
		"accepted_statuscodes": acceptedStatusCodes,
		"method":               monitor.HTTPMethod,
		"body":                 monitor.Body,
		"headers":              "",
		"authMethod":           "",
		"basic_auth_user":      monitor.BasicAuthUser,
		"basic_auth_pass":      monitor.BasicAuthPass,
		"ignoreTls":            monitor.IgnoreTLS,
		"active":               monitor.Active,
		"notificationIDList":   map[string]interface{}{}, // Use empty object like existing monitors
		"httpBodyEncoding":     "json",
		"expiryNotification":   false,
		"dns_resolve_server":   "1.1.1.1",
		"dns_resolve_type":     "A",
		"proxyId":              nil,
		"mqttUsername":         "",
		"mqttPassword":         "",
		"mqttTopic":            "",
		"mqttSuccessMessage":   "",
		"keyword":              "",
		"invertKeyword":        false,
		"packetSize":           56,
	}

	// Add notification IDs if any are specified
	if len(monitor.NotificationIDList) > 0 {
		notificationIDListMap := make(map[string]interface{})
		for _, id := range monitor.NotificationIDList {
			notificationIDListMap[fmt.Sprintf("%d", id)] = true
		}
		monitorData["notificationIDList"] = notificationIDListMap
	} else {
		monitorData["notificationIDList"] = map[string]interface{}{}
	}

	// Remove empty fields
	for key, value := range monitorData {
		switch v := value.(type) {
		case string:
			if v == "" && key != "id" {
				delete(monitorData, key)
			}
		case int:
			if v == 0 && key != "port" && key != "interval" && key != "timeout" && key != "id" {
				delete(monitorData, key)
			}
		}
	}

	// Call the "editMonitor" API endpoint
	err := c.emit("editMonitor", monitorData)
	if err != nil {
		return nil, fmt.Errorf("failed to update monitor: %w", err)
	}

	// Wait for the monitorList event to be updated (similar to create)
	time.Sleep(500 * time.Millisecond)

	return monitor, nil
}

// DeleteMonitor deletes a monitor
func (c *Client) DeleteMonitor(id int) error {
	// Call "deleteMonitor" API endpoint
	err := c.emit("deleteMonitor", id)
	if err != nil {
		return fmt.Errorf("failed to delete monitor: %w", err)
	}

	return nil
}

// RefreshNotifications requests fresh notification list from the server
func (c *Client) RefreshNotifications() error {
	// Request notification list via Socket.IO
	err := c.emit("getNotificationList", nil)
	if err != nil {
		return fmt.Errorf("failed to request notification list: %w", err)
	}

	// Wait for the notificationList event to update the cache
	time.Sleep(1 * time.Second)

	return nil
}

// GetNotifications retrieves all notifications from Uptime Kuma
func (c *Client) GetNotifications() ([]Notification, error) {
	// Force refresh to get latest data
	err := c.RefreshNotifications()
	if err != nil {
		return nil, fmt.Errorf("failed to refresh notifications: %w", err)
	}

	c.notificationsMu.RLock()
	notifications := make([]Notification, len(c.notificationCache))
	copy(notifications, c.notificationCache)
	c.notificationsMu.RUnlock()

	return notifications, nil
}

// parseNotificationMap converts a map to a Notification struct
func parseNotificationMap(notifMap map[string]interface{}) Notification {
	notification := Notification{}

	if id, ok := notifMap["id"].(float64); ok {
		notification.ID = int(id)
	}
	if name, ok := notifMap["name"].(string); ok {
		notification.Name = name
	}
	if active, ok := notifMap["active"].(bool); ok {
		notification.Active = active
	}
	if userID, ok := notifMap["userId"].(float64); ok {
		notification.UserID = int(userID)
	}
	if isDefault, ok := notifMap["isDefault"].(bool); ok {
		notification.IsDefault = isDefault
	}

	// Handle config as JSON string
	if configStr, ok := notifMap["config"].(string); ok {
		var configMap map[string]interface{}
		if err := json.Unmarshal([]byte(configStr), &configMap); err == nil {
			notification.Config = configMap

			// Extract type and other fields from config
			if notifType, exists := configMap["type"].(string); exists {
				notification.Type = notifType
			}
			if applyExisting, exists := configMap["applyExisting"].(bool); exists {
				notification.ApplyExisting = applyExisting
			}
		}
	}

	return notification
}

// GetNotification retrieves a specific notification by ID
func (c *Client) GetNotification(id int) (*Notification, error) {
	// Get all notifications and find the one with the specified ID
	notifications, err := c.GetNotifications()
	if err != nil {
		return nil, fmt.Errorf("failed to get notifications: %w", err)
	}

	for _, notification := range notifications {
		if notification.ID == id {
			// Return a copy to prevent modifications
			notifCopy := notification
			if notification.Config != nil {
				notifCopy.Config = make(map[string]interface{})
				for k, v := range notification.Config {
					notifCopy.Config[k] = v
				}
			}
			return &notifCopy, nil
		}
	}

	return nil, fmt.Errorf("notification with ID %d not found", id)
}

// CreateNotification creates a new notification
func (c *Client) CreateNotification(notification *Notification) (*Notification, error) {
	// Prepare notification data
	notificationData := map[string]interface{}{
		"name":          notification.Name,
		"type":          notification.Type,
		"isDefault":     notification.IsDefault,
		"applyExisting": notification.ApplyExisting,
	}

	// Add configuration parameters
	if notification.Config != nil {
		for key, value := range notification.Config {
			notificationData[key] = value
		}
	}

	// Use emit since notification APIs don't support callbacks
	err := c.emit("addNotification", notificationData)
	if err != nil {
		return nil, fmt.Errorf("failed to create notification: %w", err)
	}

	// Wait for the notificationList event to be updated
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		time.Sleep(500 * time.Millisecond)

		// Check if the notification appears in the cache
		c.notificationsMu.RLock()
		for _, notif := range c.notificationCache {
			if notif.Name == notification.Name {
				c.notificationsMu.RUnlock()
				// Return a copy
				result := notif
				return &result, nil
			}
		}
		c.notificationsMu.RUnlock()
	}

	return nil, fmt.Errorf("notification was created but not found in cache after %d retries", maxRetries)
}

// UpdateNotification updates an existing notification
func (c *Client) UpdateNotification(notification *Notification) (*Notification, error) {
	// Prepare notification data including the ID for updates
	notificationData := map[string]interface{}{
		"id":            notification.ID,
		"name":          notification.Name,
		"type":          notification.Type,
		"isDefault":     notification.IsDefault,
		"applyExisting": notification.ApplyExisting,
	}

	// Add configuration parameters
	if notification.Config != nil {
		for key, value := range notification.Config {
			notificationData[key] = value
		}
	}

	// Use emit for updates since notification APIs don't support callbacks
	err := c.emit("editNotification", notificationData)
	if err != nil {
		return nil, fmt.Errorf("failed to update notification: %w", err)
	}

	// Wait for the notificationList event to be updated
	maxRetries := 10
	for i := 0; i < maxRetries; i++ {
		time.Sleep(500 * time.Millisecond)

		// Check if the notification is updated in the cache
		c.notificationsMu.RLock()
		for _, notif := range c.notificationCache {
			if notif.ID == notification.ID {
				c.notificationsMu.RUnlock()
				// Return a copy
				result := notif
				return &result, nil
			}
		}
		c.notificationsMu.RUnlock()
	}

	return nil, fmt.Errorf("notification was updated but not found in cache after %d retries", maxRetries)
}

// DeleteNotification deletes a notification
func (c *Client) DeleteNotification(id int) error {
	// Use emit for deleteNotification since notification APIs don't support callbacks
	err := c.emit("deleteNotification", id)
	if err != nil {
		return fmt.Errorf("failed to delete notification: %w", err)
	}

	return nil
}

// TestNotification tests a notification configuration
func (c *Client) TestNotification(notification *Notification) error {
	// Prepare test data
	testData := map[string]interface{}{
		"name":          notification.Name,
		"type":          notification.Type,
		"isDefault":     notification.IsDefault,
		"applyExisting": notification.ApplyExisting,
	}

	// Add configuration parameters
	if notification.Config != nil {
		for key, value := range notification.Config {
			testData[key] = value
		}
	}

	err := c.emit("testNotification", testData)
	if err != nil {
		return fmt.Errorf("failed to test notification: %w", err)
	}

	return nil
}

// Close properly closes the client connection
func (c *Client) Close() {
	c.disconnect()
}
