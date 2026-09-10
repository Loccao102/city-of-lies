#!/usr/bin/env bash
set -e

echo "Running Go Unit & Domain Tests..."
cd backend
go test -v ./...
echo "All tests passed!"
