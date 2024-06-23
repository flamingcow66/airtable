package airtable

import (
	"fmt"
	"net/http"
	"os"
	"time"
)

const (
	defaultBaseURL   = "https://api.airtable.com/v0"
	defaultRateLimit = 4
)

type Client struct {
	Client      *http.Client
	BaseURL     string
	apiKey      string
	rateLimiter <-chan time.Time
}

func New(apiKey string) *Client {
	c := &Client{
		Client:  http.DefaultClient,
		BaseURL: defaultBaseURL,
		apiKey:  apiKey,
	}

	c.SetRateLimit(defaultRateLimit)

	return c
}

func NewFromEnv() (*Client, error) {
	apiKey := os.Getenv("AIRTABLE_TOKEN")
	if apiKey == "" {
		return nil, fmt.Errorf("please set $AIRTABLE_TOKEN")
	}

	return New(apiKey), nil
}

func (c *Client) SetRateLimit(rateLimit int) {
	c.rateLimiter = time.Tick(time.Second / time.Duration(rateLimit))
}

func (c *Client) waitForRateLimit() {
	<-c.rateLimiter
}
