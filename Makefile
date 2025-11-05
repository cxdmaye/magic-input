# Makefile for Magic Input
.PHONY: help build-windows build-darwin build-all clean dev install package-windows version test update-deps

# 版本信息
VERSION ?= $(shell git describe --tags --abbrev=0 2>/dev/null || echo "0.0.1")
BUILD_TIME := $(shell date -u '+%Y-%m-%d %H:%M:%S UTC')
GIT_COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")

# 构建标志
LDFLAGS := -ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"

# 默认目标
help:
	@echo "Magic Input 构建工具"
	@echo ""
	@echo "可用命令:"
	@echo "  build-windows  - 构建 Windows exe 文件"
	@echo "  build-darwin   - 构建 macOS 应用"
	@echo "  build-all      - 构建所有平台"
	@echo "  clean          - 清理构建文件"
	@echo "  dev            - 启动开发模式"
	@echo "  install        - 安装依赖"
	@echo "  package-windows- 创建 Windows 发布包"
	@echo "  test           - 运行测试"
	@echo "  update-deps    - 更新依赖"
	@echo ""
	@echo "版本信息:"
	@echo "  VERSION:    $(VERSION)"
	@echo "  BUILD_TIME: $(BUILD_TIME)"
	@echo "  GIT_COMMIT: $(GIT_COMMIT)"

# 构建 Windows 版本
build-windows:
	@echo "🚀 开始构建 Windows 版本..."
	@echo "版本: $(VERSION) | 提交: $(GIT_COMMIT)"
	@wails build -platform windows/amd64 -clean -o magic-input-app.exe $(LDFLAGS)
	@echo "✅ Windows 构建完成: build/bin/magic-input-app.exe"

# 构建 macOS 版本
build-darwin:
	@echo "🚀 开始构建 macOS 版本..."
	@echo "版本: $(VERSION) | 提交: $(GIT_COMMIT)"
	@wails build -platform darwin/amd64 -clean $(LDFLAGS)
	@echo "✅ macOS 构建完成: build/bin/magic-input-app.app"



# 构建所有平台
build-all: build-windows build-darwin

# 清理构建文件
clean:
	@echo "🧹 清理构建文件..."
	@rm -rf build/bin/
	@rm -rf frontend/dist/
	@rm -rf frontend/node_modules/.vite/
	@rm -rf release/
	@echo "✅ 清理完成"

# 开发模式
dev:
	@echo "🔧 启动开发模式..."
	@wails dev

# 安装依赖
install:
	@echo "📦 安装 Go 依赖..."
	@go mod tidy
	@echo "📦 安装前端依赖..."
	@cd frontend && npm install
	@echo "✅ 依赖安装完成"

# 更新依赖
update-deps:
	@echo "🔄 更新 Go 依赖..."
	@go get -u ./...
	@go mod tidy
	@echo "🔄 更新前端依赖..."
	@cd frontend && npm update
	@echo "✅ 依赖更新完成"

# 创建发布包
package-windows: build-windows
	@echo "📦 创建 Windows 发布包..."
	@mkdir -p release
	@cp build/bin/magic-input-app.exe release/
	@echo "✅ Windows 发布包创建完成: release/magic-input-app.exe"

package-darwin: build-darwin
	@echo "📦 创建 macOS 发布包..."
	@mkdir -p release
	@cp -r build/bin/magic-input-app.app release/
	@echo "✅ macOS 发布包创建完成: release/magic-input-app.app"

# 创建所有发布包
package-all: package-windows package-darwin package-linux

# 运行测试
test:
	@echo "🧪 运行 Go 测试..."
	@go test ./...
	@echo "🧪 运行前端测试..."
	@cd frontend && npm test
	@echo "✅ 测试完成"

# 显示版本信息
version:
	@echo "版本信息:"
	@echo "  VERSION:    $(VERSION)"
	@echo "  BUILD_TIME: $(BUILD_TIME)"
	@echo "  GIT_COMMIT: $(GIT_COMMIT)"
	@echo "  GO_VERSION: $(shell go version)"

# 创建新版本标签
tag:
	@read -p "请输入新版本号 (当前: $(VERSION)): " version; \
	if [ -n "$$version" ]; then \
		git tag "v$$version"; \
		git push origin "v$$version"; \
		echo "✅ 版本标签 v$$version 已创建并推送"; \
	else \
		echo "❌ 版本号不能为空"; \
	fi

# 生成绑定文件
generate:
	@echo "🔧 生成 TypeScript 绑定文件..."
	@wails generate bindings
	@echo "✅ 绑定文件生成完成"