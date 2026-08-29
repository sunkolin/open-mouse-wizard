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

$fyneCmd = Get-Command fyne -ErrorAction SilentlyContinue
if (-not $fyneCmd) {
    Write-Host "未检测到 fyne CLI，正在安装..." -ForegroundColor Yellow
    go install fyne.io/fyne/v2/cmd/fyne@latest
    if ($LASTEXITCODE -ne 0) {
        Write-Host "fyne CLI 安装失败" -ForegroundColor Red
        Read-Host "按回车键退出"
        exit 1
    }
}

Write-Host "[2/3] 构建 Windows 版本 (Fyne GUI)..." -ForegroundColor Yellow
Write-Host "    使用 fyne package 构建（自动嵌入图标）" -ForegroundColor DarkGray
Write-Host ""

Copy-Item "icon.png" "cmd\app\icon.png" -Force -ErrorAction SilentlyContinue

Push-Location "cmd\app"
fyne package --icon icon.png --os windows
$buildExitCode = $LASTEXITCODE
Pop-Location

Remove-Item "cmd\app\icon.png" -ErrorAction SilentlyContinue

if ($buildExitCode -ne 0) {
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

$version = ""
$fyneConfig = Get-Content "FyneApp.toml" -ErrorAction SilentlyContinue
foreach ($line in $fyneConfig) {
    if ($line -match '^\s*Version\s*=\s*"([^"]+)"') {
        $version = $Matches[1]
        break
    }
}
if (-not $version) { $version = "0.0.0" }

$exeName = "开源鼠标精灵-$version.exe"
$exePath = "bin\$exeName"

Copy-Item "cmd\app\app.exe" $exePath -Force
Remove-Item "cmd\app\app.exe" -ErrorAction SilentlyContinue

if (Test-Path "manifest.xml") {
    $mtCmd = Get-Command mt -ErrorAction SilentlyContinue
    if ($mtCmd) {
        Write-Host "    注入自定义 manifest (管理员权限)..." -ForegroundColor DarkGray
        mt -manifest "manifest.xml" -outputresource:"$exePath;#1" 2>$null
    }
}

Write-Host ""
Write-Host "构建成功" -ForegroundColor Green
Write-Host ""

Write-Host "[3/3] 验证文件..." -ForegroundColor Yellow
if (-not (Test-Path $exePath)) {
    Write-Host "$exeName 未找到！" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

$size = (Get-Item $exePath).Length
Write-Host "$exeName 已生成 (大小: $size 字节)" -ForegroundColor Green

Write-Host ""
Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  构建完成！" -ForegroundColor Green
Write-Host "  输出: $exePath" -ForegroundColor White
Write-Host ""
Write-Host "  双击 $exeName 即可运行" -ForegroundColor White
Write-Host "  如需修改配置，编辑 configs\config.yaml" -ForegroundColor White
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""