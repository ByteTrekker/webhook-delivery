# Status

**Release:** 1 — „Wyślij i zobacz” (in progress)
**Environment:** local only (no VPS/domain yet; deployment postponed by decision on 2026-10-04)

## Works

- Form at `/`, receivers from config (`demo-ok`, `demo-fail`, `demo-slow`).
- `/send` validates input (64 KiB limit, valid JSON, known receiver) and renders the result.
- Demo receiver, `/healthz`, `slog` logs, server timeouts, cross-origin protection.

## Not yet

- `webhook.Sender.Send` — task 1 (Jacek), see `docs/tasks/01-sender.md`.
- Dockerfile, Compose, Caddy, Basic Auth — before the first public deploy.

## Next step

Task 1: implement `Sender.Send` until `go test ./...` passes.
