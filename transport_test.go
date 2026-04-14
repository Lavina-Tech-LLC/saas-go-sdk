package saassupport

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Transport must treat HTTP 4xx as an error even when the envelope has no
// "code" field — which is what the lvn.ResponseCamelCase backend emits.
// Without this, a backend 403 like {"data":"","message":"...","isOk":false}
// gets unmarshalled as success and the SDK tries to decode "" into the caller's
// destination type, producing a misleading JSON parse error.
func TestTransportDetectsEnvelopeErrors(t *testing.T) {
	cases := []struct {
		name       string
		status     int
		body       string
		wantCode   int
		wantErrMsg string
	}{
		{
			name:       "lvn envelope with isOk:false and empty data string",
			status:     403,
			body:       `{"data":"","message":"Not a member of this org","isOk":false}`,
			wantCode:   403,
			wantErrMsg: "Not a member of this org",
		},
		{
			name:       "lvn envelope with isOk:false and null data",
			status:     404,
			body:       `{"data":null,"message":"not found","isOk":false}`,
			wantCode:   404,
			wantErrMsg: "not found",
		},
		{
			name:       "legacy envelope with code field",
			status:     400,
			body:       `{"code":400,"data":null,"message":"bad request"}`,
			wantCode:   400,
			wantErrMsg: "bad request",
		},
		{
			name:       "HTTP 500 with unparseable body falls back to status",
			status:     500,
			body:       `internal server error plain text`,
			wantCode:   500,
			wantErrMsg: "500 Internal Server Error",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tc.status)
				_, _ = w.Write([]byte(tc.body))
			}))
			defer srv.Close()

			tr := &Transport{
				APIKey:     "test",
				BaseURL:    srv.URL,
				HTTPClient: srv.Client(),
				UserAgent:  "test",
			}

			var dest []struct{ Foo string }
			err := tr.Request(context.Background(), "GET", "/anything", nil, &dest)
			if err == nil {
				t.Fatalf("expected error, got nil")
			}
			var apiErr *Error
			if !errors.As(err, &apiErr) {
				t.Fatalf("expected *Error, got %T: %v", err, err)
			}
			if apiErr.Code != tc.wantCode {
				t.Errorf("Code = %d, want %d", apiErr.Code, tc.wantCode)
			}
			if apiErr.Message != tc.wantErrMsg {
				t.Errorf("Message = %q, want %q", apiErr.Message, tc.wantErrMsg)
			}
		})
	}
}

// Transport must successfully decode the success shape of the lvn envelope
// (no "code" field, isOk:true).
func TestTransportDecodesLvnSuccessEnvelope(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(200)
		_, _ = w.Write([]byte(`{"data":{"name":"alice"},"message":"","isOk":true}`))
	}))
	defer srv.Close()

	tr := &Transport{
		APIKey:     "test",
		BaseURL:    srv.URL,
		HTTPClient: srv.Client(),
		UserAgent:  "test",
	}

	var dest struct {
		Name string `json:"name"`
	}
	if err := tr.Request(context.Background(), "GET", "/anything", nil, &dest); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if dest.Name != "alice" {
		t.Errorf("Name = %q, want %q", dest.Name, "alice")
	}
}
