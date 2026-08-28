#!/bin/bash

echo "🖱️ 鼠标精灵启动脚本"
echo "📌 功能：定时随机移动鼠标，防止电脑休眠"
echo ""

if [ -d "bin/MouseWizard.app" ]; then
    echo "🚀 正在启动鼠标精灵 (MouseWizard.app)..."
    echo "💡 使用说明：在界面中设置移动间隔（秒），点击启动即可"
    echo "⚠️  按 Cmd+Q 或关闭窗口停止程序"
    echo ""
    open bin/MouseWizard.app
elif [ -f "bin/mouse-wizard-macos" ]; then
    echo "🚀 正在启动鼠标精灵 (bin/mouse-wizard-macos)..."
    echo "💡 使用说明：在界面中设置移动间隔（秒），点击启动即可"
    echo "⚠️  按 Cmd+Q 或关闭窗口停止程序"
    echo ""
    ./bin/mouse-wizard-macos
else
    echo "📦 正在整理依赖..."
    go mod tidy
    if [ $? -ne 0 ]; then
        echo "❌ 依赖整理失败，请检查 go.mod 文件"
        exit 1
    fi
    echo "✅ 依赖整理完成"
    echo ""
    echo "🚀 正在启动鼠标精灵 (开发模式)..."
    echo "💡 使用说明：在界面中设置移动间隔（秒），点击启动即可"
    echo "⚠️  按 Cmd+Q 或关闭窗口停止程序"
    echo ""
    go run ./cmd/app
fi