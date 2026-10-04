package main

import (
	"embed"
	"encoding/json"
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/ByteTrekker/webhook-delivery/internal/webhook"
)

//go:embed templates/*.html
var templateFS embed.FS

var indexTmpl = template.Must(template.ParseFS(templateFS, "templates/index.html"))

const maxFormBytes = 64 << 10 // 64 KiB

type receiver struct {
	Name string
	URL  string
}

type application struct {
	logger    *slog.Logger
	sender    *webhook.Sender
	receivers []receiver
}

func (app *application) routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /{$}", app.index)
	mux.HandleFunc("POST /send", app.send)
	mux.HandleFunc("GET /healthz", healthz)

	mux.HandleFunc("POST /demo/ok", demoOK)
	mux.HandleFunc("POST /demo/fail", demoFail)
	mux.HandleFunc("POST /demo/slow", demoSlow)

	// Rejects cross-origin form submissions (CSRF) based on Sec-Fetch-Site / Origin.
	return http.NewCrossOriginProtection().Handler(mux)
}

// pageData is everything the index template renders.
type pageData struct {
	Receivers []receiver
	Selected  string
	Payload   string
	Error     string // validation error, shown above the form
	Sent      bool
	Result    webhook.Result
	SendErr   string // transport error from the sender
}

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, http.StatusOK, pageData{
		Receivers: app.receivers,
		Selected:  app.receivers[0].Name,
		Payload:   "{\n  \"event\": \"order.created\",\n  \"order_id\": 42\n}",
	})
}

func (app *application) send(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxFormBytes)
	if err := r.ParseForm(); err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			http.Error(w, "request too large", http.StatusRequestEntityTooLarge)
			return
		}
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	data := pageData{
		Receivers: app.receivers,
		Selected:  r.PostForm.Get("receiver"),
		Payload:   r.PostForm.Get("payload"),
	}

	rcv, ok := app.findReceiver(data.Selected)
	if !ok {
		data.Error = "Unknown receiver."
		app.render(w, http.StatusUnprocessableEntity, data)
		return
	}
	if !json.Valid([]byte(data.Payload)) {
		data.Error = "Payload is not valid JSON."
		app.render(w, http.StatusUnprocessableEntity, data)
		return
	}

	res, err := app.sender.Send(r.Context(), rcv.URL, []byte(data.Payload))
	data.Sent = true
	data.Result = res
	if err != nil {
		data.SendErr = err.Error()
	}

	app.logger.Info("webhook sent",
		"receiver", rcv.Name,
		"status", res.StatusCode,
		"duration_ms", res.Duration.Milliseconds(),
		"err", err,
	)
	app.render(w, http.StatusOK, data)
}

func (app *application) findReceiver(name string) (receiver, bool) {
	for _, rcv := range app.receivers {
		if rcv.Name == name {
			return rcv, true
		}
	}
	return receiver{}, false
}

func (app *application) render(w http.ResponseWriter, status int, data pageData) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	if err := indexTmpl.Execute(w, data); err != nil {
		app.logger.Error("render template", "err", err)
	}
}

func healthz(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("ok\n"))
}

// Demo receiver: a controlled endpoint to send webhooks to.

func demoOK(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"received":true}`))
}

func demoFail(w http.ResponseWriter, r *http.Request) {
	http.Error(w, `{"error":"demo failure"}`, http.StatusInternalServerError)
}

func demoSlow(w http.ResponseWriter, r *http.Request) {
	// Read the body so the server notices when the sender disconnects.
	io.Copy(io.Discard, r.Body)
	select {
	case <-time.After(5 * time.Second): // longer than the sender timeout
		w.Write([]byte(`{"received":true,"slow":true}`))
	case <-r.Context().Done(): // the sender gave up
	}
}
