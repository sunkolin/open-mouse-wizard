[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  鼠标精灵 - 构建脚本 (Windows)" -ForegroundColor Cyan
Write-Host "  Mouse Wizard - Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

if (-not (Test-Path "bin")) {
    New-Item -ItemType Directory -Path "bin" | Out-Null
}

Write-Host "[1/3] 检查依赖..." -ForegroundColor Yellow
$goCmd = Get-Command go -ErrorAction SilentlyContinue
if (-not $goCmd) {
    Write-Host "未检测到 Go，请先安装 Go" -ForegroundColor Red
    Write-Host "    下载地址: https://go.dev/dl/" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

go version
Write-Host ""

Write-Host "[2/3] 构建 Windows 版本 (Fyne GUI)..." -ForegroundColor Yellow
Write-Host "    使用 CGO_ENABLED=1 编译 Fyne GUI 应用" -ForegroundColor DarkGray
Write-Host '    添加 -ldflags="-H windowsgui" 隐藏命令行窗口' -ForegroundColor DarkGray
Write-Host ""

$env:CGO_ENABLED = "1"
go mod tidy
go build -ldflags="-H windowsgui" -o "bin\mouse-wizard.exe" ./cmd/app

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "构建失败！" -ForegroundColor Red
    Write-Host ""
    Write-Host "可能的原因：" -ForegroundColor Red
    Write-Host "  1. 未安装 MinGW-w64 (CGO 需要)" -ForegroundColor Red
    Write-Host "     下载地址: https://sourceforge.net/projects/mingw-w64/" -ForegroundColor Red
    Write-Host "  2. 未安装 OpenGL 驱动" -ForegroundColor Red
    Write-Host "  3. Go 版本不兼容" -ForegroundColor Red
    Write-Host ""
    Write-Host "安装 MinGW-w64 后，请将其 bin 目录添加到系统 PATH 环境变量" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

Write-Host ""
Write-Host "构建成功" -ForegroundColor Green
Write-Host ""

Write-Host "[3/3] 验证文件..." -ForegroundColor Yellow
$exePath = "bin\mouse-wizard.exe"
if (-not (Test-Path $exePath)) {
    Write-Host "bin\mouse-wizard.exe 未找到！" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

$size = (Get-Item $exePath).Length
Write-Host "bin\mouse-wizard.exe 已生成 (大小: $size 字节)" -ForegroundColor Green

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  构建完成！" -ForegroundColor Green
Write-Host "  输出: bin\mouse-wizard.exe" -ForegroundColor White
Write-Host ""
Write-Host "  双击 bin\mouse-wizard.exe 即可运行" -ForegroundColor White
Write-Host "  如需修改配置，编辑 configs\config.yaml" -ForegroundColor White
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""