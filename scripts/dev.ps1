# Windows PowerShell Dev Script
Write-Host "Starting City of Lies local development stack..." -ForegroundColor Cyan

if (Get-Command docker -ErrorAction SilentlyContinue) {
    Write-Host "Starting PostgreSQL container..." -ForegroundColor Green
    docker compose -f docker-compose.dev.yml up -d
}

Write-Host "To start backend: cd backend; go run ./cmd/api" -ForegroundColor Yellow
Write-Host "To start frontend: cd frontend; npm run dev" -ForegroundColor Yellow
