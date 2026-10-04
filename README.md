# Webhook Delivery

A learning project: a service that sends webhooks, with a small panel to trigger and observe them.
Built step by step, release by release (see `docs/status.md`).

## Run locally

Requires Go 1.26 or newer (older Go downloads 1.26 automatically via `GOTOOLCHAIN`).

```sh
go run ./cmd/web     # http://localhost:8080
go test ./...
```

Configuration (environment variables):

| Variable | Default | Meaning |
|---|---|---|
| `ADDR` | `localhost:8080` | listen address |
| `DEMO_BASE_URL` | `http://$ADDR` | base URL of the built-in demo receiver |

## Routes

| Route | Purpose |
|---|---|
| `GET /` | form: pick a receiver, edit JSON |
| `POST /send` | send the webhook and show the result |
| `GET /healthz` | health check |
| `POST /demo/ok`, `/demo/fail`, `/demo/slow` | built-in demo receiver |

## Layout

- `cmd/web` — HTTP server, handlers, templates
- `internal/webhook` — sending webhooks over HTTP
- `docs/` — status, decisions, learning log, tasks
