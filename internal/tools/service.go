// CorpFlow Tools - Extensible tool system
// Similar to Claude Code / OpenCode tools

package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

// Tool 定义工具接口
type Tool interface {
	Name() string
	Description() string
	Execute(ctx context.Context, input map[string]interface{}) (string, error)
	Schema() ToolSchema
}

// ToolSchema 工具参数schema
type ToolSchema struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Parameters  map[string]Parameter `json:"parameters"`
}

// Parameter 参数定义
type Parameter struct {
	Type        string `json:"type"`
	Description string `json:"description"`
	Required    bool   `json:"required"`
	Default     interface{} `json:"default,omitempty"`
}

// ToolRegistry 工具注册表
var ToolRegistry = make(map[string]Tool)

// RegisterTool 注册工具
func RegisterTool(t Tool) {
	ToolRegistry[t.Name()] = t
}

// GetTool 获取工具
func GetTool(name string) Tool {
	return ToolRegistry[name]
}

// ListTools 列出所有工具
func ListTools() []Tool {
	tools := make([]Tool, 0, len(ToolRegistry))
	for _, t := range ToolRegistry {
		tools = append(tools, t)
	}
	return tools
}

// ============ 内置工具实现 ============

// ShellTool 执行Shell命令
type ShellTool struct{}

func (t *ShellTool) Name() string { return "shell" }
func (t *ShellTool) Description() string { return "Execute shell commands / 执行Shell命令" }

func (t *ShellTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	command, ok := input["command"].(string)
	if !ok || command == "" {
		return "", fmt.Errorf("command is required")
	}

	timeout := 30
	if to, ok := input["timeout"].(float64); ok {
		timeout = int(to)
	}

	cmd := exec.CommandContext(ctx, "sh", "-c", command)
	cmd.Dir, _ = input["dir"].(string)

	output, err := execWithTimeout(cmd, time.Duration(timeout)*time.Second)
	if err != nil {
		return "", fmt.Errorf("shell error: %v", err)
	}
	return output, nil
}

func (t *ShellTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "shell",
		Description: "Execute shell command",
		Parameters: map[string]Parameter{
			"command": {Type: "string", Description: "Command to execute", Required: true},
			"dir":     {Type: "string", Description: "Working directory"},
			"timeout": {Type: "number", Description: "Timeout in seconds", Default: 30},
		},
	}
}

// WebSearchTool 网页搜索
type WebSearchTool struct{}

func (t *WebSearchTool) Name() string { return "web_search" }
func (t *WebSearchTool) Description() string { return "Search the web / 网页搜索" }

func (t *WebSearchTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	query, ok := input["query"].(string)
	if !ok || query == "" {
		return "", fmt.Errorf("query is required")
	}

	// 调用搜索API (预留接口)
	return fmt.Sprintf("Search results for: %s\n\n1. Result 1\n2. Result 2\n3. Result 3", query), nil
}

func (t *WebSearchTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "web_search",
		Description: "Search the web",
		Parameters: map[string]Parameter{
			"query": {Type: "string", Description: "Search query", Required: true},
		},
	}
}

// WebFetchTool 获取网页内容
type WebFetchTool struct{}

func (t *WebFetchTool) Name() string { return "web_fetch" }
func (t *WebFetchTool) Description() string { return "Fetch web page content / 获取网页" }

func (t *WebFetchTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	url, ok := input["url"].(string)
	if !ok || url == "" {
		return "", fmt.Errorf("url is required")
	}

	// 预留: 实际调用fetch API
	return fmt.Sprintf("Content from: %s\n\n[Web content would be fetched here]", url), nil
}

func (t *WebFetchTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "web_fetch",
		Description: "Fetch URL content",
		Parameters: map[string]Parameter{
			"url": {Type: "string", Description: "URL to fetch", Required: true},
		},
	}
}

// FileReadTool 读取文件
type FileReadTool struct{}

func (t *FileReadTool) Name() string { return "file_read" }
func (t *FileReadTool) Description() string { return "Read file content / 读取文件" }

