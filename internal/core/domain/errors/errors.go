package errors

import "errors"

var (
	// ErrToolNotFound is returned when a requested tool cannot be found.
	ErrToolNotFound = errors.New("tool not found")

	// ErrToolExecutionFailed is returned when a tool execution fails.
	ErrToolExecutionFailed = errors.New("tool execution failed")

	// ErrInvalidArguments is returned when the provided arguments for a tool are invalid.
	ErrInvalidArguments = errors.New("invalid arguments")

	// ErrUnauthorized is returned when the request is unauthorized.
	ErrUnauthorized = errors.New("unauthorized")

	// ErrAPIKeyMissing is returned when the API key is missing in the request.
	ErrAPIKeyMissing = errors.New("API key is missing")

	// ErrVerticalMissing is returned when the vertical ID is missing in the request.
	ErrVerticalMissing = errors.New("vertical ID is missing")

	// ErrInternalServer is returned when an internal server error occurs.
	ErrInternalServer = errors.New("internal server error")

	// ErrInvalidTimeout is returned when the configured timeout is invalid.
	ErrInvalidTimeout = errors.New("invalid timeout configuration")

	// ErrNetworkTimeout is returned when a network timeout occurs while communicating with the API.
	ErrNetworkTimeout = errors.New("network timeout")

	// ErrInvalidResponse is returned when the response from the API is invalid or malformed.
	ErrInvalidResponse = errors.New("invalid response from API")
)
