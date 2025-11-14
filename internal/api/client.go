package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/SaneiyanReza/smsir-cli/internal/config"
)

const (
	// defaultHTTPTimeout is the default timeout for HTTP requests
	defaultHTTPTimeout = 30 * time.Second
)

// Client represents the SMS.ir API client
type Client struct {
	config     *config.Config
	httpClient *http.Client
	baseURL    string
}

// NewClient creates a new API client
func NewClient(cfg *config.Config) *Client {
	return &Client{
		config: cfg,
		httpClient: &http.Client{
			Timeout: defaultHTTPTimeout,
		},
		baseURL: cfg.BaseURL,
	}
}

// doRequest performs an HTTP request to the API
func (c *Client) doRequest(method, endpoint string, body interface{}) (*http.Response, error) {
	url := c.baseURL + endpoint

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal request body: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequest(method, url, reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-API-KEY", c.config.APIKey)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to perform request: %w", err)
	}

	return resp, nil
}

// parseResponse parses an API response into the given type
func parseResponse[T any](resp *http.Response) (*APIResponse[T], error) {
	defer resp.Body.Close()

	// Check for API errors
	if err := HandleAPIError(resp); err != nil {
		return nil, err
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var apiResp APIResponse[T]
	if err := json.Unmarshal(body, &apiResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	return &apiResp, nil
}

// GetCredit retrieves the current credit balance
func (c *Client) GetCredit() (*APIResponse[CreditResponse], error) {
	resp, err := c.doRequest("GET", "/credit", nil)
	if err != nil {
		return nil, err
	}
	return parseResponse[CreditResponse](resp)
}

// GetLines retrieves the list of available lines
func (c *Client) GetLines() (*APIResponse[LinesResponse], error) {
	resp, err := c.doRequest("GET", "/line", nil)
	if err != nil {
		return nil, err
	}
	return parseResponse[LinesResponse](resp)
}

// SendBulk sends bulk SMS messages
func (c *Client) SendBulk(req BulkSendRequest) (*APIResponse[BulkSendResponse], error) {
	resp, err := c.doRequest("POST", "/send/bulk", req)
	if err != nil {
		return nil, err
	}
	return parseResponse[BulkSendResponse](resp)
}

// HandleAPIError handles API errors based on status codes
func HandleAPIError(resp *http.Response) error {
	switch resp.StatusCode {
	case 200:
		return nil
	case 400:
		return fmt.Errorf("logical error: invalid request")
	case 401:
		return fmt.Errorf("authentication error: invalid API key")
	case 429:
		return fmt.Errorf("rate limit exceeded: please wait a moment")
	case 500:
		return fmt.Errorf("server error: unexpected error")
	default:
		return fmt.Errorf("unknown error with status code: %d", resp.StatusCode)
	}
}
