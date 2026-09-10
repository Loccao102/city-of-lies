# Windows PowerShell Soak Test Runner (Phase 11)
param (
    [int]$Runs = 100,
    [int]$Workers = 10
)

Write-Host "Starting City of Lies Automated Soak Test ($Runs seeds, $Workers workers)..." -ForegroundColor Cyan

Push-Location backend
go run ./cmd/soak-test -runs $Runs -workers $Workers
$code = $LASTEXITCODE
Pop-Location

if ($code -ne 0) {
    Write-Host "Soak test failed!" -ForegroundColor Red
    exit $code
}

Write-Host "Soak test completed successfully!" -ForegroundColor Green
