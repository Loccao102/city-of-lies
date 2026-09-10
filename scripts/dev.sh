#!/usr/bin/env bash
set -e

echo "Starting City of Lies local development stack..."
if command -v docker &> /dev/null; then
  docker compose -f docker-compose.dev.yml up -d
fi

echo "Start backend: cd backend && go run ./cmd/api"
echo "Start frontend: cd frontend && npm run dev"
