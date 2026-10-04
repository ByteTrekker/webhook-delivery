# ADR 0002: Pull request workflow

**Date:** 2026-10-04 · **Status:** accepted

## Decision

- No commits or pushes to `main`; all changes via pull requests.
- Jacek reviews and merges every pull request (`.github/CODEOWNERS`).
- Conventional Commits for every commit.
- At most 2000 changed lines per pull request (`go.sum` excluded).

Enforced locally by `.githooks/` (see `CONTRIBUTING.md`) and on GitHub by a
branch ruleset on `main` (Settings → Rules → Rulesets → New branch ruleset):

- Target: default branch; Enforcement: Active.
- Restrict deletions; Block force pushes.
- Require a pull request before merging, with Require review from Code Owners.
- Bypass list: empty.

## Why

- Jacek is the only contributor and the only one who decides what lands in `main`.
- Small PRs are reviewable in one sitting; this is a learning project, so each
  PR should be something Jacek can read and explain.
- Conventional Commits make history readable and allow a changelog later.

## Consequences

- GitHub does not let the author of a pull request approve it. PRs opened from
  Claude sessions are created with Jacek's GitHub account, so he is their
  author and cannot click "Approve". Requiring 1 approval would block every
  merge. Jacek's review is therefore his merge, and the required approval
  count stays at 0 until PRs come from a separate account.
- Rulesets on a private repository need GitHub Pro (or a public repository).
  Without them only the local hooks enforce these rules.
- Claude sessions are denied merging, approving and pushing to `main` in
  `.claude/settings.json`.
