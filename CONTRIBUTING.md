# Contributing

Rules for everyone working in this repo, people and Claude sessions alike.

## Branches and pull requests

- Never commit or push to `main`. Every change goes through a pull request
  from a branch (`feat/...`, `fix/...`, `docs/...`, `chore/...`).
- Only Jacek (@ByteTrekker) approves and merges pull requests. Claude
  sessions never submit reviews or merge. `.github/CODEOWNERS` makes him the reviewer of
  every file.
- A pull request changes at most **2000 lines** (added + deleted, `go.sum`
  not counted). Plan work so each PR fits; split bigger work into a series of
  PRs that each work on their own. Check with `./scripts/check-pr-size.sh`.

## Documentation: English (ASD-STE100) and Polish

- Every document (`README.md`, `CONTRIBUTING.md`, `docs/**`) has an English
  version (`name.md`) and a Polish version (`name.pl.md`).
- Write the English version in ASD-STE100 Simplified Technical English.
- The Polish version is a translation of the English version.
- Update both versions in the same pull request.
- The full rules are in `CLAUDE.md`, section "Documentation language".

## Commit messages: Conventional Commits

Every commit follows [Conventional Commits 1.0](https://www.conventionalcommits.org/en/v1.0.0/):

```
<type>[(scope)][!]: <description>

[body]
```

- Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`,
  `build`, `ci`, `chore`, `revert`.
- Scope is optional, lowercase, e.g. a package: `feat(webhook): ...`.
- `!` marks a breaking change.
- Description in English, imperative, no trailing period:
  `fix(web): limit request body to 64 KiB`.

## Enforced locally by git hooks

Enable once per clone: `git config core.hooksPath .githooks`.

| Hook | Checks |
|---|---|
| `pre-commit` | not on `main`; `scripts/check.sh` (gofmt, mod tidy, vet, golangci-lint, tests) |
| `commit-msg` | Conventional Commits format |
| `pre-push` | not pushing to `main`; PR size ≤ 2000 lines |

## Enforced on GitHub (repository settings)

Local hooks can be skipped, so `main` is also protected on GitHub; see
`docs/adr/0002-pull-request-workflow.md` for the exact settings.
