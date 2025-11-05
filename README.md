# Magic Input

一个本地剪切板 + AI 查询分析的 APP，可以对任意的输入进行分析，可接入 AI 大模型进行进一步的分析查询

## 关于项目

这是一个使用 Wails v2 + React + TypeScript 构建的桌面应用程序。

项目配置可以通过编辑 `wails.json` 文件进行调整。更多项目设置信息请参考：https://wails.io/docs/reference/project-config

## 开发模式

要在实时开发模式下运行，请在项目目录中运行 `wails dev`。这将启动一个 Vite 开发服务器，提供前端代码的快速热重载。

如果你想在浏览器中开发并访问 Go 方法，还有一个运行在 http://localhost:34115 的开发服务器。在浏览器中连接到此地址，你可以从开发工具调用 Go 代码。

## 构建

### 快速构建

使用 Makefile 命令（推荐）：

```bash
# 构建 Windows exe 文件
make build-windows

# 构建 macOS 应用
make build-darwin

# 构建所有平台
make build-all

# 清理构建文件
make clean

# 启动开发模式
make dev
```

### 手动构建

#### Windows 版本

**在 macOS/Linux 上交叉编译：**

```bash
# 使用构建脚本
./build-windows.sh

# 或者直接使用 wails 命令
wails build -platform windows/amd64 -clean -o magic-input-app.exe
```

**在 Windows 上构建：**

```powershell
# 使用 PowerShell 脚本
.\build-windows.ps1

# 或者使用参数
.\build-windows.ps1 -Clean -Debug

# 或者直接使用 wails 命令
wails build -platform windows/amd64 -clean
```

#### macOS 版本

```bash
wails build -platform darwin/amd64 -clean
```

#### Linux 版本

```bash
wails build -platform linux/amd64 -clean
```

## 输出文件

- **Windows**: `build/bin/magic-input-app.exe`
- **macOS**: `build/bin/magic-input-app.app`
- **Linux**: `build/bin/magic-input-app`

## 依赖安装

```bash
# 安装所有依赖
make install

# 或者手动安装
go mod tidy
cd frontend && npm install
```

## 系统要求

- Go 1.23+
- Node.js 16+
- Wails CLI v2.x

## 安装 Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## 特性

- 🚀 现代化的 React + TypeScript 前端
- ⚡ 高性能的 Go 后端
- 🎨 原生桌面应用体验
- 📦 一键打包多平台应用
- 🔧 热重载开发模式
- 🔄 **自动更新系统**
- 🛡️ 安全的版本管理
- 📊 实时更新进度
- ⚙️ 可配置更新选项

## 🔄 自动更新功能

### 功能特点

- ✅ **自动检查更新**：应用启动时自动检查新版本
- ✅ **手动检查**：用户可随时检查更新
- ✅ **下载进度**：实时显示下载进度和速度
- ✅ **一键更新**：自动下载、安装和重启
- ✅ **版本管理**：支持跳过特定版本
- ✅ **配置管理**：自定义检查间隔和更新渠道
- ✅ **多平台支持**：Windows、macOS、Linux

### 快速开始

1. **配置更新源**（首次使用）：

```bash
# 编辑 updater.go 和 config.go 中的 GitHub 仓库地址
# 将 YOUR_USERNAME/YOUR_REPO 替换为实际仓库
```

2. **发布新版本**：

```bash
# 创建版本标签
git tag v1.0.1
git push origin v1.0.1

# GitHub Actions 会自动构建并发布
```

3. **用户更新体验**：
   - 应用启动自动检查更新
   - 发现新版本时显示更新对话框
   - 一键下载安装，自动重启

### 详细文档

- 📖 [自动更新实现指南](AUTO_UPDATE_GUIDE.md)
- 🚀 [快速开始指南](QUICK_START_UPDATE.md)
- 🔧 [构建优化文档](BUILD_OPTIMIZATION.md)
