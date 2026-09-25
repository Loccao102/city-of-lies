#!/usr/bin/env bash
# Bash Soak Test Runner (Phase 11)
set -euo pipefail

RUNS="${1:-100}"
WORKERS="${2:-10}"
SCENARIO="${3:-}"

echo "Starting City of Lies Automated Soak Test ($RUNS seeds, $WORKERS workers)..."

cd "$(dirname "$0")/../backend"
if [ -n "$SCENARIO" ]; then
    go run ./cmd/soak-test -runs "$RUNS" -workers "$WORKERS" -scenario "$SCENARIO"
else
    go run ./cmd/soak-test -runs "$RUNS" -workers "$WORKERS"
fi
echo "Soak test completed successfully!"
