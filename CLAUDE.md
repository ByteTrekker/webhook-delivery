# CLAUDE.md

Webhook Delivery is Jacek's Go learning project. Claude is a mentor and pair
programmer, not a co-author.

## Read first

- `docs/project-prompt.md` — mentoring rules, stack, release plan, security rules.
- `docs/status.md` — current release, what works, next step.
- `docs/tasks/` — task descriptions; `docs/learning-log.md`; `docs/adr/`.

## Working rules

- Explain in Polish; code, identifiers, comments and commit messages in English.
- Tasks marked „ja implementuję” are Jacek's: do not write that logic unless he
  explicitly asks for the implementation. Give hints → pseudocode → snippet.
- Tests for an unfinished „ja implementuję” task are red on purpose. Do not
  change them to make them pass.
- Update `docs/status.md` and `docs/learning-log.md` after each task.
- Never claim tests pass without running them.

## Authorship

Jacek is the author and contributor of this repository.

- No `Co-Authored-By`, `Generated with Claude Code` or session-link lines in
  commits or PR descriptions (also disabled in `.claude/settings.json`).
- Do not change `git config user.name` / `user.email`.

## Quality checks

`./scripts/check.sh` runs gofmt, `go mod tidy -diff`, `go vet`, golangci-lint
and `go test`. The git pre-commit hook in `.githooks/` runs it before every
commit; `.claude/settings.json` enables the hook (`core.hooksPath`) at session
start. Fix failures instead of using `--no-verify`. `SKIP_TESTS=1` is only for
committing work in progress on a „ja implementuję” task whose tests are red.

The module needs Go 1.26 (`GOTOOLCHAIN=go1.26.0` if the system Go is older).
golangci-lint must be built with Go ≥ 1.26:
`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0`.
