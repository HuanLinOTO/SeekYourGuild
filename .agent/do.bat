@echo off
chcp 65001 >nul
cd /d d:\Projects\SeekYourGuild\backend
go run ./cmd/server/main.go
