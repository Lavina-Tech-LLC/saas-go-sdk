package saassupport

import (
	"errors"
	"fmt"
)

// Error represents an API error returned by the SaaS Support backend.
type Error struct {
	// Code is the HTTP status code (e.g. 400, 401, 404, 500).
	Code int `json:"code"`
	// Message is the human-readable error message from the API.
	Message string `json:"message"`
	// RawBody is the full response body for debugging.
	RawBody []byte `json:"-"`
}

func (e *Error) Error() string {
	return fmt.Sprintf("saassupport: %d %s", e.Code, e.Message)
}

// IsNotFound returns true if the error is a 404.
func IsNotFound(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == 404
}

// IsUnauthorized returns true if the error is a 401.
func IsUnauthorized(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == 401
}

// IsConflict returns true if the error is a 409.
func IsConflict(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == 409
}

// IsForbidden returns true if the error is a 403.
func IsForbidden(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == 403
}

// IsRateLimited returns true if the error is a 429.
func IsRateLimited(err error) bool {
	var e *Error
	return errors.As(err, &e) && e.Code == 429
}
