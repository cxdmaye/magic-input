# GitHub Actions 触发机制

## 📋 工作流概述

项目包含两个主要的 GitHub Actions 工作流：

### 1. 🚀 构建和发布 (`build.yml`)

**用途**: 正式版本的构建、打包和发布

**触发条件**:
- ✅ **仅限推送版本标签** (格式: `v*`)
  ```bash
  git tag v1.0.0
  git push origin v1.0.0
  ```

**执行内容**:
- 多平台构建 (Windows, Linux, macOS)
- 创建 GitHub Release
- 上传构建文件
- 自动生成发布说明

**文件输出**:
- `magic-input-app-windows-amd64.exe`
- `magic-input-app-linux-amd64`
- `magic-input-app-macos-amd64.tar.gz`

### 2. 🧪 测试构建 (`test-build.yml`)

**用途**: 开发过程中的快速验证和测试

**触发条件**:
- ✅ 推送到 `main` 分支
- ✅ 手动触发 (GitHub 网页操作)

**执行内容**:
- 代码编译验证
- Linux 平台构建测试
- 基本功能验证

## 🎯 使用场景

### 开发阶段

```bash
# 1. 日常开发提交 (触发测试构建)
git add .
git commit -m "feat: 添加新功能"
git push origin main
# → 自动运行测试构建验证代码

# 2. 手动测试特定分支
# 在 GitHub Actions 页面手动触发 "Test Build"
```

### 发布阶段

```bash
# 1. 准备发布版本
git add .
git commit -m "release: 准备发布 v1.2.0"
git push origin main

# 2. 创建发布标签 (触发构建和发布)
git tag v1.2.0
git push origin v1.2.0
# → 自动构建所有平台版本并创建 Release
```

## 🔒 安全考虑

### 为什么限制发布触发条件？

1. **防止意外发布**: 避免开发提交意外触发正式发布
2. **版本控制**: 确保每个 Release 都对应一个明确的版本标签
3. **资源节约**: 正式构建消耗更多 CI 资源，只在需要时运行
4. **发布质量**: 强制通过版本标签进行发布，保证版本管理规范

### 标签命名规范

支持的标签格式 (触发构建):
- ✅ `v1.0.0` (稳定版)
- ✅ `v1.0.0-beta.1` (测试版)
- ✅ `v1.0.0-alpha.1` (预览版)
- ✅ `v1.0.0-rc.1` (候选版)

不支持的格式 (不会触发):
- ❌ `1.0.0` (缺少 v 前缀)
- ❌ `release-1.0.0` (不匹配 v* 模式)
- ❌ `latest` (非版本号)

## 🔄 工作流程图

```
开发提交 → main 分支 → 测试构建 (验证代码)
                         ↓
                    开发继续...
                         ↓
版本准备完成 → 创建标签 → 构建发布 → GitHub Release
```

## 📊 构建矩阵

### 测试构建
- 平台: Ubuntu (Linux)
- 目的: 快速验证
- 时间: ~2-3 分钟

### 正式构建
- 平台: Ubuntu, Windows, macOS
- 目的: 完整发布包
- 时间: ~8-12 分钟

## 🛠️ 故障排除

### 常见问题

**Q: 推送了标签但没有触发构建？**
A: 检查标签格式是否为 `v*` 模式，确保推送到了正确的远程仓库

**Q: 想要取消错误的发布构建？**
A: 在 GitHub Actions 页面可以取消正在运行的工作流，并删除错误的标签

**Q: 如何修改已发布的版本？**
A: 删除标签和 Release，修复后重新创建标签

```bash
# 删除本地和远程标签
git tag -d v1.0.0
git push origin :refs/tags/v1.0.0

# 修复代码后重新发布
git tag v1.0.0
git push origin v1.0.0
```

---

这种设计确保了发布流程的严谨性和可控性，同时保持了开发过程的灵活性。