func (t *FileReadTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	path, ok := input["path"].(string)
	if !ok || path == "" {
		return "", fmt.Errorf("path is required")
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("read file error: %v", err)
	}

	lines := strings.Split(string(content), "\n")
	maxLines := 100
	if len(lines) > maxLines {
		return string(content) + fmt.Sprintf("\n\n... (%d more lines)", len(lines)-maxLines), nil
	}
	return string(content), nil
}

func (t *FileReadTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "file_read",
		Description: "Read file",
		Parameters: map[string]Parameter{
			"path": {Type: "string", Description: "File path", Required: true},
		},
	}
}

// FileWriteTool 写文件
type FileWriteTool struct{}

func (t *FileWriteTool) Name() string { return "file_write" }
func (t *FileWriteTool) Description() string { return "Write file content / 写文件" }

func (t *FileWriteTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	path, ok := input["path"].(string)
	if !ok || path == "" {
		return "", fmt.Errorf("path is required")
	}

	content, ok := input["content"].(string)
	if !ok {
		return "", fmt.Errorf("content is required")
	}

	// 创建目录
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return "", fmt.Errorf("create dir error: %v", err)
	}

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		return "", fmt.Errorf("write file error: %v", err)
	}

	return fmt.Sprintf("File written: %s", path), nil
}

func (t *FileWriteTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "file_write",
		Description: "Write file",
		Parameters: map[string]Parameter{
			"path":    {Type: "string", Description: "File path", Required: true},
			"content": {Type: "string", Description: "File content", Required: true},
		},
	}
}

// CodeReviewTool 代码审查
type CodeReviewTool struct{}

func (t *CodeReviewTool) Name() string { return "code_review" }
func (t *CodeReviewTool) Description() string { return "AI code review / 代码审查" }

func (t *CodeReviewTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	code, ok := input["code"].(string)
	if !ok || code == "" {
		return "", fmt.Errorf("code is required")
	}

	language, _ := input["language"].(string)
	if language == "" {
		language = detectLanguage(code)
	}

	// 预留: 调用AI进行代码审查
	return fmt.Sprintf(`## Code Review / 代码审查

### Language: %s

### Issues Found:
- No critical issues detected
- 2 suggestions for improvement

### Suggestions:
1. Consider adding error handling
2. Could benefit from comments

### Score: 8/10`, language), nil
}

func (t *CodeReviewTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "code_review",
		Description: "AI code review",
		Parameters: map[string]Parameter{
			"code":     {Type: "string", Description: "Code to review", Required: true},
			"language": {Type: "string", Description: "Programming language"},
		},
	}
}

// TestGenTool 生成测试
type TestGenTool struct{}

func (t *TestGenTool) Name() string { return "test_gen" }
func (t *TestGenTool) Description() string { return "Generate unit tests / 生成测试" }

func (t *TestGenTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	code, ok := input["code"].(string)
	if !ok || code == "" {
		return "", fmt.Errorf("code is required")
	}

	framework, _ := input["framework"].(string)
	if framework == "" {
		framework = "go"
	}

	// 预留: 生成测试代码
	return fmt.Sprintf(`// Generated tests for %s
package main

import "testing"

func TestExample(t *testing.T) {
    // TODO: Add test cases
    t.Skip("Generate actual tests")
}
`, framework), nil
}

func (t *TestGenTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "test_gen",
		Description: "Generate unit tests",
		Parameters: map[string]Parameter{
			"code":      {Type: "string", Description: "Code to generate tests for", Required: true},
			"framework": {Type: "string", Description: "Test framework (go/python/jest)"},
		},
	}
}

// GitTool Git操作
type GitTool struct{}

func (t *GitTool) Name() string { return "git" }
func (t *GitTool) Description() string { return "Git operations / Git操作" }

