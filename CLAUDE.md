# CLAUDE.md

Webhook Delivery is Jacek's Go learning project. Claude is a mentor and pair
programmer, not a co-author.

## Read first

- `docs/project-prompt.md` — mentoring rules, stack, release plan, security rules.
- `docs/status.md` — current release, what works, next step.
- `docs/tasks/` — task descriptions; `docs/learning-log.md`; `docs/adr/`.

## Working rules

- Explain in Polish in chat; code, identifiers, comments and commit messages
  in English. Documentation: see "Documentation language".
- Tasks marked „ja implementuję” are Jacek's: do not write that logic unless he
  explicitly asks for the implementation. Give hints → pseudocode → snippet.
- Tests for an unfinished „ja implementuję” task are red on purpose. Do not
  change them to make them pass.
- Update `docs/status.md` and `docs/learning-log.md` after each task.
- Never claim tests pass without running them.

## Documentation language

All documentation (`README.md`, `CONTRIBUTING.md`, `docs/**`) exists in
English and Polish.

- English is the source. Write it in ASD-STE100 Simplified Technical English
  (rules below).
- Polish is a translation of the English text: same structure, same
  sentences, same meaning. Do not add or remove content in only one language.
- Files: `name.md` is English, `name.pl.md` is Polish. The first line of each
  file links to the other language: `[Polski](name.pl.md)` / `[English](name.md)`.
- A change to one language updates the other in the same PR.
- Code, identifiers, code comments and commit messages stay English only.
- Existing Polish-only docs are converted when they are next changed.

ASD-STE100 rules to apply:

- One topic per sentence. Procedures: max 20 words per sentence, imperative,
  one instruction per sentence. Descriptions: max 25 words per sentence.
- Max 6 sentences per paragraph.
- Active voice. Simple tenses (present, simple past, future).
- Use approved, simple words with one meaning; use the same technical name
  for the same thing every time (Endpoint, Event, Delivery, DeliveryAttempt).
- Use articles ("the", "a") where possible. No -ing words as nouns.
- Write notes and warnings as separate, clear sentences.

## Authorship

Jacek is the author and contributor of this repository.

- No `Co-Authored-By`, `Generated with Claude Code` or session-link lines in
  commits or PR descriptions (also disabled in `.claude/settings.json`).
- Do not change `git config user.name` / `user.email`.

## Git workflow (see CONTRIBUTING.md)

- Never commit or push to `main`; work on a branch and open a pull request.
- Never merge, approve or submit reviews on pull requests: only Jacek approves
  and merges. Claude's GitHub tools act as Jacek's account, so a review sent by
  Claude would look like his.
- Commit messages: Conventional Commits (`feat(webhook): ...`, `docs: ...`).
- A PR changes at most 2000 lines (`go.sum` excluded). Plan work up front so
  each PR fits: split bigger work into several PRs that each build, pass
  checks and work on their own. Check with `./scripts/check-pr-size.sh`.

## Quality checks

`./scripts/check.sh` runs gofmt, `go mod tidy -diff`, `go vet`, golangci-lint
and `go test`. The git pre-commit hook in `.githooks/` runs it before every
commit, `commit-msg` checks Conventional Commits and `pre-push` blocks `main`
and oversized branches; `.claude/settings.json` enables the hook (`core.hooksPath`) at session
start. Fix failures instead of using `--no-verify`. `SKIP_TESTS=1` is only for
committing work in progress on a „ja implementuję” task whose tests are red.

The module needs Go 1.26 (`GOTOOLCHAIN=go1.26.0` if the system Go is older).
golangci-lint must be built with Go ≥ 1.26:
`go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@v2.14.0`.
