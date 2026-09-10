# Windows PowerShell Test Script
Write-Host "Running City of Lies test suite..." -ForegroundColor Cyan

Push-Location backend
Write-Host "Running Go Unit & Domain Tests..." -ForegroundColor Green
go test -v ./...
if ($LASTEXITCODE -ne 0) {
    Write-Host "Go tests failed!" -ForegroundColor Red
    Pop-Location
    exit 1
}
Pop-Location

Write-Host "All tests passed successfully!" -ForegroundColor Green
