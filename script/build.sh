#!/bin/bash

set -e

echo "========================================"
echo "  鼠标精灵 - 构建脚本 (Electron)"
echo "  Open Mouse Wizard - Build Script"
echo "========================================"
echo ""

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_DIR"

export ELECTRON_MIRROR="${ELECTRON_MIRROR:-https://npmmirror.com/mirrors/electron/}"

if [ ! -d "node_modules" ]; then
  echo "[1/4] 安装依赖..."
  npm install
else
  echo "[1/4] 依赖已安装，跳过..."
fi

echo ""
echo "[2/4] 检查构建环境..."

if ! command -v node &> /dev/null; then
  echo "未检测到 Node.js，请先安装 Node.js" >&2
  echo "    下载地址: https://nodejs.org/" >&2
  exit 1
fi

if ! command -v npm &> /dev/null; then
  echo "未检测到 npm，请先安装 npm" >&2
  exit 1
fi

node --version
npm --version

echo ""
echo "[3/4] 读取版本号..."

VERSION=$(node -e "
  try {
    const pkg = require('./package.json');
    console.log(pkg.version || '0.0.0');
  } catch(e) {
    console.log('0.0.0');
  }
")

echo "    版本: $VERSION"

echo ""
echo "[4/4] 开始构建..."

PLATFORM="${1:-all}"

case "$PLATFORM" in
  mac|macos|darwin)
    echo "    构建 macOS 版本..."
    npm run build:mac
    ;;
  win|windows|win32)
    echo "    构建 Windows 版本..."
    npm run build:win
    ;;
  all)
    echo "    构建所有平台版本 (macOS + Windows)..."
    npm run build:all
    ;;
  *)
    echo "    未知平台: $PLATFORM" >&2
    echo "    用法: $0 [mac|win|all]" >&2
    exit 1
    ;;
esac

echo ""
echo "========================================"
echo "  构建完成！" 
echo "  版本: $VERSION"
echo "  输出目录: dist/"
echo "========================================"
echo ""

ls -la dist/ 2>/dev/null || echo "请检查 dist/ 目录获取构建产物"