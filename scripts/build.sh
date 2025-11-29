#!/bin/bash
# 本地构建脚本 - 将前端嵌入到后端二进制中

set -e

echo "=== SeekYourGuild Build Script ==="

# 构建前端
echo ""
echo ">>> Building frontend..."
cd frontend
yarn install
yarn build
cd ..

# 复制前端到后端静态目录
echo ""
echo ">>> Copying frontend to backend..."
rm -rf backend/internal/static/dist
cp -r frontend/dist backend/internal/static/dist

# 构建后端
echo ""
echo ">>> Building backend..."
cd backend
go build -ldflags="-s -w" -o ../seekyourguild ./cmd/server/
cd ..

echo ""
echo "=== Build complete! ==="
echo "Binary: ./seekyourguild"
echo ""
echo "Run with:"
echo "  ./seekyourguild"
echo ""
echo "Or with custom settings:"
echo "  API_KEY=your-key ALLOWED_GROUPS=123,456 ./seekyourguild"
