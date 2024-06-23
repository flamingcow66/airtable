package airtable

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func get[T any](ctx context.Context, c *Client, path string, params url.Values) (*T, error) {
	return do[T](ctx, c, "GET", path, params, nil)
}

func post[T any](ctx context.Context, c *Client, path string, data any) (*T, error) {
	return do[T](ctx, c, "POST", path, nil, data)
}

func del[T any](ctx context.Context, c *Client, path string, params url.Values) (*T, error) {
	return do[T](ctx, c, "DELETE", path, params, nil)
}

func patch[T any](ctx context.Context, c *Client, path string, data any) (*T, error) {
	return do[T](ctx, c, "PATCH", path, nil, data)
}

func put[T any](ctx context.Context, c *Client, path string, data any) (*T, error) {
	return do[T](ctx, c, "PUT", path, nil, data)
}

func do[T any](ctx context.Context, c *Client, method, path string, params url.Values, data any) (*T, error) {
	var err error

	c.waitForRateLimit()

	url := fmt.Sprintf("%s/%s", c.BaseURL, path)

	body := []byte{}

	if data != nil {
		body, err = json.Marshal(data)
		if err != nil {
			return nil, fmt.Errorf("marshalling message body: %w", err)
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewReader(body))
	if err != nil {
		return nil, fmt.Errorf("cannot create request: %w", err)
	}

	req.URL.RawQuery = params.Encode()

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", c.apiKey))

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("http request failed: %w", err)
	}

	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("http error: %d %s", resp.StatusCode, http.StatusText(resp.StatusCode))
	}

	dec := json.NewDecoder(resp.Body)

	target := new(T)

	err = dec.Decode(target)
	if err != nil {
		return nil, fmt.Errorf("json decode failed: %w", err)
	}

	return target, nil
}

func listAll[T any](ctx context.Context, c *Client, path string, params url.Values, key string, cb func(*T) error) ([]*T, error) {
	ret := []*T{}

	if params == nil {
		params = url.Values{}
	}

	for {
		resp, err := get[map[string]any](ctx, c, path, params)
		if err != nil {
			return nil, err
		}

		subresp, err := json.Marshal((*resp)[key])
		if err != nil {
			return nil, err
		}

		objs := []*T{}

		err = json.Unmarshal(subresp, &objs)
		if err != nil {
			return nil, err
		}

		if cb != nil {
			for _, obj := range objs {
				err = cb(obj)
				if err != nil {
					return nil, err
				}
			}
		}

		ret = append(ret, objs...)

		off, found := (*resp)["offset"]
		if !found {
			return ret, nil
		}

		params.Set("offset", off.(string))

	}
}
