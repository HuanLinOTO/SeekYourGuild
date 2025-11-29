@echo off
chcp 65001 >nul
echo ====================================
echo SeekYourGuild - 前端启动脚本
echo ====================================

cd /d %~dp0..\..\..\frontend

echo [1/2] 安装依赖...
call npm install

echo [2/2] 启动开发服务器...
call npm run dev

pause
