// Copyright © 2020 Mike Berezin
//
// Use of this source code is governed by an MIT license.
// Details in the LICENSE file.

package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

const (
	defaultBaseURL   = "https://api.airtable.com/v0"
	defaultRateLimit = 4
)

// Client client for airtable api.
type Client struct {
	Client      *http.Client
	BaseURL     string
	APIKey      string
	rateLimiter <-chan time.Time
}

// New airtable client constructor
// your API KEY you can get on your account page
// https://airtable.com/account
func New(apiKey string) *Client {
	c := &Client{
		Client:  http.DefaultClient,
		APIKey:  apiKey,
		BaseURL: defaultBaseURL,
	}

	c.SetRateLimit(defaultRateLimit)

	return c
}

// SetRateLimit rate limit setter for custom usage
// Airtable limit is 5 requests per second (we use 4)
// https://airtable.com/{yourDatabaseID}/api/docs#curl/ratelimits
func (at *Client) SetRateLimit(rateLimit int) {
	at.rateLimiter = time.Tick(time.Second / time.Duration(rateLimit))
}

func (at *Client) waitForRateLimit() {
	<-at.rateLimiter
}

func (at *Client) get(ctx context.Context, db, table, recordID string, params url.Values, target any) error {
	return at.do(ctx, "GET", db, table, recordID, params, nil, target)
}

func (at *Client) post(ctx context.Context, db, table string, data, target any) error {
	return at.do(ctx, "POST", db, table, "", nil, data, target)
}

func (at *Client) delete(ctx context.Context, db, table string, recordIDs []string, target any) error {
	params := url.Values{}

	for _, recordID := range recordIDs {
		params.Add("records[]", recordID)
	}

	return at.do(ctx, "DELETE", db, table, "", params, nil, target)
}

func (at *Client) patch(ctx context.Context, db, table string, data, target any) error {
	return at.do(ctx, "PATCH", db, table, "", nil, data, target)
}

func (at *Client) put(ctx context.Context, db, table string, data, target any) error {
	return at.do(ctx, "PUT", db, table, "", nil, data, target)
}

func (at *Client) do(ctx context.Context, method, db, table, recordID string, params url.Values, data, target any) error {
	var err error

	at.waitForRateLimit()

	url := fmt.Sprintf("%s/%s/%s", at.BaseURL, db, table)

	if recordID != "" {
		url += fmt.Sprintf("/%s", recordID)
	}

	body := []byte{}

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return fmt.Errorf("marshalling message body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("cannot create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", at.APIKey))

	resp, err := at.Client.Do(req)
	if err != nil {
		return fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return makeHTTPClientError(url, resp)
	}

	dec := json.NewDecoder(resp.Body)

	err = dec.Decode(target)
	if err != nil {
		return fmt.Errorf("json decode failed: %w", err)
	}

	return nil
}
