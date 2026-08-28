# 🖱️ Open Mouse Wizard - 开源鼠标精灵

一个开源的鼠标精灵工具，可以定时随机移动鼠标，防止电脑休眠或屏幕保护程序启动。

## ✨ 功能特性

- ⏱️ 可自定义移动间隔时间（分钟）
- 🖱️ 定时随机移动鼠标，防止休眠
- 🎨 简洁美观的图形界面
- 🖥️ 跨平台支持（macOS / Windows / Linux）

## 🚀 快速开始

### 前置要求

- [Go](https://go.dev/) 1.24 或更高版本
- macOS 用户需要在 **系统设置 > 隐私与安全性 > 辅助功能** 中授权终端访问权限

### 运行方式

#### macOS / Linux

```bash
# 方式一：使用启动脚本
./start.sh

# 方式二：直接运行
go mod tidy
go run main.go
```

#### Windows

```bash
# 使用启动脚本
build.bat

# 或直接运行
go mod tidy
go run main.go
```

## 🏗️ 构建

```bash
# macOS (Apple Silicon)
GOOS=darwin GOARCH=arm64 go build -o mouse-wizard main.go

# macOS (Intel)
GOOS=darwin GOARCH=amd64 go build -o mouse-wizard main.go

# Windows
GOOS=windows GOARCH=amd64 go build -o mouse-wizard.exe main.go
```

## 📖 使用说明

1. 启动程序后会显示图形界面
2. 在输入框中设置移动间隔（单位：分钟）
3. 点击 **"▶️ 启动"** 按钮开始
4. 程序会在设定的间隔后随机移动鼠标
5. 点击 **"⏹️ 停止"** 按钮或关闭窗口即可停止

## ⚠️ 注意事项

- macOS 用户首次使用需要授权辅助功能权限
- 程序会通过系统命令移动鼠标位置（相对偏移）
- 请合理设置间隔时间，建议 1-10 分钟

## 📁 项目结构

```
├── main.go          # 主程序入口
├── config.yaml      # 配置文件
├── bundled.go       # 图标资源打包
├── start.sh         # macOS/Linux 启动脚本
├── build.bat        # Windows 启动脚本
├── icon.png         # 应用图标
└── README.md        # 项目文档
```

## 📜 许可证

本项目基于 MIT 许可证开源。