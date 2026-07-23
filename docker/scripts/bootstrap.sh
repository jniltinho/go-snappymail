#!/bin/bash
# Bootstrap lab stack: domains, mailboxes, SnappyMail presets.
set -euo pipefail

cd "$(dirname "$0")/.."

echo "==> Waiting for Go-PostfixAdmin..."
for i in $(seq 1 60); do
  if docker compose exec -T postfixadmin wget -q -O /dev/null http://localhost:8080/ 2>/dev/null; then
    break
  fi
  sleep 2
done

bash "$(dirname "$0")/seed-lab.sh"

echo "==> Bootstrap complete"
