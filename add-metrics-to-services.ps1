# Script to add Prometheus metrics to all Go services

$services = @(
    "user-service",
    "savings-service",
    "budget-service",
    "goal-service",
    "education-service",
    "notification-service",
    "analytics-service"
)

Write-Host "Adding Prometheus dependencies to all services..." -ForegroundColor Green

foreach ($service in $services) {
    Write-Host "`nProcessing $service..." -ForegroundColor Cyan
    
    # Add Prometheus dependencies
    Push-Location $service
    go get github.com/prometheus/client_golang/prometheus github.com/prometheus/client_golang/prometheus/promauto github.com/prometheus/client_golang/prometheus/promhttp
    Pop-Location
    
    Write-Host "Done with $service" -ForegroundColor Green
}

Write-Host "`nAll services updated with Prometheus dependencies!" -ForegroundColor Green
