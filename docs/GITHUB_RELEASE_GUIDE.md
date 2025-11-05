# GitHub Actions 发布指南

## 🔧 已修复的问题

✅ **升级了所有 GitHub Actions 到最新版本**：

- `actions/checkout@v4`
- `actions/setup-go@v5`
- `actions/setup-node@v4`
- `actions/upload-artifact@v4`
- `actions/download-artifact@v4`
- `softprops/action-gh-release@v2`

✅ **解决了 deprecated actions 警告**

✅ **添加了更多错误处理和验证**

## 🚀 如何发布新版本

### 📋 触发条件

**构建和发布** (`build.yml`) 仅在以下情况触发：
- ✅ 推送版本标签 (格式: `v*`，如 `v1.0.0`)

**测试构建** (`test-build.yml`) 在以下情况触发：
- ✅ 推送到 `main` 分支
- ✅ 手动触发 (GitHub 网页操作)

### 方法一：创建 Git 标签（唯一发布方式）

```bash
# 1. 确保所有代码已提交
git add .
git commit -m "feat: 准备发布 v0.1.0"
git push origin main

# 2. 创建版本标签
git tag v0.1.0
git push origin v0.1.0
```

推送标签后，GitHub Actions 会自动：
- 构建 Windows、Linux、macOS 版本
- 创建 GitHub Release
- 上传所有构建文件
- 生成发布说明

### 开发测试方法

如需在开发过程中测试构建，有以下选项：

#### 1. 自动测试构建
推送代码到 `main` 分支会自动触发测试构建（仅 Linux 版本）：
```bash
git push origin main
```

#### 2. 手动触发测试

1. 访问 GitHub 仓库页面
2. 点击 "Actions" 标签
3. 选择 "Test Build" 工作流
4. 点击 "Run workflow"
5. 选择分支并运行

**注意**: 只有测试构建支持手动触发，正式发布必须通过推送标签。

## 📦 构建产物

构建完成后，Release 中会包含：

- **Windows**: `magic-input-app-windows-amd64.exe`
- **Linux**: `magic-input-app-linux-amd64`
- **macOS**: `magic-input-app-macos-amd64.tar.gz`

## 🧪 测试构建

我们还创建了一个测试工作流 (`test-build.yml`)，用于：

- 验证代码编译
- 测试 Linux 构建
- 检查构建产物

## ⚙️ 高级功能

### 自动版本注入

构建时会自动注入以下信息：

- 版本号（来自 Git 标签）
- 构建时间
- Git 提交哈希

### 预发布版本

如果标签包含以下关键词，会标记为预发布：

- `beta`
- `alpha`
- `rc`

例如：`v1.0.0-beta.1`

### 发布说明

自动生成的发布说明包含：

- 下载链接
- 安装说明
- 构建信息
- 功能列表

## 🔍 故障排除

### 查看构建日志

1. 访问 GitHub 仓库
2. 点击 "Actions" 标签
3. 选择失败的工作流
4. 查看详细日志

### 常见问题

**问题**: 构建失败 "command not found: wails"
**解决**: Wails CLI 安装失败，通常是网络问题，重新运行即可

**问题**: 前端依赖安装失败
**解决**: 检查 `package-lock.json` 是否存在，或删除后重新生成

**问题**: 找不到构建文件
**解决**: 检查 `wails.json` 配置和构建输出路径

## 📝 版本命名建议

- **主版本**: `v1.0.0`（重大更新）
- **次版本**: `v1.1.0`（新功能）
- **补丁版本**: `v1.0.1`（bug 修复）
- **预发布**: `v1.0.0-beta.1`（测试版本）

## 🎯 下一步

1. **立即测试**：推送一个测试标签验证构建流程

   ```bash
   git tag v0.0.2-test
   git push origin v0.0.2-test
   ```

2. **发布正式版本**：确认一切正常后发布 v1.0.0

3. **配置自动更新**：用户下载后可以自动检测新版本

---

现在你的 GitHub Actions 已经完全现代化，支持自动构建和发布！🎉
