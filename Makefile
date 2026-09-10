.PHONY: all dev backend frontend test test-race lint db-up db-down migrate-up sqlc scenario-check sim clean

all: backend frontend

dev:
	docker compose -f docker-compose.dev.yml up --build

backend:
	cd backend && go build -o bin/api.exe ./cmd/api

sim:
	cd backend && go build -o bin/sim.exe ./cmd/sim && ./bin/sim.exe

scenario-check:
	cd backend && go run ./cmd/scenario-check ../data/scenarios/riverside-factory

frontend:
	cd frontend && npm run build

test:
	cd backend && go test -v ./...

test-race:
	cd backend && go test -race -v ./...

lint:
	cd backend && go vet ./...

db-up:
	docker compose up -d postgres

db-down:
	docker compose down

migrate-up:
	@echo "Running migrations..."

sqlc:
	cd backend && sqlc generate

clean:
	rm -rf backend/bin frontend/.next frontend/out
