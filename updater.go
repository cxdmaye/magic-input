package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	goruntime "runtime"
	"strings"
	"time"

	"github.com/Masterminds/semver/v3"
	"github.com/inconshreveable/go-update"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// UpdateInfo 更新信息结构
type UpdateInfo struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
	Changelog   string `json:"changelog"`
	Required    bool   `json:"required"`
	PublishedAt string `json:"published_at"`
}

// UpdateStatus 更新状态
type UpdateStatus struct {
	Available   bool   `json:"available"`
	CurrentVer  string `json:"current_version"`
	LatestVer   string `json:"latest_version"`
	DownloadURL string `json:"download_url"`
	Changelog   string `json:"changelog"`
	Required    bool   `json:"required"`
}

// UpdateProgress 更新进度
type UpdateProgress struct {
	Percent     int    `json:"percent"`
	Speed       string `json:"speed"`
	ETA         string `json:"eta"`
	Status      string `json:"status"`
	BytesTotal  int64  `json:"bytes_total"`
	BytesLoaded int64  `json:"bytes_loaded"`
}

// Updater 更新器结构
type Updater struct {
	app        *App
	updateURL  string
	currentVer *semver.Version
}

// NewUpdater 创建新的更新器
func NewUpdater(app *App) *Updater {
	currentVer, err := semver.NewVersion(Version)
	if err != nil {
		// 如果版本格式不正确，使用默认版本
		currentVer = semver.MustParse("0.0.1")
	}
	
	return &Updater{
		app:        app,
		updateURL:  "https://api.github.com/repos/cxdmaye/magic-input/releases/latest",
		currentVer: currentVer,
	}
}

// CheckForUpdate 检查更新
func (a *App) CheckForUpdate() (*UpdateStatus, error) {
	updater := NewUpdater(a)
	
	// 获取最新版本信息
	updateInfo, err := updater.fetchLatestVersion()
	if err != nil {
		return nil, fmt.Errorf("获取版本信息失败: %v", err)
	}
	
	// 比较版本
	latestVer, err := semver.NewVersion(updateInfo.Version)
	if err != nil {
		return nil, fmt.Errorf("解析版本号失败: %v", err)
	}
	
	available := latestVer.GreaterThan(updater.currentVer)
	
	return &UpdateStatus{
		Available:   available,
		CurrentVer:  updater.currentVer.String(),
		LatestVer:   latestVer.String(),
		DownloadURL: updateInfo.DownloadURL,
		Changelog:   updateInfo.Changelog,
		Required:    updateInfo.Required,
	}, nil
}

// DownloadAndInstallUpdate 下载并安装更新
func (a *App) DownloadAndInstallUpdate(downloadURL string) error {
	// 发送开始下载事件
	runtime.EventsEmit(a.ctx, "update:download:start", map[string]interface{}{
		"status": "开始下载更新...",
	})
	
	// 下载更新文件
	resp, err := http.Get(downloadURL)
	if err != nil {
		return fmt.Errorf("下载失败: %v", err)
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("下载失败，状态码: %d", resp.StatusCode)
	}
	
	// 创建进度读取器
	totalBytes := resp.ContentLength
	progressReader := &ProgressReader{
		Reader:     resp.Body,
		Total:      totalBytes,
		OnProgress: func(current, total int64, percent int) {
			runtime.EventsEmit(a.ctx, "update:download:progress", UpdateProgress{
				Percent:     percent,
				Status:      "正在下载...",
				BytesTotal:  total,
				BytesLoaded: current,
			})
		},
	}
	
	// 根据平台选择不同的更新策略
	if goruntime.GOOS == "darwin" {
		err = a.applyMacOSUpdate(progressReader, downloadURL)
	} else {
		// Windows 和 Linux 使用标准方式
		err = update.Apply(progressReader, update.Options{})
	}
	
	if err != nil {
		runtime.EventsEmit(a.ctx, "update:download:error", map[string]interface{}{
			"error": err.Error(),
		})
		return fmt.Errorf("应用更新失败: %v", err)
	}
	
	// 发送完成事件
	runtime.EventsEmit(a.ctx, "update:download:complete", map[string]interface{}{
		"status": "更新下载完成，请重启应用",
	})
	
	return nil
}

// applyMacOSUpdate macOS 专用的更新逻辑
func (a *App) applyMacOSUpdate(reader io.Reader, downloadURL string) error {
	// 检查是否是 pkg 文件
	if strings.Contains(downloadURL, ".pkg") {
		return a.downloadAndInstallPkg(reader, downloadURL)
	}
	
	// 对于其他文件类型，尝试标准更新，如果失败则提供手动安装指导
	err := update.Apply(reader, update.Options{})
	if err != nil {
		// 如果标准更新失败，提供更详细的错误信息和解决方案
		if strings.Contains(err.Error(), "permission") || strings.Contains(err.Error(), "denied") {
			return fmt.Errorf("权限不足，无法自动更新。请：\n1. 手动下载新版本\n2. 或者以管理员权限运行应用\n3. 或者将应用移动到应用程序文件夹\n\n原始错误: %v", err)
		}
		return err
	}
	return nil
}

