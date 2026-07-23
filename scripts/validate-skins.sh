#!/usr/bin/env bash
# Verify Go catalog, TS manifest, and CSS imports stay in sync.
# Usage: scripts/validate-skins.sh

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
GO_FILE="$ROOT/internal/ui/skins.go"
TS_FILE="$ROOT/frontend/src/skins/manifest.ts"
CSS_INDEX="$ROOT/frontend/src/skins/index.css"
SKINS_DIR="$ROOT/frontend/src/skins"

fail=0

err() {
  echo "ERROR: $*" >&2
  fail=1
}

extract_go_ids() {
  sed -n '/catalog-begin/,/catalog-end/p' "$GO_FILE" \
    | grep 'ID:' \
    | sed -E 's/.*"([^"]+)".*/\1/' \
    | sort
}

extract_ts_ids() {
  sed -n '/manifest-begin/,/manifest-end/p' "$TS_FILE" \
    | grep "id:" \
    | sed -E "s/.*'([^']+)'.*/\1/" \
    | sort
}

extract_css_imports() {
  grep -E '@import "\./[^"]+\.css"' "$CSS_INDEX" \
    | sed -E 's/@import "\.\/([^"]+)\.css";/\1/' \
    | grep -v '^_' \
    | sort
}

go_ids="$(extract_go_ids)"
ts_ids="$(extract_ts_ids)"
css_ids="$(extract_css_imports)"

echo "Go catalog:  $(echo "$go_ids" | tr '\n' ' ')"
echo "TS manifest: $(echo "$ts_ids" | tr '\n' ' ')"
echo "CSS imports: $(echo "$css_ids" | tr '\n' ' ')"

if [[ "$go_ids" != "$ts_ids" ]]; then
  err "Go catalog and TS manifest IDs differ"
  comm -3 <(echo "$go_ids") <(echo "$ts_ids") | sed 's/^/  /' >&2 || true
fi

if [[ "$go_ids" != "$css_ids" ]]; then
  err "Go catalog and CSS imports differ"
  comm -3 <(echo "$go_ids") <(echo "$css_ids") | sed 's/^/  /' >&2 || true
fi

while IFS= read -r id; do
  [[ -z "$id" ]] && continue
  css_file="$SKINS_DIR/${id}.css"
  if [[ ! -f "$css_file" ]]; then
    err "Missing CSS file: frontend/src/skins/${id}.css"
    continue
  fi
  if ! grep -q "\[data-skin='${id}'\]" "$css_file"; then
    err "frontend/src/skins/${id}.css missing [data-skin='${id}'] selector"
  fi
  if ! grep -q "\[data-skin='${id}'\].dark" "$css_file"; then
    err "frontend/src/skins/${id}.css missing dark mode block"
  fi
done <<< "$go_ids"

if [[ "$fail" -ne 0 ]]; then
  echo ""
  echo "Fix with: make new-skin ID=<name> REGISTER=1"
  echo "Or edit: internal/ui/skins.go, frontend/src/skins/manifest.ts, frontend/src/skins/index.css"
  exit 1
fi

echo "OK — skin catalog, manifest, and CSS are in sync."
