package auth

import (
	"errors"
	"net/http"
	"strings"
)

// GetAPIKey extracts an API key from the headers of an HTTP request.
// Example:
// Authorization: ApiKey {insert apikey here}
func GetAPIKey(headers http.Header) (string, error) {
	val := headers.Get("Authorization")
	if val == "" {
		return "", errors.New("authorization header is missing")
	}
	vals := strings.Split(val, " ")
	if len(vals) != 2 {
		return "", errors.New("no authentication info found")
	}
	// Check if apikey is correct
	if vals[0] != "ApiKey" {
		return "", errors.New("malformed first part of auth reader ")
	}
	return vals[1], nil
}
