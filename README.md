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

## Contributing

Branches + pull requests only (no commits to `main`), Conventional Commits,
at most 2000 changed lines per PR. Details: `CONTRIBUTING.md`.

## Code quality checks

Every commit runs gofmt, `go mod tidy -diff`, `go vet`, golangci-lint and
`go test` (see `scripts/check.sh`, decision in `docs/adr/0001-quality-checks-before-commit.md`).

One-time setup after cloning:

```sh
git config core.hooksPath .githooks
go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0
```

Make sure `$(go env GOPATH)/bin` is on your `PATH`. Run the checks by hand with
`./scripts/check.sh`. To commit work in progress while a task's tests are
still red: `SKIP_TESTS=1 git commit ...` (formatting, vet and lint still run).

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
- `scripts/`, `.githooks/` — quality checks, commit message and PR size checks
- `docs/` — status, decisions, learning log, tasks
- `CLAUDE.md`, `.claude/` — rules for Claude sessions
