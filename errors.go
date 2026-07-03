package moyasar

import (
	"encoding/json"
	"fmt"
	"strings"
)

// APIError is returned when Moyasar responds with a non-2xx status.
//
// Moyasar error responses commonly include a type, a short message, and an
// optional errors object with field-level validation details. Body preserves
// the raw response body for diagnostics.
type APIError struct {
	// StatusCode is the HTTP response status code.
	StatusCode int `json:"-"`
	// Type is Moyasar's machine-readable error type.
	Type string `json:"type"`
	// Message is Moyasar's short human-readable error message.
	Message string `json:"message"`
	// Errors contains field-level validation errors when present.
	Errors map[string]interface{} `json:"errors"`
	// Body is the raw response body returned by Moyasar.
	Body []byte `json:"-"`
}

func (e *APIError) Error() string {
	parts := []string{fmt.Sprintf("moyasar: status %d", e.StatusCode)}
	if e.Type != "" {
		parts = append(parts, e.Type)
	}
	if e.Message != "" {
		parts = append(parts, e.Message)
	}
	return strings.Join(parts, ": ")
}

func decodeAPIError(statusCode int, body []byte) error {
	apiErr := &APIError{
		StatusCode: statusCode,
		Body:       append([]byte(nil), body...),
	}
	if len(body) == 0 {
		return apiErr
	}
	if err := json.Unmarshal(body, apiErr); err != nil {
		apiErr.Message = string(body)
	}
	return apiErr
}
