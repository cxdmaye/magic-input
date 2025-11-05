package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/openai/openai-go"
	"github.com/openai/openai-go/option"
)

const defaultPrompt = `# 角色
你是一个专业的搜索助手，能够根据用户的输入内容来猜测该内容可能是什么样格式内容，并提供精准的判断和建议方向。

## 技能：判断内容可能的格式
用户输入的内容可能来自于其他地方的粘贴，在粘贴之前可能不知道这段内容代表什么，以下是可能的一些类型：
1. 网络资源类
	- 网址
	- 图片、视频、音频等
	- PDF、EXCEL、CSV 等文档
2. 纯文本
	- 颜色值，rgb,rgba,hsl等
	- 数学公式，Latex公式等
	- 外语单词/句子等
	- 文言文 出处查询
	- 诗词 查询上下文
	- 普通文本 如表情，人名，地名等
	- 时间戳
	- 除 10 进制以外的进制数  转换为 10进制
3. 可执行文本
	- 网络资源请求
	- bash 命令
	- ws 等连接
	- curl 请求
4. 代码
	- 脚本语言
	- 样式/标记语言 HTML/CSS
	- SQL 语言
	- XML、YAML 等
	- 配置脚本 如 dockerFile/package.json等
    - JSON 代码
5. 加密文本
	- MD5
	- base64
	- 虚拟币地址
  用户的输入包括但不限于以上的各种类型，你可以通过分析来判断用户的输入可能的类型，并提供每一种类型的可能概率

## 输出格式
` + "```json" + `
[{
    "classification":"plaintext",
    "type":"color",
    "classify":"颜色值",
    "percent":"10%",
    "result":"可能的结果",
    "suggestion":"针对该类型的建议"
}]
` + "```" + `

## 限制:
- 输出内容尽可能简洁，最好是 JSON 格式。输出的 type 按照一定的格式进行英文转换，result使用中文，展示这种猜测可能的结果`

// AIConfig AI配置结构
type AIConfig struct {
	APIKey  string `json:"api_key"`
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
}

// AIAnalysisResult AI分析结果
type AIAnalysisResult struct {
	Classification string `json:"classification"`
	Type           string `json:"type"`
	Classify       string `json:"classify"`
	Percent        string `json:"percent"`
	Result         string `json:"result"`
	Suggestion     string `json:"suggestion"`
}

// AIService AI服务
type AIService struct {
	app    *App
	config *AIConfig
}

// NewAIService 创建AI服务
func NewAIService(app *App) *AIService {
	return &AIService{
		app: app,
		config: &AIConfig{
			BaseURL: "https://dashscope.aliyuncs.com/compatible-mode/v1/",
			Model:   "qwen-plus",
		},
	}
}

// GetAIConfig 获取AI配置
func (a *App) GetAIConfig() (*AIConfig, error) {
	return a.aiService.config, nil
}

// SaveAIConfig 保存AI配置
func (a *App) SaveAIConfig(config AIConfig) error {
	a.aiService.config = &config
	
	// 保存到本地配置文件
	configData, err := json.Marshal(config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %v", err)
	}
	
	// 直接保存到文件
	return a.saveConfigToFile("ai_config.json", configData)
}

// LoadAIConfig 加载AI配置
func (a *App) LoadAIConfig() error {
	// 直接从文件加载
	return a.loadConfigFromFile("ai_config.json")
}

// AnalyzeContent 分析内容 - 实际调用 AI 接口
func (a *App) AnalyzeContent(content string) ([]AIAnalysisResult, error) {
	if a.aiService.config.APIKey == "" {
		return nil, fmt.Errorf("请先配置 API Key")
	}
	
	// 调用 AI 接口
	response, err := a.callAIAPI(content)
	if err != nil {
		return nil, fmt.Errorf("调用 AI 接口失败: %v", err)
	}
	
	// 解析 AI 响应
	results, err := a.parseAIResponse(response)
	if err != nil {
		return nil, fmt.Errorf("解析 AI 响应失败: %v", err)
	}
	
	return results, nil
}

