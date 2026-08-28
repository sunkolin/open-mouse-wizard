#!/bin/bash

echo "🖱️ 鼠标精灵构建脚本"
echo ""

mkdir -p bin

OS="$(uname -s)"
ARCH="$(uname -m)"

APP_NAME="mouse-wizard"
BUNDLE_ID="com.mousewizard.app"

build_macos() {
    HELPER_SRC="internal/service/helper/mouse_helper.c"
    HELPER_BIN="internal/service/helper/mouse_helper_darwin"

    if [ -f "$HELPER_SRC" ]; then
        echo "[1/4] 编译 macOS 鼠标辅助程序..."
        ARCH_FLAG=""
        if [ "$ARCH" = "arm64" ]; then
            ARCH_FLAG="-arch arm64"
        elif [ "$ARCH" = "x86_64" ]; then
            ARCH_FLAG="-arch x86_64"
        fi
        clang $ARCH_FLAG -o "$HELPER_BIN" "$HELPER_SRC" -framework CoreGraphics -framework CoreFoundation 2>/dev/null
        if [ $? -ne 0 ]; then
            echo "  ⚠️ 辅助程序编译失败，尝试通用编译..."
            clang -o "$HELPER_BIN" "$HELPER_SRC" -framework CoreGraphics -framework CoreFoundation 2>/dev/null
        fi
        if [ -f "$HELPER_BIN" ]; then
            echo "  ✅ 辅助程序编译成功"
        else
            echo "  ⚠️ 辅助程序编译失败，鼠标移动功能可能不可用"
        fi
    fi

    echo "[2/4] 构建 macOS 版本..."
    CGO_ENABLED=1 go build -o bin/mouse-wizard-macos ./cmd/app
    if [ $? -ne 0 ]; then
        echo "❌ 构建失败"
        return 1
    fi
    echo "✅ 二进制构建成功"

    echo "[3/4] 打包 macOS .app 应用..."
    APP_DIR="bin/${APP_NAME}.app"
    rm -rf "$APP_DIR"

    mkdir -p "$APP_DIR/Contents/MacOS"
    mkdir -p "$APP_DIR/Contents/Resources"

    cp bin/mouse-wizard-macos "$APP_DIR/Contents/MacOS/mouse-wizard"
    chmod +x "$APP_DIR/Contents/MacOS/mouse-wizard"

    if [ -f "configs/config.yaml" ]; then
        cp configs/config.yaml "$APP_DIR/Contents/Resources/config.yaml"
    fi

    ICON_PNG="internal/api/icon.png"
    if [ -f "$ICON_PNG" ]; then
        echo "  转换图标 PNG → ICNS..."
        ICONSET_DIR="/tmp/mousewizard.iconset"
        rm -rf "$ICONSET_DIR"
        mkdir -p "$ICONSET_DIR"
        sips -z 16 16     "$ICON_PNG" --out "$ICONSET_DIR/icon_16x16.png" 2>/dev/null
        sips -z 32 32     "$ICON_PNG" --out "$ICONSET_DIR/icon_16x16@2x.png" 2>/dev/null
        sips -z 32 32     "$ICON_PNG" --out "$ICONSET_DIR/icon_32x32.png" 2>/dev/null
        sips -z 64 64     "$ICON_PNG" --out "$ICONSET_DIR/icon_32x32@2x.png" 2>/dev/null
        sips -z 128 128   "$ICON_PNG" --out "$ICONSET_DIR/icon_128x128.png" 2>/dev/null
        sips -z 256 256   "$ICON_PNG" --out "$ICONSET_DIR/icon_128x128@2x.png" 2>/dev/null
        sips -z 256 256   "$ICON_PNG" --out "$ICONSET_DIR/icon_256x256.png" 2>/dev/null
        sips -z 512 512   "$ICON_PNG" --out "$ICONSET_DIR/icon_256x256@2x.png" 2>/dev/null
        sips -z 512 512   "$ICON_PNG" --out "$ICONSET_DIR/icon_512x512.png" 2>/dev/null
        sips -z 1024 1024 "$ICON_PNG" --out "$ICONSET_DIR/icon_512x512@2x.png" 2>/dev/null

        if iconutil -c icns "$ICONSET_DIR" -o "$APP_DIR/Contents/Resources/icon.icns" 2>/dev/null; then
            echo "  ✅ icon.icns 生成成功"
        else
            echo "  ⚠️ iconutil 转换失败"
        fi
        rm -rf "$ICONSET_DIR"
    fi

    cat > "$APP_DIR/Contents/Info.plist" << 'PLIST'
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
    <key>CFBundleDevelopmentRegion</key>
    <string>zh_CN</string>
    <key>CFBundleExecutable</key>
    <string>mouse-wizard</string>
    <key>CFBundleIconFile</key>
    <string>icon.icns</string>
    <key>CFBundleIdentifier</key>
    <string>com.mousewizard.app</string>
    <key>CFBundleName</key>
    <string>MouseWizard</string>
    <key>CFBundlePackageType</key>
    <string>APPL</string>
    <key>CFBundleShortVersionString</key>
    <string>1.0.0</string>
    <key>CFBundleVersion</key>
    <string>1</string>
    <key>LSMinimumSystemVersion</key>
    <string>10.13</string>
    <key>NSHighResolutionCapable</key>
    <true/>
    <key>NSPrincipalClass</key>
    <string>NSApplication</string>
</dict>
</plist>
PLIST

    echo "  ✅ .app 包构建成功"
    echo "[4/4] macOS 构建完成"
}

build_windows() {
    echo "[1/1] 交叉编译 Windows amd64 版本..."
    echo "💡 使用 Gio 软件渲染，无需 CGO，可直接在 macOS 上交叉编译"

    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build \
        -ldflags="-H windowsgui" \
        -o bin/mouse-wizard.exe \
        ./cmd/app

    if [ $? -eq 0 ]; then
        echo "✅ Windows amd64 版本构建成功"
        return 0
    fi

    echo "❌ Windows 版本构建失败"
    echo "💡 如果需要硬件加速渲染，可安装 mingw-w64："
    echo "   brew install mingw-w64"
    echo "   然后使用 CGO_ENABLED=1 重新构建"
    return 1
}

if [ "$OS" = "Darwin" ]; then
    echo "当前系统: macOS ($ARCH)"
    echo ""
    build_macos
    echo ""
    build_windows
else
    echo "❌ 当前系统为 $OS，本脚本仅支持 macOS"
    exit 1
fi

echo ""
echo "📁 bin/ 目录内容:"
ls -lh bin/ 2>/dev/null || echo "   (空)"
echo ""
echo "✅ 构建完成"
echo ""
echo "💡 macOS 用户: 双击 bin/MouseWizard.app 即可运行"
echo "💡 Windows 用户: 将 bin/mouse-wizard.exe 复制到 Windows 上运行"
echo "💡 如需修改配置，编辑 configs/config.yaml"