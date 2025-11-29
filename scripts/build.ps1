# Windows 本地构建脚本 - 将前端嵌入到后端二进制中

Write-Host "=== SeekYourGuild Build Script ===" -ForegroundColor Cyan

# 构建前端
Write-Host ""
Write-Host ">>> Building frontend..." -ForegroundColor Yellow
Push-Location frontend
yarn install
yarn build
Pop-Location

# 复制前端到后端静态目录
Write-Host ""
Write-Host ">>> Copying frontend to backend..." -ForegroundColor Yellow
if (Test-Path "backend/internal/static/dist") {
    Remove-Item -Recurse -Force "backend/internal/static/dist"
}
Copy-Item -Recurse "frontend/dist" "backend/internal/static/dist"

# 构建后端
Write-Host ""
Write-Host ">>> Building backend..." -ForegroundColor Yellow
Push-Location backend
go build -ldflags="-s -w" -o ../seekyourguild.exe ./cmd/server/
Pop-Location

Write-Host ""
Write-Host "=== Build complete! ===" -ForegroundColor Green
Write-Host "Binary: ./seekyourguild.exe"
Write-Host ""
Write-Host "Run with:"
Write-Host "  .\seekyourguild.exe"
Write-Host ""
Write-Host "Or with custom settings:"
Write-Host '  $env:API_KEY="your-key"; $env:ALLOWED_GROUPS="123,456"; .\seekyourguild.exe'
