#!/usr/bin/env bash
# Scaffold a new webmail skin and optionally register it everywhere.
# Usage:
#   scripts/new-skin.sh <skin-id> [--register]
#   make new-skin ID=acme
#   make new-skin ID=acme REGISTER=1

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
ID="${1:-}"
REGISTER="${REGISTER:-0}"
if [[ "${2:-}" == "--register" ]]; then
  REGISTER=1
fi

if [[ -z "$ID" ]]; then
  echo "Usage: make new-skin ID=<skin-id> [REGISTER=1]"
  echo "Example: make new-skin ID=acme REGISTER=1"
  exit 1
fi

if [[ ! "$ID" =~ ^[a-z][a-z0-9-]*$ ]]; then
  echo "Error: skin id must be lowercase letters, digits, or hyphens (start with a letter)."
  echo "Got: $ID"
  exit 1
fi

TARGET="$ROOT/frontend/src/skins/${ID}.css"
TEMPLATE="$ROOT/frontend/src/skins/_template.css"
GO_FILE="$ROOT/internal/ui/skins.go"
TS_FILE="$ROOT/frontend/src/skins/manifest.ts"
CSS_INDEX="$ROOT/frontend/src/skins/index.css"

if [[ -f "$TARGET" ]]; then
  echo "Error: $TARGET already exists."
  exit 1
fi

# Human label: acme-corp → Acme Corp
label="$(echo "$ID" | sed -E 's/-/ /g' | awk '{for(i=1;i<=NF;i++) $i=toupper(substr($i,1,1)) tolower(substr($i,2)); print}')"

sed "s/__SKIN_ID__/${ID}/g" "$TEMPLATE" > "$TARGET"
echo "Created: frontend/src/skins/${ID}.css"

insert_before_marker() {
  local file="$1"
  local marker="$2"
  local block="$3"
  local tmp
  tmp="$(mktemp)"
  printf '%s' "$block" > "$tmp"
  python3 - "$file" "$marker" "$tmp" <<'PY'
import sys

path, marker, block_path = sys.argv[1], sys.argv[2], sys.argv[3]
with open(block_path, encoding="utf-8") as f:
    block = f.read()
if block and not block.endswith("\n"):
    block += "\n"

with open(path, encoding="utf-8") as f:
    lines = f.readlines()

out: list[str] = []
for line in lines:
    if marker in line:
        out.append(block)
    out.append(line)

with open(path, "w", encoding="utf-8") as f:
    f.writelines(out)
PY
  rm -f "$tmp"
}

register_go() {
  if grep -q "ID:.*\"${ID}\"" "$GO_FILE"; then
    echo "Go catalog: ${ID} already registered"
    return
  fi
  insert_before_marker "$GO_FILE" "// catalog-end" "$(cat <<EOF
	{
		ID:      "${ID}",
		Label:   "${label}",
		Ready:   false,
		Aliases: []string{},
	},
EOF
)"
  echo "Registered in internal/ui/skins.go"
}

register_ts() {
  if grep -q "id: '${ID}'" "$TS_FILE"; then
    echo "TS manifest: ${ID} already registered"
    return
  fi
  insert_before_marker "$TS_FILE" "// manifest-end" "$(cat <<EOF
  {
    id: '${ID}',
    label: '${label}',
    ready: false,
    aliases: [],
  },
EOF
)"
  echo "Registered in frontend/src/skins/manifest.ts"
}

register_css_import() {
  if grep -q "@import \"./${ID}.css\"" "$CSS_INDEX"; then
    echo "CSS index: ${ID} already imported"
    return
  fi
  insert_before_marker "$CSS_INDEX" "/* imports-begin" "@import \"./${ID}.css\";
"
  echo "Added @import in frontend/src/skins/index.css"
}

if [[ "$REGISTER" == "1" ]]; then
  register_go
  register_ts
  register_css_import
  echo ""
  bash "$ROOT/scripts/validate-skins.sh"
else
  echo ""
  echo "Next steps:"
  echo "  1. Edit frontend/src/skins/${ID}.css — adjust CSS variables"
  echo "  2. Register everywhere: make new-skin ID=${ID} REGISTER=1  (or edit catalog/manifest/index.css manually)"
  echo "  3. config.toml → [ui] skin = \"${ID}\""
  echo "  4. make validate-skins && make frontend-dev"
  echo ""
  echo "Full guide: docs/skins.md"
fi
