package whatsapp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
)

// WhatsappClient defines a client for communicating with the WhatsApp API
type WhatsappClient struct {
	URL   string
	Token string
}

// NewWhatsappClient initializes and returns a new instance of the WhatsApp client
func NewWhatsappClient(url, token string) *WhatsappClient {
	return &WhatsappClient{
		URL:   url,
		Token: token,
	}
}

// SendMessage sends a text message
func (wc *WhatsappClient) SendMessage(message, to string) error {
	// Construct the payload data according to WhatsApp API structure
	data := map[string]interface{}{
		"messaging_product": "whatsapp",
		"recipient_type":    "individual",
		"to":                to,
		"type":              "text",
		"text": map[string]interface{}{
			"preview_url": true,
			"body":        message,
		},
	}

	// Marshal the data into JSON
	payload, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("failed to marshal payload: %w", err)
	}

	// Create a new HTTP request
	req, err := http.NewRequest("POST", wc.URL, bytes.NewBuffer(payload))
	if err != nil {
		return fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+wc.Token)

	// Perform the HTTP request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body for logging/debugging
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %w", err)
	}

	// Check for non-2xx responses
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("unexpected status code: %d - %s", resp.StatusCode, string(body))
	}

	log.Printf("Message sent successfully: %s", string(body))
	return nil
}

// ValidateConfiguration checks if the client configuration is valid
func (wc *WhatsappClient) ValidateConfiguration() error {
	if wc.URL == "" || wc.Token == "" {
		return errors.New("invalid client configuration: URL and Token are required")
	}
	return nil
}
