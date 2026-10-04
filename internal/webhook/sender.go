// Package webhook sends webhook payloads to receivers over HTTP.
package webhook

import (
	"context"
	"errors"
	"net/http"
	"time"
)

// Result describes one HTTP delivery attempt.
type Result struct {
	// StatusCode is the HTTP status returned by the receiver.
	// It is 0 when no response was received.
	StatusCode int
	// Duration is the time from sending the request to reading the response
	// (or to the error). It is set even when Send returns an error.
	Duration time.Duration
	// Body is the beginning of the response body, at most MaxResponseBytes long.
	Body string
	// Truncated reports whether the response body was longer than MaxResponseBytes.
	Truncated bool
}

// Sender posts JSON payloads to webhook receivers.
type Sender struct {
	// Client is the HTTP client used for requests.
	Client *http.Client
	// Timeout limits the whole attempt: connecting, sending and reading the response.
	Timeout time.Duration
	// MaxResponseBytes limits how much of the response body is read.
	MaxResponseBytes int64
}

// Send posts payload to url as JSON.
//
// A response with any HTTP status (including 4xx and 5xx) is a successful
// attempt from the transport point of view: Send returns it in Result with a
// nil error. Send returns a non-nil error only when no complete response was
// received (invalid URL, connection refused, timeout, ...).
//
// TODO(Jacek): implement. See docs/tasks/01-sender.md.
func (s *Sender) Send(ctx context.Context, url string, payload []byte) (Result, error) {
	return Result{}, errors.New("Send: not implemented yet")
}
