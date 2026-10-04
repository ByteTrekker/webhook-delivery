#!/bin/sh
# Fails when a branch changes more than MAX_LINES lines (added + deleted)
# compared with main. go.sum is not counted.
#
# Usage: ./scripts/check-pr-size.sh [commit]   (default: HEAD)
set -u

MAX_LINES=2000
head=${1:-HEAD}

base=origin/main
git rev-parse -q --verify "$base" >/dev/null || base=main
if ! git rev-parse -q --verify "$base" >/dev/null; then
	echo "check-pr-size: no main branch found, skipping" >&2
	exit 0
fi

lines=$(git diff --numstat "$base...$head" -- . ':(exclude)go.sum' |
	awk '$1 != "-" { n += $1 + $2 } END { print n + 0 }')

if [ "$lines" -gt "$MAX_LINES" ]; then
	echo "This branch changes $lines lines vs $base (limit $MAX_LINES). Split it into smaller pull requests." >&2
	exit 1
fi
echo "PR size: $lines/$MAX_LINES lines changed vs $base"