// callAIAPI 调用 AI 接口
func (a *App) callAIAPI(userContent string) (string, error) {
	config := a.aiService.config
	
	fmt.Println("AI 原始请求:", userContent)
	
	// 创建 OpenAI 客户端
	client := openai.NewClient(
		option.WithAPIKey(config.APIKey),
		option.WithBaseURL(config.BaseURL),
	)
	
	// 调用 Chat Completions API
	chatCompletion, err := client.Chat.Completions.New(
		context.TODO(), 
		openai.ChatCompletionNewParams{
			Messages: []openai.ChatCompletionMessageParamUnion{
				openai.SystemMessage(defaultPrompt),
				openai.UserMessage(userContent),
			},
			Model:       config.Model,
			Temperature: openai.Float(0.7),
			MaxTokens:   openai.Int(2000),
		},
	)
	
	if err != nil {
		return "", fmt.Errorf("调用 OpenAI API 失败: %v", err)
	}
	
	fmt.Println("AI 原始响应:", chatCompletion.Choices[0].Message.Content)
	
	return chatCompletion.Choices[0].Message.Content, nil
}

// parseAIResponse 解析 AI 响应
func (a *App) parseAIResponse(aiResponse string) ([]AIAnalysisResult, error) {
	// 尝试从响应中提取 JSON 部分
	jsonStart := strings.Index(aiResponse, "[")
	jsonEnd := strings.LastIndex(aiResponse, "]")
	
	if jsonStart == -1 || jsonEnd == -1 || jsonStart >= jsonEnd {
		// 如果没有找到 JSON 格式，返回一个包含原始响应的结果
		return []AIAnalysisResult{
			{
				Classification: "unknown",
				Type:           "text",
				Classify:       "未知类型",
				Percent:        "50%",
				Result:         "AI 分析结果",
				Suggestion:     aiResponse,
			},
		}, nil
	}
	
	jsonStr := aiResponse[jsonStart : jsonEnd+1]
	
	// 解析 JSON
	var results []AIAnalysisResult
	err := json.Unmarshal([]byte(jsonStr), &results)
	if err != nil {
		// 如果 JSON 解析失败，返回包含原始响应的结果
		return []AIAnalysisResult{
			{
				Classification: "unknown",
				Type:           "text",
				Classify:       "解析错误",
				Percent:        "0%",
				Result:         "JSON 解析失败",
				Suggestion:     fmt.Sprintf("原始响应: %s", aiResponse),
			},
		}, nil
	}
	
	// 验证和修正结果
	for i := range results {
		if results[i].Classification == "" {
			results[i].Classification = "unknown"
		}
		if results[i].Type == "" {
			results[i].Type = "text"
		}
		if results[i].Classify == "" {
			results[i].Classify = "未知类型"
		}
		if results[i].Percent == "" {
			results[i].Percent = "50%"
		}
		if results[i].Result == "" {
			results[i].Result = "无法确定"
		}
		if results[i].Suggestion == "" {
			results[i].Suggestion = "建议进一步分析"
		}
	}
	
	return results, nil
}

// getConfigDir 获取配置目录路径
func (a *App) getConfigDir() (string, error) {
	// 优先使用用户配置目录
	configDir, err := os.UserConfigDir()
	if err != nil {
		// 如果获取用户配置目录失败，使用用户主目录
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("无法获取用户目录: %v", err)
		}
		configDir = homeDir
	}
	
	// 创建应用专用配置目录
	appConfigDir := filepath.Join(configDir, "magic-input-app")
	return appConfigDir, nil
}

// 辅助函数
func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 辅助方法：保存配置到文件
func (a *App) saveConfigToFile(filename string, data []byte) error {
	configDir, err := a.getConfigDir()
	if err != nil {
		return err
	}
	
	// 创建配置目录
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("创建配置目录失败: %v", err)
	}
	
	filePath := filepath.Join(configDir, filename)
	return os.WriteFile(filePath, data, 0644)
}

// 辅助方法：从文件加载配置
func (a *App) loadConfigFromFile(filename string) error {
	configDir, err := a.getConfigDir()
	if err != nil {
		return err
	}
	
	filePath := filepath.Join(configDir, filename)
	data, err := os.ReadFile(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，使用默认配置
		}
		return err
	}
	
	var config AIConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return fmt.Errorf("解析配置文件失败: %v", err)
	}
	
	a.aiService.config = &config
	return nil
}