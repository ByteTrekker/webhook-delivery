// Command web serves the webhook panel and the demo receiver.
package main

import (
	"cmp"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/ByteTrekker/webhook-delivery/internal/webhook"
)

func main() {
	logger := slog.New(slog.NewTextHandler(os.Stdout, nil))

	addr := cmp.Or(os.Getenv("ADDR"), "localhost:8080")
	demoBaseURL := cmp.Or(os.Getenv("DEMO_BASE_URL"), "http://"+addr)

	app := &application{
		logger: logger,
		sender: &webhook.Sender{
			Client:           &http.Client{},
			Timeout:          3 * time.Second,
			MaxResponseBytes: 2 << 10, // 2 KiB
		},
		// Receivers come from configuration, never from the user.
		receivers: []receiver{
			{Name: "demo-ok", URL: demoBaseURL + "/demo/ok"},
			{Name: "demo-fail", URL: demoBaseURL + "/demo/fail"},
			{Name: "demo-slow", URL: demoBaseURL + "/demo/slow"},
		},
	}

	srv := &http.Server{
		Addr:              addr,
		Handler:           app.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       10 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       60 * time.Second,
		ErrorLog:          slog.NewLogLogger(logger.Handler(), slog.LevelError),
	}

	logger.Info("starting server", "addr", "http://"+addr)
	if err := srv.ListenAndServe(); err != nil {
		logger.Error("server stopped", "err", err)
		os.Exit(1)
	}
}
