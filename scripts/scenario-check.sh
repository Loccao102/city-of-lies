#!/usr/bin/env bash
set -e

echo "Validating Riverside Factory Scenario JSON bundle..."
cd backend
go run ./cmd/scenario-check ../data/scenarios/riverside-factory
