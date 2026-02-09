package exitcode

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

// Exit codes aligned with docs/USAGE.md.
const (
	ExitOK             = 0
	ExitAuthFailure    = 10
	ExitPermission     = 11
	ExitNotFound       = 12
	ExitBadArgs        = 13
	ExitTransport      = 14
	ExitRateLimit      = 15
	ExitLocalWriteErr  = 16
	ExitUserAbort      = 17
)

// Category classifies an error for exit-code mapping.
type Category int

const (
	CatAuthFailure   Category = iota // 401 or missing token
	CatPermission                     // 403 permission denied
	CatNotFound                       // 404
	CatBadArgs                        // invalid CLI arguments
	CatTransport                      // timeout / network
	CatRateLimit                      // 403 rate-limited
	CatLocalWriteErr                  // local I/O failure
	CatUserAbort                      // user cancelled operation
)

// CLIError is the single error type that reaches main and maps to an exit code.
type CLIError struct {
	Cat     Category
	Message string
	Err     error
}

func (e *CLIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *CLIError) Unwrap() error { return e.Err }

// ExitCode returns the numeric exit code for this error.
func (e *CLIError) ExitCode() int {
	switch e.Cat {
	case CatAuthFailure:
		return ExitAuthFailure
	case CatPermission:
		return ExitPermission
	case CatNotFound:
		return ExitNotFound
	case CatBadArgs:
		return ExitBadArgs
	case CatTransport:
		return ExitTransport
	case CatRateLimit:
		return ExitRateLimit
	case CatLocalWriteErr:
		return ExitLocalWriteErr
	case CatUserAbort:
		return ExitUserAbort
	default:
		return 1
	}
}

// Convenience constructors.

func NewAuthFailure(msg string, err error) *CLIError {
	return &CLIError{Cat: CatAuthFailure, Message: msg, Err: err}
}

func NewPermission(msg string, err error) *CLIError {
	return &CLIError{Cat: CatPermission, Message: msg, Err: err}
}

func NewRateLimit(msg string, err error) *CLIError {
	return &CLIError{Cat: CatRateLimit, Message: msg, Err: err}
}

func NewTransport(msg string, err error) *CLIError {
	return &CLIError{Cat: CatTransport, Message: msg, Err: err}
}

func NewBadArgs(msg string, err error) *CLIError {
	return &CLIError{Cat: CatBadArgs, Message: msg, Err: err}
}

func NewNotFound(msg string, err error) *CLIError {
	return &CLIError{Cat: CatNotFound, Message: msg, Err: err}
}

func NewLocalWriteErr(msg string, err error) *CLIError {
	return &CLIError{Cat: CatLocalWriteErr, Message: msg, Err: err}
}

func NewUserAbort(msg string, err error) *CLIError {
	return &CLIError{Cat: CatUserAbort, Message: msg, Err: err}
}

// ClassifyHTTP converts an HTTP status code plus optional context into a CLIError.
// rateLimited should be true when response headers indicate rate limiting.
func ClassifyHTTP(status int, rateLimited bool, body string) *CLIError {
	switch {
	case status == 401:
		return NewAuthFailure("authentication failed: invalid or expired token", nil)
	case status == 403 && rateLimited:
		return NewRateLimit("rate limit exceeded", nil)
	case status == 403:
		return NewPermission("permission denied: insufficient token scope", nil)
	case status == 404:
		return &CLIError{Cat: CatNotFound, Message: "not found: " + body}
	default:
		return &CLIError{Cat: CatTransport, Message: fmt.Sprintf("unexpected HTTP %d: %s", status, body)}
	}
}

// ClassifyTransportErr converts a Go network/timeout error into a CLIError.
func ClassifyTransportErr(err error) *CLIError {
	if err == nil {
		return nil
	}
	var netErr net.Error
	if errors.As(err, &netErr) && netErr.Timeout() {
		return NewTransport("request timed out", err)
	}
	if isNetworkErr(err) {
		return NewTransport("network error", err)
	}
	return NewTransport("request failed", err)
}

func isNetworkErr(err error) bool {
	var opErr *net.OpError
	if errors.As(err, &opErr) {
		return true
	}
	msg := err.Error()
	return strings.Contains(msg, "connection refused") ||
		strings.Contains(msg, "no such host") ||
		strings.Contains(msg, "dial tcp")
}
