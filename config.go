package main

import (
	"encoding/json"
	"os"
	"path/filepath"
	"time"
)

// UpdateConfig 更新配置
type UpdateConfig struct {
	AutoCheck     bool   `json:"auto_check"`     // 是否自动检查更新
	CheckInterval int    `json:"check_interval"` // 检查间隔（小时）
	UpdateURL     string `json:"update_url"`     // 更新服务器地址
	SkipVersion   string `json:"skip_version"`   // 跳过的版本
	LastCheck     int64  `json:"last_check"`     // 上次检查时间戳
	UpdateChannel string `json:"update_channel"` // 更新渠道 (stable, beta, alpha)
}

// GetDefaultUpdateConfig 获取默认更新配置
func GetDefaultUpdateConfig() *UpdateConfig {
	return &UpdateConfig{
		AutoCheck:     true,
		CheckInterval: 24, // 每24小时检查一次
		UpdateURL:     "https://api.github.com/repos/cxdmaye/magic-input/releases",
		SkipVersion:   "",
		LastCheck:     0,
		UpdateChannel: "stable",
	}
}

// LoadUpdateConfig 加载更新配置
func LoadUpdateConfig() (*UpdateConfig, error) {
	configPath := getConfigPath()
	
	// 如果配置文件不存在，返回默认配置
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return GetDefaultUpdateConfig(), nil
	}
	
	data, err := os.ReadFile(configPath)
	if err != nil {
		return GetDefaultUpdateConfig(), err
	}
	
	var config UpdateConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return GetDefaultUpdateConfig(), err
	}
	
	return &config, nil
}

// SaveUpdateConfig 保存更新配置
func SaveUpdateConfig(config *UpdateConfig) error {
	configPath := getConfigPath()
	
	// 确保配置目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return err
	}
	
	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return err
	}
	
	return os.WriteFile(configPath, data, 0644)
}

// getConfigPath 获取配置文件路径
func getConfigPath() string {
	homeDir, _ := os.UserHomeDir()
	return filepath.Join(homeDir, ".magic-app", "update_config.json")
}

// ShouldCheckUpdate 判断是否应该检查更新
func ShouldCheckUpdate(config *UpdateConfig) bool {
	if !config.AutoCheck {
		return false
	}
	
	now := time.Now().Unix()
	checkInterval := int64(config.CheckInterval * 3600) // 转换为秒
	
	return now-config.LastCheck > checkInterval
}

// UpdateLastCheckTime 更新最后检查时间
func UpdateLastCheckTime(config *UpdateConfig) error {
	config.LastCheck = time.Now().Unix()
	return SaveUpdateConfig(config)
}

// GetUpdateConfig 获取更新配置 (前端调用)
func (a *App) GetUpdateConfig() (*UpdateConfig, error) {
	return LoadUpdateConfig()
}

// SetUpdateConfig 设置更新配置 (前端调用)
func (a *App) SetUpdateConfig(config *UpdateConfig) error {
	return SaveUpdateConfig(config)
}