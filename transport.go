package saassupport

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// envelope is the standard API response wrapper.
type envelope struct {
	Code    int             `json:"code"`
	Data    json.RawMessage `json:"data"`
	Message string          `json:"message"`
}

// Transport handles HTTP request construction and response parsing.
// It is shared across all module services.
type Transport struct {
	APIKey     string
	BaseURL    string
	HTTPClient *http.Client
	UserAgent  string
}

// Request sends an HTTP request and unmarshals the response data into dest.
func (t *Transport) Request(ctx context.Context, method, path string, body interface{}, dest interface{}) error {
	return t.RequestWithHeaders(ctx, method, path, body, dest, nil)
}

// RequestWithHeaders is like Request but also accepts extra headers (e.g. Authorization).
func (t *Transport) RequestWithHeaders(ctx context.Context, method, path string, body interface{}, dest interface{}, headers map[string]string) error {
	req, err := t.buildRequest(ctx, method, path, body)
	if err != nil {
		return fmt.Errorf("saassupport: building request: %w", err)
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := t.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("saassupport: sending request: %w", err)
	}
	defer resp.Body.Close()

	return t.parseResponse(resp, dest)
}

func (t *Transport) buildRequest(ctx context.Context, method, path string, body interface{}) (*http.Request, error) {
	url := t.BaseURL + path

	var bodyReader io.Reader
	if body != nil {
		b, err := json.Marshal(body)
		if err != nil {
			return nil, fmt.Errorf("marshaling body: %w", err)
		}
		bodyReader = bytes.NewReader(b)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bodyReader)
	if err != nil {
		return nil, err
	}

	req.Header.Set("X-API-Key", t.APIKey)
	req.Header.Set("User-Agent", t.UserAgent)
	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return req, nil
}

func (t *Transport) parseResponse(resp *http.Response, dest interface{}) error {
	rawBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("saassupport: reading response: %w", err)
	}

	var env envelope
	if err := json.Unmarshal(rawBody, &env); err != nil {
		// If we can't parse the envelope, return the raw status
		if resp.StatusCode >= 400 {
			return &Error{
				Code:    resp.StatusCode,
				Message: resp.Status,
				RawBody: rawBody,
			}
		}
		return fmt.Errorf("saassupport: parsing response: %w", err)
	}

	if env.Code >= 400 {
		return &Error{
			Code:    env.Code,
			Message: env.Message,
			RawBody: rawBody,
		}
	}

	if dest != nil && len(env.Data) > 0 {
		if err := json.Unmarshal(env.Data, dest); err != nil {
			return fmt.Errorf("saassupport: parsing data: %w", err)
		}
	}

	return nil
}
