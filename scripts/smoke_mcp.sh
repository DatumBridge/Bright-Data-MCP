#!/usr/bin/env bash
set -euo pipefail
BASE="${1:-http://localhost:8011}"

echo "== health =="
curl -sf "$BASE/health" | tee /dev/stderr
echo

echo "== initialize =="
INIT=$(curl -sf -D - -o /tmp/bd-mcp-init.json -X POST "$BASE/mcp" \
  -H 'Content-Type: application/json' \
  -d '{"jsonrpc":"2.0","id":1,"method":"initialize","params":{}}')
SID=$(echo "$INIT" | awk -F': ' 'tolower($1)=="mcp-session-id"{print $2}' | tr -d '\r')
echo "session=$SID"
cat /tmp/bd-mcp-init.json
echo

echo "== tools/list =="
curl -sf -X POST "$BASE/mcp" \
  -H 'Content-Type: application/json' \
  -H "Mcp-Session-Id: $SID" \
  -d '{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}'
echo
