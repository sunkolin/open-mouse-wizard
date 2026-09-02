[Console]::OutputEncoding = [System.Text.Encoding]::UTF8

Write-Host "========================================" -ForegroundColor Cyan
Write-Host "  鼠标精灵 - 构建脚本 (Windows)" -ForegroundColor Cyan
Write-Host "  Open Mouse Wizard - Build Script" -ForegroundColor Cyan
Write-Host "========================================" -ForegroundColor Cyan
Write-Host ""

$PROJECT_DIR = Split-Path -Parent $PSScriptRoot
Set-Location $PROJECT_DIR

if (-not (Test-Path "node_modules")) {
    Write-Host "[1/4] 安装依赖..." -ForegroundColor Yellow
    npm install
    if ($LASTEXITCODE -ne 0) {
        Write-Host "依赖安装失败" -ForegroundColor Red
        Read-Host "按回车键退出"
        exit 1
    }
} else {
    Write-Host "[1/4] 依赖已安装，跳过..." -ForegroundColor Yellow
}

Write-Host ""
Write-Host "[2/4] 检查构建环境..." -ForegroundColor Yellow

$nodeCmd = Get-Command node -ErrorAction SilentlyContinue
if (-not $nodeCmd) {
    Write-Host "未检测到 Node.js，请先安装 Node.js" -ForegroundColor Red
    Write-Host "    下载地址: https://nodejs.org/" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

node --version
npm --version

Write-Host ""
Write-Host "[3/4] 读取版本号..." -ForegroundColor Yellow

$VERSION = ""
try {
    $pkg = Get-Content "package.json" -Raw | ConvertFrom-Json
    $VERSION = $pkg.version
} catch {}
if (-not $VERSION) { $VERSION = "0.0.0" }

Write-Host "    版本: $VERSION" -ForegroundColor Green

Write-Host ""
Write-Host "[4/4] 构建 Windows 版本..." -ForegroundColor Yellow

$PLATFORM = if ($args.Count -gt 0) { $args[0] } else { "win" }

switch ($PLATFORM) {
    "mac" {
        Write-Host "    构建 macOS 版本..." -ForegroundColor DarkGray
        npm run build:mac
    }
    "win" {
        Write-Host "    构建 Windows 版本..." -ForegroundColor DarkGray
        npm run build:win
    }
    "all" {
        Write-Host "    构建所有平台版本..." -ForegroundColor DarkGray
        npm run build:all
    }
    default {
        Write-Host "    未知平台: $PLATFORM" -ForegroundColor Red
        Write-Host "    用法: .\build.ps1 [mac|win|all]" -ForegroundColor Red
        Read-Host "按回车键退出"
        exit 1
    }
}

if ($LASTEXITCODE -ne 0) {
    Write-Host ""
    Write-Host "构建失败！" -ForegroundColor Red
    Read-Host "按回车键退出"
    exit 1
}

Write-Host ""
Write-Host "========================================" -ForegroundColor Green
Write-Host "  构建完成！" -ForegroundColor Green
Write-Host "  版本: $VERSION" -ForegroundColor Green
Write-Host "  输出目录: dist\" -ForegroundColor Green
Write-Host "========================================" -ForegroundColor Green
Write-Host ""

if (Test-Path "dist") {
    Get-ChildItem "dist" -Recurse -File | ForEach-Object {
        $sizeMB = [math]::Round($_.Length / 1MB, 2)
        Write-Host "  $($_.Name) ($sizeMB MB)" -ForegroundColor Green
    }
}