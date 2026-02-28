package redcap

import (
	"net/http"
	"time"
)

type Option func(*Client) error

type Client struct {
	baseURL 		string
	token   		string
	httpClient *http.Client
	rateLimiter RateLimiter
	maxRetries  int
	retryDelay  time.Duration
}


func NewClient(baseURL, token string, opts ...Option ) (*Client, error) {
	return nil, nil
}

