# ADR 0001: Quality checks before every commit

**Date:** 2026-10-04 · **Status:** accepted

## Decision

A versioned git pre-commit hook (`.githooks/pre-commit`) runs
`scripts/check.sh`: gofmt, `go mod tidy -diff`, `go vet`, golangci-lint
(`standard` linters: errcheck, govet, ineffassign, staticcheck, unused) and
`go test ./...`. It is enabled per clone with `git config core.hooksPath .githooks`.

## Why

- Errors are caught on the laptop, before CI exists (CI comes later).
- One script for people, Claude sessions and, later, CI.
- golangci-lint instead of staticcheck alone: it already includes staticcheck,
  and errcheck teaches handling every error explicitly.

## Consequences

- golangci-lint is a dev tool installed with `go install`; it is not a module
  dependency, so `go.mod` stays stdlib-only.
- Red tests block a commit. For work in progress on a task whose tests are red
  on purpose: `SKIP_TESTS=1 git commit ...` (format, vet and lint still run).
- `git commit --no-verify` skips everything; avoid it.
