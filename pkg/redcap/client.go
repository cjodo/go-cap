package redcap

import (
	"net/http"
	"time"
)

const DefaultMaxRetries = 3

type Option func(*Client) error

func WithHttpClient(httpClient *http.Client) Option {
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

type Client struct {
	baseURL 		string
	token   		string
	httpClient *http.Client
	rateLimiter RateLimiter
	maxRetries  int
	retryDelay  time.Duration
}


func NewClient(baseURL, token string, opts ...Option ) (*Client, error) {
	hc := &http.Client{}

	c := &Client{
		baseURL: baseURL,
		token: token,
		httpClient: hc,
		rateLimiter: NewRateLimiterWithDefaultOpts(),
		maxRetries: DefaultMaxRetries,
		retryDelay: time.Millisecond * 2,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, err
		}
	}

	return c, nil
}