// downloadAndInstallPkg 下载并安装 pkg 文件（macOS 专用）
func (a *App) downloadAndInstallPkg(reader io.Reader, downloadURL string) error {
	// 创建临时文件
	tmpDir := os.TempDir()
	tmpFile := filepath.Join(tmpDir, "magic-input-update.pkg")
	
	// 将下载内容写入临时文件
	file, err := os.Create(tmpFile)
	if err != nil {
		return fmt.Errorf("创建临时文件失败: %v", err)
	}
	defer file.Close()
	defer os.Remove(tmpFile) // 清理临时文件
	
	_, err = io.Copy(file, reader)
	if err != nil {
		return fmt.Errorf("保存下载文件失败: %v", err)
	}
	file.Close()
	
	// 使用系统的 installer 命令安装 pkg
	// 这需要用户输入管理员密码
	cmd := exec.Command("osascript", "-e", fmt.Sprintf(`
		do shell script "installer -pkg '%s' -target /" with administrator privileges
	`, tmpFile))
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("安装失败: %v\n输出: %s", err, string(output))
	}
	
	return nil
}

// AutoCheckUpdate 自动检查更新
func (a *App) AutoCheckUpdate() {
	go func() {
		// 延迟5秒后开始检查，避免影响应用启动
		time.Sleep(5 * time.Second)
		
		updateStatus, err := a.CheckForUpdate()
		if err != nil {
			runtime.EventsEmit(a.ctx, "update:check:error", map[string]interface{}{
				"error": err.Error(),
			})
			return
		}
		
		if updateStatus.Available {
			runtime.EventsEmit(a.ctx, "update:available", updateStatus)
		} else {
			runtime.EventsEmit(a.ctx, "update:no-update", updateStatus)
		}
	}()
}

// RestartApp 重启应用
func (a *App) RestartApp() error {
	// 获取当前执行文件路径
	executable, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取执行文件路径失败: %v", err)
	}
	
	var cmd *exec.Cmd
	
	if goruntime.GOOS == "darwin" {
		// macOS 特殊处理
		if strings.Contains(executable, ".app/Contents/MacOS/") {
			// 如果是在 .app 包内运行，使用 open 命令重启整个应用包
			appPath := executable
			for strings.Contains(appPath, "/Contents/MacOS/") {
				appPath = filepath.Dir(appPath)
			}
			// 找到 .app 目录
			for !strings.HasSuffix(appPath, ".app") && appPath != "/" {
				appPath = filepath.Dir(appPath)
			}
			
			if strings.HasSuffix(appPath, ".app") {
				// 使用 open 命令重新启动应用
				cmd = exec.Command("open", "-n", appPath)
			} else {
				// 回退到直接执行
				cmd = exec.Command(executable, os.Args[1:]...)
			}
		} else {
			// 非 .app 包情况，直接执行
			cmd = exec.Command(executable, os.Args[1:]...)
		}
	} else {
		// Windows 和 Linux
		cmd = exec.Command(executable, os.Args[1:]...)
	}
	
	// 设置命令属性
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	
	// 启动新进程
	err = cmd.Start()
	if err != nil {
		return fmt.Errorf("重启应用失败: %v", err)
	}
	
	// 退出当前进程
	runtime.Quit(a.ctx)
	return nil
}

// fetchLatestVersion 获取最新版本信息
func (u *Updater) fetchLatestVersion() (*UpdateInfo, error) {
	fmt.Println("更新请求:", u.updateURL)
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	
	resp, err := client.Get(u.updateURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}
	
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	
	// 解析GitHub API响应
	var release struct {
		TagName     string `json:"tag_name"`
		Body        string `json:"body"`
		Assets      []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		} `json:"assets"`
		PublishedAt string `json:"published_at"`
		Prerelease  bool   `json:"prerelease"`
	}
	
	err = json.Unmarshal(body, &release)
	if err != nil {
		return nil, err
	}
	
	// 查找对应平台的下载链接
	var downloadURL string
	platform := goruntime.GOOS
	
	for _, asset := range release.Assets {
		fmt.Println("Asset Name:", asset.Name)
		if platform == "windows" && (asset.Name == "magic-input-app-windows-amd64.exe" || 
			fmt.Sprintf("magic-input-app-%s.exe", platform) == asset.Name) {
			downloadURL = asset.BrowserDownloadURL
			break
		} else if platform == "darwin" && (asset.Name == "magic-input-app-macos-amd64.pkg" ||
			fmt.Sprintf("magic-input-app-%s.app", platform) == asset.Name) {
			downloadURL = asset.BrowserDownloadURL
			break
		} else if platform == "linux" && (asset.Name == "magic-input-app" ||
			fmt.Sprintf("magic-input-app-%s", platform) == asset.Name) {
			downloadURL = asset.BrowserDownloadURL
			break
		}
	}
	
	if downloadURL == "" {
		return nil, fmt.Errorf("未找到适用于 %s 平台的更新文件", platform)
	}
	
	// 移除版本号前的 'v' 前缀（如果存在）
	version := release.TagName
	if len(version) > 0 && version[0] == 'v' {
		version = version[1:]
	}
	
	return &UpdateInfo{
		Version:     version,
		DownloadURL: downloadURL,
		Changelog:   release.Body,
		Required:    false, // 可以根据需要设置为强制更新
		PublishedAt: release.PublishedAt,
	}, nil
}

// ProgressReader 进度读取器
type ProgressReader struct {
	io.Reader
	Total      int64
	Current    int64
	OnProgress func(current, total int64, percent int)
}

func (pr *ProgressReader) Read(p []byte) (n int, err error) {
	n, err = pr.Reader.Read(p)
	pr.Current += int64(n)
	
	if pr.OnProgress != nil && pr.Total > 0 {
		percent := int((pr.Current * 100) / pr.Total)
		if percent > 100 {
			percent = 100
		}
		pr.OnProgress(pr.Current, pr.Total, percent)
	}
	
	return
}