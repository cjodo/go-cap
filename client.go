package redcap

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	DefaultMaxRetries = 3
	DefaultRetryDelay = time.Second
	MaxRetryDelay     = 30 * time.Second
)

type Option func(*Client) error

func WithHTTPClient(httpClient *http.Client) Option {
	return func(c *Client) error {
		c.httpClient = httpClient
		return nil
	}
}

func WithMaxRetries(n int) Option {
	return func(c *Client) error {
		c.maxRetries = n
		return nil
	}
}

func WithRetryDelay(d time.Duration) Option {
	return func(c *Client) error {
		c.retryDelay = d
		return nil
	}
}

func WithRateLimiter(r RateLimiter) Option {
	return func(c *Client) error {
		c.rateLimiter = r
		return nil
	}
}

type Client struct {
	baseURL     string
	token       string
	httpClient  *http.Client
	rateLimiter RateLimiter
	maxRetries  int
	retryDelay  time.Duration
	logLevel    string
}

func NewClient(baseURL, token string, opts ...Option) (*Client, error) {
	hc := &http.Client{
		Timeout: 30 * time.Second,
	}

	c := &Client{
		baseURL:     baseURL,
		token:       token,
		httpClient:  hc,
		rateLimiter: NewRateLimiterWithDefaultOpts(),
		maxRetries:  DefaultMaxRetries,
		retryDelay:  DefaultRetryDelay,
		logLevel:    "info",
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

// Request makes a REDCap API request with retry logic and rate limiting.
// The content parameter specifies the API endpoint (e.g., "record", "metadata").
// Additional params are merged with the standard parameters (token, content).
func (c *Client) Request(ctx context.Context, content string, params map[string]string) ([]byte, error) {
	var lastErr error

	for attempt := 0; attempt <= c.maxRetries; attempt++ {
		if attempt > 0 {
			delay := c.calculateBackoff(attempt)
			select {
			case <-ctx.Done():
				return nil, ctx.Err()
			case <-time.After(delay):
			}
		}

		// Wait for rate limiter
		if err := c.rateLimiter.Wait(ctx); err != nil {
			if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
				return nil, err
			}
			continue
		}

		body, err := c.doRequest(ctx, content, params)
		if err == nil {
			return body, nil
		}

		lastErr = err

		// Check if we should retry
		var redcapErr *Error
		if errors.As(err, &redcapErr) {
			if !redcapErr.IsRetryable() {
				return nil, err
			}
			// Rate limited - maybe increase delay
			if redcapErr.Code == ErrCodeRateLimit {
				c.rateLimiter.SetRate(c.rateLimiter.GetRate() * 0.8)
			}
		} else if errors.Is(err, context.Canceled) || errors.Is(err, context.DeadlineExceeded) {
			return nil, err
		}
		// Network errors are retryable
	}

	return nil, lastErr
}

// doRequest performs a single HTTP request to the REDCap API.
func (c *Client) doRequest(ctx context.Context, content string, params map[string]string) ([]byte, error) {
	form := url.Values{}
	form.Add("token", c.token)
	form.Add("returnFormat", "json")
	if content != "" {
		form.Add("content", content)
	}
	for k, v := range params {
		if v != "" {
			form.Add(k, v)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, c.baseURL, strings.NewReader(form.Encode()))
	if err != nil {
		return nil, fmt.Errorf("creating request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Accept", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("reading response body: %w", err)
	}

	// Check for REDCap API errors
	if resp.StatusCode != http.StatusOK {
		return nil, c.parseError(resp.StatusCode, body)
	}

	// Check for error in response body (REDCap sometimes returns errors as JSON with 200)
	var apiErr struct {
		Error string `json:"error"`
	}
	if json.Unmarshal(body, &apiErr); apiErr.Error != "" {
		return nil, &Error{
			Code:       ErrCodeInvalidRequest,
			Message:    apiErr.Error,
			StatusCode: resp.StatusCode,
		}
	}

	return body, nil
}

// parseError converts an HTTP response into a redcap.Error.
func (c *Client) parseError(statusCode int, body []byte) *Error {
	code := ErrCodeUnknown
	message := string(body)

	// Try to parse REDCap's error format
	var redcapErr struct {
		Error   string `json:"error"`
		Message string `json:"message"`
	}
	if json.Unmarshal(body, &redcapErr); redcapErr.Error != "" {
		message = redcapErr.Error
	} else if redcapErr.Message != "" {
		message = redcapErr.Message
	}

	switch statusCode {
	case http.StatusBadRequest:
		code = ErrCodeInvalidRequest
	case http.StatusUnauthorized:
		code = ErrCodeUnauthorized
	case http.StatusForbidden:
		code = ErrCodeForbidden
	case http.StatusNotFound:
		code = ErrCodeNotFound
	case http.StatusTooManyRequests:
		code = ErrCodeRateLimit
	case 500, 502, 503, 504:
		code = ErrCodeServerError
	}

	return &Error{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}

// calculateBackoff returns the delay for exponential backoff with jitter.
func (c *Client) calculateBackoff(attempt int) time.Duration {
	delay := c.retryDelay * time.Duration(math.Pow(2, float64(attempt-1)))
	if delay > MaxRetryDelay {
		delay = MaxRetryDelay
	}
	// Add jitter (±25%)
	jitter := float64(delay) * 0.25 * (rand.Float64()*2 - 1)
	return delay + time.Duration(jitter)
}
