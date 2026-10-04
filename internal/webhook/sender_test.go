package webhook

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func newTestSender() *Sender {
	return &Sender{
		Client:           &http.Client{},
		Timeout:          time.Second,
		MaxResponseBytes: 1024,
	}
}

func TestSendSuccess(t *testing.T) {
	payload := []byte(`{"event":"order.created","id":42}`)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s, want POST", r.Method)
		}
		if ct := r.Header.Get("Content-Type"); ct != "application/json" {
			t.Errorf("Content-Type = %q, want application/json", ct)
		}
		body, _ := io.ReadAll(r.Body)
		if string(body) != string(payload) {
			t.Errorf("body = %s, want %s", body, payload)
		}
		w.Write([]byte(`{"received":true}`))
	}))
	defer srv.Close()

	res, err := newTestSender().Send(context.Background(), srv.URL, payload)
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if res.StatusCode != http.StatusOK {
		t.Errorf("StatusCode = %d, want 200", res.StatusCode)
	}
	if res.Body != `{"received":true}` {
		t.Errorf("Body = %q", res.Body)
	}
	if res.Truncated {
		t.Error("Truncated = true, want false")
	}
	if res.Duration <= 0 {
		t.Errorf("Duration = %v, want > 0", res.Duration)
	}
}

// An HTTP 500 is still a response: Send returns it without an error.
func TestSendServerError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "boom", http.StatusInternalServerError)
	}))
	defer srv.Close()

	res, err := newTestSender().Send(context.Background(), srv.URL, []byte(`{}`))
	if err != nil {
		t.Fatalf("Send() error = %v, want nil for an HTTP 500 response", err)
	}
	if res.StatusCode != http.StatusInternalServerError {
		t.Errorf("StatusCode = %d, want 500", res.StatusCode)
	}
	if !strings.Contains(res.Body, "boom") {
		t.Errorf("Body = %q, want it to contain %q", res.Body, "boom")
	}
}

func TestSendTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Read the body so the server notices when the client disconnects.
		io.Copy(io.Discard, r.Body)
		select {
		case <-time.After(2 * time.Second):
		case <-r.Context().Done(): // the client gave up
		}
	}))
	defer srv.Close()

	s := newTestSender()
	s.Timeout = 50 * time.Millisecond

	start := time.Now()
	res, err := s.Send(context.Background(), srv.URL, []byte(`{}`))
	if err == nil {
		t.Fatal("Send() error = nil, want a timeout error")
	}
	var netErr net.Error
	if !errors.As(err, &netErr) || !netErr.Timeout() {
		t.Errorf("Send() error = %v, want a timeout error", err)
	}
	if elapsed := time.Since(start); elapsed > time.Second {
		t.Errorf("Send took %v, the timeout was not applied", elapsed)
	}
	if res.StatusCode != 0 {
		t.Errorf("StatusCode = %d, want 0 when there is no response", res.StatusCode)
	}
	if res.Duration <= 0 {
		t.Errorf("Duration = %v, want > 0 even on error", res.Duration)
	}
}

func TestSendTruncatesLargeResponse(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(strings.Repeat("x", 10_000)))
	}))
	defer srv.Close()

	s := newTestSender()
	s.MaxResponseBytes = 100

	res, err := s.Send(context.Background(), srv.URL, []byte(`{}`))
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if len(res.Body) != 100 {
		t.Errorf("len(Body) = %d, want 100", len(res.Body))
	}
	if !res.Truncated {
		t.Error("Truncated = false, want true")
	}
}

func TestSendConnectionRefused(t *testing.T) {
	srv := httptest.NewServer(http.NotFoundHandler())
	url := srv.URL
	srv.Close() // nothing listens on this address anymore

	_, err := newTestSender().Send(context.Background(), url, []byte(`{}`))
	if err == nil {
		t.Fatal("Send() error = nil, want an error when the receiver is down")
	}
}
