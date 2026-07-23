#!/usr/bin/env bash
# Fail if base/ or compiled binaries would be committed.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

fail=0

check_paths() {
  local label="$1"
  shift
  while IFS= read -r path; do
    [[ -z "$path" ]] && continue
    echo "ERROR: $label must not be in git: $path"
    fail=1
  done
}

# Tracked files (entire history on current branch is what we care about for CI;
# for pre-push we check index + working tree names that would be added)
check_paths "base/" < <(git ls-files 'base/*' 2>/dev/null || true)
check_paths "dist binary" < <(git ls-files 'dist/*' 2>/dev/null || true)
check_paths "root binary" < <(git ls-files 'go-snappymail' 'go-snappymail_*' 2>/dev/null || true)

# Staged additions
if git rev-parse --git-dir >/dev/null 2>&1; then
  check_paths "staged base/" < <(git diff --cached --name-only --diff-filter=A | rg '^base/' || true)
  check_paths "staged dist/" < <(git diff --cached --name-only --diff-filter=A | rg '^dist/' || true)
  check_paths "staged binary" < <(git diff --cached --name-only --diff-filter=A | rg '^(go-snappymail|go-snappymail_)' || true)
fi

if [[ "$fail" -ne 0 ]]; then
  echo ""
  echo "Fix: git rm --cached <path> and keep base/ + dist/ in .gitignore"
  exit 1
fi

echo "OK: no base/ or binaries tracked"
