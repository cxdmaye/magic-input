#!/bin/bash

# Windows 构建脚本
# 用于在 macOS 上交叉编译 Windows 应用

echo "开始构建 Windows 版本..."

# 检查是否安装了 wails
if ! command -v wails &> /dev/null; then
    echo "错误: 未找到 wails 命令。请先安装 Wails CLI："
    echo "go install github.com/wailsapp/wails/v2/cmd/wails@latest"
    exit 1
fi

# 清理之前的构建
echo "清理构建目录..."
rm -rf build/bin/

# 构建 Windows 版本
echo "构建 Windows exe..."
wails build -platform windows/amd64 -clean

# 检查构建结果
if [ -f "build/bin/magic-input-app.exe" ]; then
    echo "✅ Windows 构建成功！"
    echo "输出文件: build/bin/magic-input-app.exe"
    ls -la build/bin/
else
    echo "❌ Windows 构建失败！"
    exit 1
fi

echo "构建完成！"