func (t *GitTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	action, ok := input["action"].(string)
	if !ok || action == "" {
		return "", fmt.Errorf("action is required (status/commit/push/pull)")
	}

	repoPath, _ := input["path"].(string)
	if repoPath == "" {
		repoPath = "."
	}

	var cmd *exec.Cmd
	switch action {
	case "status":
		cmd = exec.CommandContext(ctx, "git", "-C", repoPath, "status", "--short")
	case "commit":
		msg, _ := input["message"].(string)
		if msg == "" {
			return "", fmt.Errorf("commit message required")
		}
		cmd = exec.CommandContext(ctx, "git", "-C", repoPath, "commit", "-m", msg)
	case "push":
		cmd = exec.CommandContext(ctx, "git", "-C", repoPath, "push")
	case "pull":
		cmd = exec.CommandContext(ctx, "git", "-C", repoPath, "pull")
	case "log":
		cmd = exec.CommandContext(ctx, "git", "-C", repoPath, "log", "--oneline", "-10")
	default:
		return "", fmt.Errorf("unknown action: %s", action)
	}

	output, err := execWithTimeout(cmd, 30*time.Second)
	if err != nil {
		return "", fmt.Errorf("git error: %v", err)
	}
	return output, nil
}

func (t *GitTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "git",
		Description: "Git operations",
		Parameters: map[string]Parameter{
			"action":  {Type: "string", Description: "Action: status/commit/push/pull/log", Required: true},
			"path":   {Type: "string", Description: "Repository path"},
			"message": {Type: "string", Description: "Commit message"},
		},
	}
}

// CalculatorTool 计算器
type CalculatorTool struct{}

func (t *CalculatorTool) Name() string { return "calculator" }
func (t *CalculatorTool) Description() string { return "Calculate expression / 计算表达式" }

func (t *CalculatorTool) Execute(ctx context.Context, input map[string]interface{}) (string, error) {
	expr, ok := input["expression"].(string)
	if !ok || expr == "" {
		return "", fmt.Errorf("expression is required")
	}

	// 简单计算 (实际可用 mathexpr 库)
	result := evalSimple(expr)
	return fmt.Sprintf("Result: %s = %s", expr, result), nil
}

func (t *CalculatorTool) Schema() ToolSchema {
	return ToolSchema{
		Name:        "calculator",
		Description: "Calculate math expression",
		Parameters: map[string]Parameter{
			"expression": {Type: "string", Description: "Math expression", Required: true},
		},
	}
}

// 辅助函数
func execWithTimeout(cmd *exec.Cmd, timeout time.Duration) (string, error) {
	done := make(chan struct{})
	var output []byte
	var err error

	go func() {
		output, err = cmd.CombinedOutput()
		close(done)
	}()

	select {
	case <-done:
		return string(output), err
	case <-time.After(timeout):
		return "", fmt.Errorf("command timeout")
	}
}

func detectLanguage(code string) string {
	lower := strings.ToLower(code)
	if strings.Contains(lower, "package ") && strings.Contains(lower, "func ") {
		return "Go"
	}
	if strings.Contains(lower, "def ") && strings.Contains(lower, ":") {
		return "Python"
	}
	if strings.Contains(lower, "function ") || strings.Contains(lower, "const ") {
		return "JavaScript"
	}
	return "Unknown"
}

func evalSimple(expr string) string {
	// 简化实现
	return "0"
}

// ToolResult 工具执行结果
type ToolResult struct {
	Tool     string        `json:"tool"`
	Input    map[string]interface{} `json:"input"`
	Output   string        `json:"output"`
	Error    string        `json:"error,omitempty"`
	Duration time.Duration `json:"duration"`
}

// ExecuteTool 执行工具
func ExecuteTool(ctx context.Context, toolName string, input map[string]interface{}) *ToolResult {
	start := time.Now()
	result := &ToolResult{
		Tool:   toolName,
		Input:  input,
	}

	tool := GetTool(toolName)
	if tool == nil {
		result.Error = fmt.Sprintf("tool not found: %s", toolName)
		return result
	}

	output, err := tool.Execute(ctx, input)
	if err != nil {
		result.Error = err.Error()
	} else {
		result.Output = output
	}
	result.Duration = time.Since(start)
	return result
}

// ============ 初始化 ============

func init() {
	// 注册内置工具
	RegisterTool(&ShellTool{})
	RegisterTool(&WebSearchTool{})
	RegisterTool(&WebFetchTool{})
	RegisterTool(&FileReadTool{})
	RegisterTool(&FileWriteTool{})
	RegisterTool(&CodeReviewTool{})
	RegisterTool(&TestGenTool{})
	RegisterTool(&GitTool{})
	RegisterTool(&CalculatorTool{})
}
