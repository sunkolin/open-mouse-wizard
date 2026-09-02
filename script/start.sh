#!/bin/bash

set -e

SCRIPT_DIR="$(cd "$(dirname "$0")" && pwd)"
PROJECT_DIR="$(cd "$SCRIPT_DIR/.." && pwd)"

cd "$PROJECT_DIR"

export ELECTRON_MIRROR="${ELECTRON_MIRROR:-https://npmmirror.com/mirrors/electron/}"

if [ ! -d "node_modules" ]; then
  echo "安装依赖..."
  npm install
fi

if [ ! -f "node_modules/electron/path.txt" ]; then
  echo "electron 未正确安装（缺少 path.txt），重新执行安装脚本..."
  node node_modules/electron/install.js
fi

npx electron .