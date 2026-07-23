#!/usr/bin/env bash
# Reject base/, binaries, and sensitive data from git.
set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

fail=0

err() {
  echo "ERROR: $*"
  fail=1
}

# --- paths that must never be tracked ---
while IFS= read -r path; do
  [[ -z "$path" ]] && continue
  err "must not be tracked: $path"
done < <(git ls-files 2>/dev/null | rg -i '^(base/|dist/)' || true)

while IFS= read -r path; do
  [[ -z "$path" ]] && continue
  err "must not be tracked: $path"
done < <(git ls-files 2>/dev/null | rg -i '^(go-snappymail|go-snappymail_)' || true)

while IFS= read -r path; do
  [[ -z "$path" ]] && continue
  err "env file must not be tracked (use .env.example): $path"
done < <(git ls-files 2>/dev/null | rg '\.env$' | rg -v '\.env\.example$' || true)

while IFS= read -r path; do
  [[ -z "$path" ]] && continue
  err "private key must not be tracked: $path"
done < <(git ls-files 2>/dev/null | rg -i 'id_rsa$|id_ed25519$|\.pem$|\.p12$|private_key' || true)

# --- staged additions ---
if git rev-parse --git-dir >/dev/null 2>&1; then
  while IFS= read -r path; do
    [[ -z "$path" ]] && continue
    err "staged path blocked: $path"
  done < <(git diff --cached --name-only --diff-filter=A 2>/dev/null | rg '^(base/|dist/)' || true)

  while IFS= read -r path; do
    [[ -z "$path" ]] && continue
    err "staged env file blocked: $path"
  done < <(git diff --cached --name-only --diff-filter=A 2>/dev/null | rg '\.env$' | rg -v '\.env\.example$' || true)
fi

# --- secret patterns in staged diff ---
if git rev-parse --git-dir >/dev/null 2>&1; then
  staged="$(git diff --cached 2>/dev/null || true)"
  if [[ -n "$staged" ]]; then
    if echo "$staged" | rg -q 'BEGIN (OPENSSH|RSA|EC) PRIVATE KEY'; then
      err "staged diff contains a private key block"
    fi
    if echo "$staged" | rg -q 'ghp_[A-Za-z0-9]{20,}|gho_[A-Za-z0-9]{20,}|github_pat_[A-Za-z0-9_]{20,}'; then
      err "staged diff contains a GitHub token"
    fi
    if echo "$staged" | rg -q 'AKIA[0-9A-Z]{16}'; then
      err "staged diff contains an AWS access key"
    fi
  fi
fi

# --- scan tracked files for real tokens (not lab placeholders) ---
while IFS= read -r file; do
  [[ -z "$file" || ! -f "$file" ]] && continue
  if rg -q 'ghp_[A-Za-z0-9]{20,}|gho_[A-Za-z0-9]{20,}|github_pat_[A-Za-z0-9_]{20,}' "$file" 2>/dev/null; then
    err "GitHub token pattern in tracked file: $file"
  fi
  if rg -q 'BEGIN (OPENSSH|RSA|EC) PRIVATE KEY' "$file" 2>/dev/null; then
    err "private key in tracked file: $file"
  fi
done < <(git ls-files 2>/dev/null || true)

if [[ "$fail" -ne 0 ]]; then
  echo ""
  echo "Secrets belong in .env (gitignored) or local config.toml — never in the repo."
  echo "Lab defaults: copy docker/.env.example → docker/.env"
  exit 1
fi

echo "OK: no base/, binaries, .env, or secret patterns in git"
