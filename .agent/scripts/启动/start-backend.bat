@echo off
chcp 65001 >nul
echo ====================================
echo SeekYourGuild - 后端启动脚本
echo ====================================

cd /d %~dp0backend

echo [1/2] 安装依赖...
go mod tidy

echo [2/2] 启动服务器...
go run ./cmd/server/main.go

pause
