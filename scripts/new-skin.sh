#!/usr/bin/env bash
# Scaffold a new webmail skin CSS file.
# Usage: scripts/new-skin.sh <skin-id>
# Full guide: docs/skins.md

set -euo pipefail

ROOT="$(cd "$(dirname "$0")/.." && pwd)"
ID="${1:-}"

if [[ -z "$ID" ]]; then
  echo "Usage: make new-skin ID=<skin-id>"
  echo "Example: make new-skin ID=acme"
  exit 1
fi

if [[ ! "$ID" =~ ^[a-z][a-z0-9-]*$ ]]; then
  echo "Error: skin id must be lowercase letters, digits, or hyphens (start with a letter)."
  echo "Got: $ID"
  exit 1
fi

TARGET="$ROOT/frontend/src/skins/${ID}.css"
TEMPLATE="$ROOT/frontend/src/skins/_template.css"

if [[ -f "$TARGET" ]]; then
  echo "Error: $TARGET already exists."
  exit 1
fi

sed "s/MY_SKIN_ID/${ID}/g" "$TEMPLATE" > "$TARGET"
echo "Created: frontend/src/skins/${ID}.css"
echo ""
echo "Next steps (see docs/skins.md):"
echo "  1. Edit frontend/src/skins/${ID}.css — adjust CSS variables"
echo "  2. internal/ui/skins.go — add constant, available[], NormalizeSkin case"
echo "  3. frontend/src/skins/types.ts — add '${ID}' to SkinId"
echo "  4. frontend/src/skins/registry.ts — SKIN_REGISTRY + normalizeSkinId()"
echo "  5. frontend/src/style.css — @import \"./skins/${ID}.css\";"
echo "  6. config.toml — [ui] skin = \"${ID}\""
echo "  7. make frontend-dev && make run"
