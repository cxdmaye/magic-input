# Wails 项目自动更新实现指南

## 概述

本项目实现了一个完整的自动更新系统，支持：

- 自动检查更新
- 手动检查更新
- 下载进度显示
- 自动安装和重启
- 更新配置管理
- 版本跳过功能

## 架构组件

### 后端组件

#### 1. updater.go

- `UpdateInfo`: 更新信息结构
- `UpdateStatus`: 更新状态
- `UpdateProgress`: 下载进度
- `Updater`: 核心更新器类
- 提供版本检查、下载安装等功能

#### 2. update_manager.go

- `UpdateManager`: 更新管理器
- 负责定时检查、初始化检查
- 管理更新流程和事件

#### 3. config.go

- `UpdateConfig`: 更新配置结构
- 配置文件管理（加载、保存）
- 自动更新设置管理

#### 4. version.go

- 版本信息管理
- 构建时注入版本信息

### 前端组件

#### 1. UpdateDialog.tsx

- 更新对话框 UI
- 显示版本信息、更新内容
- 下载进度条
- 用户交互（更新、跳过、取消）

#### 2. UpdateSettings.tsx

- 更新设置页面
- 自动检查开关
- 检查间隔设置
- 更新渠道选择

#### 3. useUpdater.ts

- 更新功能 Hook
- 事件监听和状态管理
- 与后端 API 交互

## 实现流程

### 1. 应用启动

```
应用启动 → 创建UpdateManager → 启动定时检查 → 初始检查（延迟5秒）
```

### 2. 检查更新

```
检查时间间隔 → 请求GitHub API → 比较版本 → 发送事件通知
```

### 3. 下载更新

```
用户确认 → 下载文件 → 显示进度 → 应用更新 → 提示重启
```

### 4. 配置管理

```
加载配置 → 显示设置界面 → 用户修改 → 保存配置 → 重新启动管理器
```

## 使用方法

### 配置更新源

1. 修改 `updater.go` 中的 `updateURL`：

```go
updateURL: "https://api.github.com/repos/YOUR_USERNAME/YOUR_REPO/releases/latest"
```

2. 修改 `config.go` 中的默认配置：

```go
UpdateURL: "https://api.github.com/repos/YOUR_USERNAME/YOUR_REPO/releases"
```

### 发布新版本

1. **创建 Git 标签**：

```bash
git tag v1.0.1
git push origin v1.0.1
```

2. **GitHub Actions 自动构建**：

   - 推送标签后自动触发构建
   - 生成多平台应用文件
   - 创建 GitHub Release

3. **手动发布**：

```bash
# 构建应用
make build-windows
make build-darwin

# 创建 Release 并上传文件
```

### 版本命名规范

- 使用语义化版本：`v1.0.0`, `v1.0.1`, `v1.1.0`
- 测试版：`v1.0.0-beta.1`
- 预览版：`v1.0.0-alpha.1`

## API 接口

### 后端方法

```go
// 检查更新
func (a *App) CheckForUpdate() (*UpdateStatus, error)

// 下载并安装更新
func (a *App) DownloadAndInstallUpdate(downloadURL string) error

// 重启应用
func (a *App) RestartApp() error

// 获取版本信息
func (a *App) GetVersionInfo() map[string]string

// 配置管理
func (a *App) GetUpdateConfig() (*UpdateConfig, error)
func (a *App) SetUpdateConfig(config *UpdateConfig) error

// 版本管理
func (a *App) SkipVersion(version string) error
func (a *App) ResetSkipVersion() error
func (a *App) SetAutoUpdate(enabled bool) error
func (a *App) GetAutoUpdateEnabled() (bool, error)
```

### 前端事件

```typescript
// 更新可用
EventsOn("update:available", (data: UpdateStatus) => {});

// 无更新
EventsOn("update:no-update", (data: UpdateStatus) => {});

// 检查错误
EventsOn("update:check:error", (data: { error: string }) => {});

// 下载开始
EventsOn("update:download:start", (data: { status: string }) => {});

// 下载进度
EventsOn("update:download:progress", (progress: UpdateProgress) => {});

// 下载完成
EventsOn("update:download:complete", (data: { status: string }) => {});

// 下载错误
EventsOn("update:download:error", (data: { error: string }) => {});
```

## 配置选项

### 自动更新配置

```json
{
  "auto_check": true,
  "check_interval": 24,
  "update_url": "https://api.github.com/repos/USER/REPO/releases",
  "skip_version": "",
  "last_check": 1635724800,
  "update_channel": "stable"
}
```

### 参数说明

- `auto_check`: 是否启用自动检查
- `check_interval`: 检查间隔（小时）
- `update_url`: 更新服务器地址
- `skip_version`: 跳过的版本号
- `last_check`: 上次检查时间戳
- `update_channel`: 更新渠道（stable/beta/alpha）

## 构建配置

### Makefile 更新

已集成版本信息注入：

```makefile
LDFLAGS := -ldflags "-X 'main.Version=$(VERSION)' -X 'main.BuildTime=$(BUILD_TIME)' -X 'main.GitCommit=$(GIT_COMMIT)'"
```

### GitHub Actions

自动构建工作流支持：

- 多平台构建（Windows/macOS/Linux）
- 自动发布 Release
- 版本标签触发

## 安全考虑

### 1. 下载验证

- 验证下载文件来源
- 检查文件完整性
- HTTPS 传输

### 2. 权限管理

- 最小权限原则
- 用户确认更新
- 安全的文件替换

### 3. 错误处理

- 网络错误重试
- 下载失败回滚
- 更新失败恢复

## 故障排除

### 常见问题

1. **无法检查更新**

   - 检查网络连接
   - 验证 GitHub API 地址
   - 查看错误日志

2. **下载失败**

   - 检查存储空间
   - 验证文件权限
   - 重试下载

3. **更新失败**
   - 检查应用权限
   - 确认文件未被占用
   - 手动重启应用

### 调试模式

启用详细日志：

```go
log.SetLevel(log.DebugLevel)
```

查看配置文件：

```
~/.magic-input-app/update_config.json
```

## 最佳实践

### 1. 版本发布

- 遵循语义化版本
- 提供详细的更新日志
- 测试所有平台

### 2. 用户体验

- 非侵入式更新提醒
- 清晰的进度指示
- 可选的自动更新

### 3. 错误处理

- 优雅的降级机制
- 详细的错误信息
- 用户友好的提示

## 扩展功能

### 可能的增强

- 增量更新支持
- 更新回滚功能
- 自定义更新服务器
- 企业版本管理
- 更新统计分析

这个自动更新系统为 Wails 应用提供了完整的更新解决方案，确保用户始终使用最新版本的应用程序。
