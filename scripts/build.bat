@echo off
chcp 65001 >nul
echo ========================================
echo   鼠标精灵 - 构建脚本 (Windows)
echo   Mouse Wizard - Build Script
echo ========================================
echo.

if not exist bin mkdir bin

echo [1/3] 检查依赖...
where go >nul 2>&1
if %errorlevel% neq 0 (
    echo ❌ 未检测到 Go，请先安装 Go
    echo    下载地址: https://go.dev/dl/
    pause
    exit /b 1
)

go version
echo.

echo [2/3] 构建 Windows 版本 (Fyne GUI)...
echo    使用 CGO_ENABLED=1 编译 Fyne GUI 应用
echo    添加 -ldflags="-H windowsgui" 隐藏命令行窗口
echo.

set CGO_ENABLED=1
go mod tidy
go build -ldflags="-H windowsgui" -o bin\mouse-wizard.exe ./cmd/app
if %errorlevel% neq 0 (
    echo.
    echo ❌ 构建失败！
    echo.
    echo 可能的原因：
    echo   1. 未安装 MinGW-w64 (CGO 需要)
    echo      下载地址: https://sourceforge.net/projects/mingw-w64/
    echo   2. 未安装 OpenGL 驱动
    echo   3. Go 版本不兼容
    echo.
    echo 安装 MinGW-w64 后，请将其 bin 目录添加到系统 PATH 环境变量
    pause
    exit /b 1
)
echo.
echo ✅ 构建成功
echo.

echo [3/3] 验证文件...
if exist bin\mouse-wizard.exe (
    for %%I in (bin\mouse-wizard.exe) do set size=%%~zI
    echo bin\mouse-wizard.exe 已生成 (大小: %size% 字节)
) else (
    echo ❌ bin\mouse-wizard.exe 未找到！
    pause
    exit /b 1
)
echo.

echo ========================================
echo   ✅ 构建完成！
echo   输出: bin\mouse-wizard.exe
echo.
echo   💡 双击 bin\mouse-wizard.exe 即可运行
echo   💡 如需修改配置，编辑 configs\config.yaml
echo ========================================
echo.
pause