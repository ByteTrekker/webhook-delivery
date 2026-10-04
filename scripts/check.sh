#!/bin/sh
# Runs every code-quality check. Used by the git pre-commit hook
# (.githooks/pre-commit); you can also run it by hand: ./scripts/check.sh
#
# SKIP_TESTS=1 skips `go test` (for committing work in progress while a
# task's tests are still red). Formatting, vet and lint always run.
set -u

GOLANGCI_LINT_VERSION=v2.14.0

cd "$(git rev-parse --show-toplevel)" || exit 1

failed=0
step() {
	printf '\n==> %s\n' "$1"
}
fail() {
	printf 'FAIL: %s\n' "$1"
	failed=1
}

step "gofmt"
unformatted=$(gofmt -l .)
if [ -n "$unformatted" ]; then
	printf '%s\n' "$unformatted"
	fail "files above are not formatted (fix: gofmt -w .)"
fi

step "go mod tidy"
go mod tidy -diff || fail "go.mod/go.sum are not tidy (fix: go mod tidy)"

step "go vet"
go vet ./... || fail "go vet"

step "golangci-lint"
if command -v golangci-lint >/dev/null 2>&1; then
	case "$(golangci-lint version 2>&1)" in
	*"version ${GOLANGCI_LINT_VERSION#v} "*) ;;
	*) echo "warning: expected golangci-lint $GOLANGCI_LINT_VERSION (install below)" ;;
	esac
	golangci-lint run ./... || fail "golangci-lint"
else
	fail "golangci-lint is not installed. Install it once with:
  go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$GOLANGCI_LINT_VERSION"
fi

step "go test"
if [ "${SKIP_TESTS:-}" = "1" ]; then
	echo "skipped (SKIP_TESTS=1)"
else
	go test ./... || fail "go test (commit work in progress with: SKIP_TESTS=1 git commit ...)"
fi

if [ "$failed" -ne 0 ]; then
	printf '\nChecks failed, commit aborted.\n'
	exit 1
fi
printf '\nAll checks passed.\n'
