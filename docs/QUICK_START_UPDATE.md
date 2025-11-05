# Wails 自动更新 - 快速开始

## 🚀 快速设置自动更新

### 1. 环境准备

确保已安装必要工具：

```bash
# 检查 Wails 环境
wails doctor

# 安装依赖
make install
```

### 2. 配置更新源

#### GitHub 仓库设置

1. **更新 updater.go 中的仓库地址**：

```go
// 在 updater.go 第52行左右
updateURL: "https://api.github.com/repos/YOUR_USERNAME/YOUR_REPO/releases/latest"

// 在 fetchLatestVersion 方法中也要更新
```

2. **更新 config.go 中的默认配置**：

```go
// 在 GetDefaultUpdateConfig 函数中
UpdateURL: "https://api.github.com/repos/YOUR_USERNAME/YOUR_REPO/releases"
```

#### 示例配置

```go
// 将 YOUR_USERNAME 替换为你的 GitHub 用户名
// 将 YOUR_REPO 替换为你的仓库名
// 例如：
updateURL: "https://api.github.com/repos/john/magic-input-app/releases/latest"
UpdateURL: "https://api.github.com/repos/john/magic-input-app/releases"
```

### 3. 构建和测试

#### 本地构建测试

```bash
# 构建应用
make build-windows

# 查看版本信息
make version

# 运行应用测试更新功能
```

#### 发布新版本

```bash
# 1. 修改版本号（在 version.go 中）
# 2. 提交代码
git add .
git commit -m "feat: 更新到 v0.0.2"

# 3. 创建版本标签
git tag v0.0.2
git push origin v0.0.2

# 4. GitHub Actions 会自动构建并发布
```

### 4. 用户界面集成

应用已自动集成更新功能：

1. **自动检查**：应用启动 5 秒后自动检查更新
2. **手动检查**：右上角"检查更新"按钮
3. **更新对话框**：发现更新时自动显示
4. **设置页面**：管理自动更新配置

## 📋 功能清单

### ✅ 已实现功能

- [x] 自动版本检查
- [x] 手动检查更新
- [x] 下载进度显示
- [x] 自动安装更新
- [x] 应用重启
- [x] 版本跳过功能
- [x] 更新配置管理
- [x] 多平台支持
- [x] GitHub Actions 集成
- [x] 错误处理和重试

### 🎨 用户界面

- [x] 更新检查按钮
- [x] 更新对话框
- [x] 进度条显示
- [x] 设置页面
- [x] 响应式设计

### 🔧 开发工具

- [x] Makefile 构建脚本
- [x] 版本信息注入
- [x] 自动化构建流水线
- [x] 多平台交叉编译

## 🔍 测试自动更新

### 模拟更新流程

1. **创建测试版本**：

```bash
# 修改 version.go 中的版本号
Version = "0.0.1"

# 构建并运行应用
make build-windows
./build/bin/magic-input-app.exe
```

2. **发布新版本**：

```bash
# 修改版本号为更高版本
Version = "0.0.2"

# 创建 GitHub Release
git tag v0.0.2
git push origin v0.0.2
```

3. **测试更新检查**：
   - 点击"检查更新"按钮
   - 应用会发现新版本
   - 显示更新对话框

### 调试更新功能

#### 启用调试日志

```go
// 在 main.go 中添加
import "log"
log.SetFlags(log.LstdFlags | log.Lshortfile)
```

#### 查看配置文件

```bash
# macOS/Linux
cat ~/.magic-input-app/update_config.json

# Windows
type %USERPROFILE%\.magic-input-app\update_config.json
```

#### 手动触发检查

```javascript
// 在浏览器控制台中
window.go.main.App.CheckForUpdate().then(console.log);
```

## 🛠️ 自定义配置

### 修改检查间隔

```go
// 在 config.go 中修改默认配置
CheckInterval: 6, // 改为6小时检查一次
```

### 添加更新渠道

```go
// 支持不同的更新渠道
switch config.UpdateChannel {
case "beta":
    updateURL = "https://api.github.com/repos/USER/REPO/releases" // 包含预发布版本
case "stable":
    updateURL = "https://api.github.com/repos/USER/REPO/releases/latest" // 仅稳定版
}
```

### 自定义更新服务器

```go
// 可以替换为自己的更新服务器
UpdateURL: "https://your-update-server.com/api/check"
```

## 📦 发布流程

### 自动发布（推荐）

1. **配置 GitHub Actions**（已包含在项目中）
2. **推送版本标签**：

```bash
git tag v1.0.0
git push origin v1.0.0
```

3. **自动构建发布**：GitHub Actions 会自动构建并创建 Release

### 手动发布

1. **构建所有平台**：

```bash
make build-all
```

2. **创建发布包**：

```bash
make package-all
```

3. **上传到 GitHub Release**：
   - 手动创建 Release
   - 上传 `release/` 目录中的文件

## 🔒 安全注意事项

### 1. 验证更新来源

- 仅从官方 GitHub 仓库下载
- 验证 HTTPS 连接
- 检查文件完整性

### 2. 用户权限

- 需要应用写入权限
- 用户确认更新操作
- 安全的文件替换

### 3. 错误恢复

- 下载失败自动重试
- 更新失败回滚机制
- 详细错误日志

## 🎯 最佳实践

### 版本管理

- 使用语义化版本号
- 提供详细的更新日志
- 测试所有目标平台

### 用户体验

- 非侵入式更新提醒
- 清晰的进度指示
- 可配置的自动更新

### 开发流程

- 自动化构建和发布
- 完整的测试覆盖
- 统一的版本管理

## 🆘 故障排除

### 常见问题

1. **无法检查更新**

   - 检查网络连接
   - 验证 GitHub API 地址
   - 查看应用日志

2. **下载失败**

   - 检查磁盘空间
   - 确认网络稳定
   - 重试下载

3. **更新安装失败**
   - 检查应用权限
   - 确认文件未被占用
   - 手动重启应用

### 获取帮助

- 查看项目文档：`AUTO_UPDATE_GUIDE.md`
- 检查 GitHub Issues
- 查看应用日志文件

---

🎉 **恭喜！** 你的 Wails 应用现在具备了完整的自动更新功能！

用户现在可以：

- ✅ 自动获取最新版本
- ✅ 一键更新应用
- ✅ 自定义更新设置
- ✅ 享受持续的功能改进
