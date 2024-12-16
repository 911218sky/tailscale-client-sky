param(
    [switch]$UseUpx = $false
)

Write-Host "🚀 Starting build process..." -ForegroundColor Green
Write-Host "UPX Compression: $(if ($UseUpx) { 'Enabled' } else { 'Disabled' })" -ForegroundColor Yellow

# Check if go-winres is installed
$goWinres = Get-Command go-winres -ErrorAction SilentlyContinue
if (-not $goWinres) {
    Write-Host "Installing go-winres..." -ForegroundColor Yellow
    go install github.com/tc-hib/go-winres@latest
}

# Add icon to Windows executable
Write-Host "Adding icon to Windows executable..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o sky-tailscale.exe main.go
go-winres simply --icon ./img/sky-tailscale-icon.png

# Build Windows version
Write-Host "Building Windows executable..." -ForegroundColor Cyan
go build -ldflags "-s -w" -o sky-tailscale.exe

# Build Linux version
Write-Host "Building Linux executable..." -ForegroundColor Cyan
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -ldflags "-s -w" -o sky-tailscale-linux main.go

# Reset GOOS and GOARCH
$env:GOOS = "windows"
$env:GOARCH = "amd64"

# Check if UPX compression is requested
if ($UseUpx) {
    $upx = Get-Command upx -ErrorAction SilentlyContinue
    if ($upx) {
        Write-Host "Compressing executables with UPX..." -ForegroundColor Cyan
        upx -9 sky-tailscale.exe
        upx -9 sky-tailscale-linux
    } else {
        Write-Host "⚠️ UPX not found. Skipping compression." -ForegroundColor Yellow
        Write-Host "To enable compression, please install UPX and add it to your PATH." -ForegroundColor Yellow
    }
} else {
    Write-Host "Skipping UPX compression (not requested)." -ForegroundColor Yellow
}

Write-Host "✅ Build process completed!" -ForegroundColor Green
Write-Host "Generated files:" -ForegroundColor Cyan
Write-Host "- sky-tailscale.exe (Windows)" -ForegroundColor White
Write-Host "- sky-tailscale-linux (Linux)" -ForegroundColor White

# Clean up resource files
Remove-Item -Path "rsrc_windows_amd64.syso" -ErrorAction SilentlyContinue
Remove-Item -Path "rsrc_windows_386.syso" -ErrorAction SilentlyContinue