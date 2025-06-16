package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path"
	"strings"
	"time"
)

// Client represents a Slurm REST API client.
type Client struct {
	BaseURL         string
	HTTPClient      *http.Client
	Auth            Authenticator
	lastResponseBody string
}

// Authenticator is an interface for different authentication methods.
type Authenticator interface {
	AddAuthHeaders(req *http.Request) error
}

// NewClient creates a new Slurm REST API client.
func NewClient(baseURL string, auth Authenticator) *Client {
	return &Client{
		BaseURL: baseURL,
		HTTPClient: &http.Client{
			Timeout: time.Second * 30,
		},
		Auth: auth,
	}
}

// Request makes an HTTP request to the Slurm REST API.
func (c *Client) Request(method, endpoint string, body interface{}, result interface{}) error {
	// Build the full URL
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	// Create the request body if provided
	var bodyReader io.Reader
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("error marshaling request body: %w", err)
		}
		bodyReader = bytes.NewReader(bodyBytes)
	}

	// Create the request
	req, err := http.NewRequest(method, u.String(), bodyReader)
	if err != nil {
		return fmt.Errorf("error creating request: %w", err)
	}

	// Set content type for requests with body
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// Add authentication headers
	if err := c.Auth.AddAuthHeaders(req); err != nil {
		return fmt.Errorf("error adding authentication headers: %w", err)
	}

	// Send the request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Read the response body
	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %w", err)
	}
	
	// Store the response body for later retrieval
	c.lastResponseBody = string(respBody)

	// Check for error status codes
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		// Special handling for 511 status code (Network Authentication Required)
		// This often happens with partition endpoints due to slurmdbd permissions
		if resp.StatusCode == 511 && strings.Contains(endpoint, "/partition") {
			// Print response body for debugging
			fmt.Fprintf(os.Stderr, "Response body: %s\n", string(respBody))
			
			// Try to unmarshal the response anyway, as it might contain partial data
			if result != nil && len(respBody) > 0 {
				if err := json.Unmarshal(respBody, result); err == nil {
					// Successfully unmarshaled, so we can continue with warnings
					fmt.Fprintf(os.Stderr, "Warning: Authentication error with slurmdbd (status code: %d)\n", resp.StatusCode)
					return nil
				} else {
					fmt.Fprintf(os.Stderr, "Error unmarshaling response: %v\n", err)
				}
			}
		}
		return fmt.Errorf("API error: %s (status code: %d)", string(respBody), resp.StatusCode)
	}

	// Parse the response if a result container was provided
	if result != nil && len(respBody) > 0 {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("error unmarshaling response: %w", err)
		}
	}

	return nil
}

// GetLastResponseBody returns the body of the last response received
func (c *Client) GetLastResponseBody() string {
	return c.lastResponseBody
}

// GenerateCurlCommand generates an equivalent curl command for a request.
func (c *Client) GenerateCurlCommand(method, endpoint string, body interface{}) (string, error) {
	// Build the full URL
	u, err := url.Parse(c.BaseURL)
	if err != nil {
		return "", fmt.Errorf("invalid base URL: %w", err)
	}
	u.Path = path.Join(u.Path, endpoint)

	// Create a dummy request to get the auth headers
	req, err := http.NewRequest(method, u.String(), nil)
	if err != nil {
		return "", fmt.Errorf("error creating request: %w", err)
	}

	// Add authentication headers
	if err := c.Auth.AddAuthHeaders(req); err != nil {
		return "", fmt.Errorf("error adding authentication headers: %w", err)
	}

	// Build the curl command
	curl := fmt.Sprintf("curl -X %s", method)

	// Add headers
	for k, v := range req.Header {
		curl += fmt.Sprintf(" -H '%s: %s'", k, v[0])
	}

	// Add body if provided
	if body != nil {
		bodyBytes, err := json.Marshal(body)
		if err != nil {
			return "", fmt.Errorf("error marshaling request body: %w", err)
		}
		curl += fmt.Sprintf(" -d '%s'", string(bodyBytes))
	}

	// Add URL
	curl += fmt.Sprintf(" '%s'", u.String())

	return curl, nil
}
