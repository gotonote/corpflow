// CorpFlow Execution Logs - Track workflow runs
// Similar to Claude Code / OpenCode execution tracking

package logs

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

// LogLevel 日志级别
type LogLevel string

const (
	LevelDebug LogLevel = "DEBUG"
	LevelInfo  LogLevel = "INFO"
	LevelWarn  LogLevel = "WARN"
	LevelError LogLevel = "ERROR"
)

// ExecutionLog 执行日志
type ExecutionLog struct {
	ID          string      `json:"id"`           // 日志ID
	FlowID      string      `json:"flow_id"`      // 流程ID
	FlowName    string      `json:"flow_name"`    // 流程名称
	StartedAt   time.Time   `json:"started_at"`   // 开始时间
	EndedAt     *time.Time  `json:"ended_at"`     // 结束时间
	Status      string      `json:"status"`       // running/success/failed
	Trigger     string      `json:"trigger"`       // 触发方式
	Input       interface{} `json:"input"`         // 输入数据
	Output      interface{} `json:"output"`        // 输出数据
	Steps       []StepLog   `json:"steps"`         // 步骤日志
	Metadata    map[string]interface{} `json:"metadata"`
}

// StepLog 步骤日志
type StepLog struct {
	ID          string                 `json:"id"`           // 步骤ID
	NodeID      string                 `json:"node_id"`      // 节点ID
	NodeName    string                 `json:"node_name"`     // 节点名称
	NodeType    string                 `json:"node_type"`     // 节点类型
	StartedAt   time.Time              `json:"started_at"`   // 开始时间
	EndedAt     *time.Time             `json:"ended_at"`     // 结束时间
	Status      string                 `json:"status"`       // pending/running/success/failed/skipped
	Input       map[string]interface{} `json:"input"`        // 输入
	Output      map[string]interface{}  `json:"output"`       // 输出
	Error       string                 `json:"error"`         // 错误信息
	Duration    time.Duration          `json:"duration"`     // 耗时
}

// Service 日志服务
type Service struct {
	mu       sync.RWMutex
	logs     map[string]*ExecutionLog
	storage  string // 存储路径
}

// NewService 创建日志服务
func NewService(storagePath string) *Service {
	if storagePath == "" {
		storagePath = "./data/logs"
	}
	
	s := &Service{
		logs:    make(map[string]*ExecutionLog),
		storage: storagePath,
	}
	
	// 创建存储目录
	os.MkdirAll(storagePath, 0755)
	
	// 加载历史日志
	s.loadLogs()
	
	return s
}

// StartExecution 开始执行日志
func (s *Service) StartExecution(flowID, flowName, trigger string, input interface{}) *ExecutionLog {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log := &ExecutionLog{
		ID:        generateID(),
		FlowID:    flowID,
		FlowName:  flowName,
		StartedAt: time.Now(),
		Status:    "running",
		Trigger:   trigger,
		Input:     input,
		Steps:     []StepLog{},
		Metadata:  make(map[string]interface{}),
	}
	
	s.logs[log.ID] = log
	go s.persistLog(log)
	
	return log
}

// EndExecution 结束执行日志
func (s *Service) EndExecution(logID string, status string, output interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log, ok := s.logs[logID]
	if !ok {
		return fmt.Errorf("log not found: %s", logID)
	}
	
	now := time.Now()
	log.EndedAt = &now
	log.Status = status
	log.Output = output
	
	go s.persistLog(log)
	
	return nil
}

// AddStep 添加步骤日志
func (s *Service) AddStep(logID string, stepID, nodeID, nodeName, nodeType string, input map[string]interface{}) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log, ok := s.logs[logID]
	if !ok {
		return
	}
	
	step := StepLog{
		ID:        stepID,
		NodeID:    nodeID,
		NodeName:  nodeName,
		NodeType:  nodeType,
		StartedAt: time.Now(),
		Status:    "running",
		Input:     input,
	}
	
	log.Steps = append(log.Steps, step)
}

// EndStep 结束步骤日志
func (s *Service) EndStep(logID, stepID, status string, output map[string]interface{}, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	
	log, ok := s.logs[logID]
	if !ok {
		return
	}
	
	for i := range log.Steps {
		if log.Steps[i].ID == stepID {
			now := time.Now()
			log.Steps[i].EndedAt = &now
			log.Steps[i].Status = status
			log.Steps[i].Output = output
			log.Steps[i].Duration = now.Sub(log.Steps[i].StartedAt)
			
			if err != nil {
				log.Steps[i].Error = err.Error()
			}
			break
		}
	}
}

// GetLog 获取日志
func (s *Service) GetLog(logID string) (*ExecutionLog, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	log, ok := s.logs[logID]
	if !ok {
		return nil, fmt.Errorf("log not found")
	}
	
	return log, nil
}

// ListLogs 列出日志
func (s *Service) ListLogs(flowID string, limit int) []*ExecutionLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	var result []*ExecutionLog
	count := 0
	
	// 按时间倒序
	for i := len(s.logs) - 1; i >= 0; i-- {
		if limit > 0 && count >= limit {
			break
		}
		
		log := s.logs[generateID()] // 需要正确遍历
		if flowID == "" || log.FlowID == flowID {
			result = append(result, log)
			count++
		}
	}
	
	return result
}

// GetRecentLogs 获取最近日志
func (s *Service) GetRecentLogs(limit int) []*ExecutionLog {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	logs := make([]*ExecutionLog, 0, len(s.logs))
	for _, log := range s.logs {
		logs = append(logs, log)
	}
	
	// 排序: 最新的在前
	for i := 0; i < len(logs)-1; i++ {
		for j := i + 1; j < len(logs); j++ {
			if logs[j].StartedAt.After(logs[i].StartedAt) {
				logs[i], logs[j] = logs[j], logs[i]
			}
		}
	}
	
	if limit > 0 && len(logs) > limit {
		logs = logs[:limit]
	}
	
	return logs
}

// persistLog 持久化日志
func (s *Service) persistLog(log *ExecutionLog) {
	filename := filepath.Join(s.storage, fmt.Sprintf("%s.json", log.ID))
	data, err := json.MarshalIndent(log, "", "  ")
	if err != nil {
		return
	}
	os.WriteFile(filename, data, 0645)
}

// loadLogs 加载历史日志
func (s *Service) loadLogs() {
	files, err := os.ReadDir(s.storage)
	if err != nil {
		return
	}
	
	for _, f := range files {
		if filepath.Ext(f.Name()) != ".json" {
			continue
		}
		
		data, err := os.ReadFile(filepath.Join(s.storage, f.Name()))
		if err != nil {
			continue
		}
		
		var log ExecutionLog
		if err := json.Unmarshal(data, &log); err == nil {
			s.logs[log.ID] = &log
		}
	}
}

// LogFormat 格式化日志输出
func (log *ExecutionLog) String() string {
	duration := ""
	if log.EndedAt != nil {
		duration = log.EndedAt.Sub(log.StartedAt).String()
	}
	
	output := fmt.Sprintf(`
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━
📊 Execution: %s
🌊 Flow: %s (%s)
⏱️ Duration: %s
📌 Status: %s
🔔 Trigger: %s
`, log.ID[:8], log.FlowName, log.FlowID, duration, log.Status, log.Trigger)

	if len(log.Steps) > 0 {
		output += "\n📝 Steps:\n"
		for _, step := range log.Steps {
			statusIcon := "⏳"
			switch step.Status {
			case "success":
				statusIcon = "✅"
			case "failed":
				statusIcon = "❌"
			case "running":
				statusIcon = "🔄"
			case "skipped":
				statusIcon = "⏭️"
			}
			
			output += fmt.Sprintf("  %s %s (%s) - %s\n", 
				statusIcon, step.NodeName, step.NodeType, step.Duration)
		}
	}
	
	return output
}

// GetExecutionStats 获取执行统计
func (s *Service) GetExecutionStats(days int) map[string]interface{} {
	s.mu.RLock()
	defer s.mu.RUnlock()
	
	stats := map[string]interface{}{
		"total":      len(s.logs),
		"success":    0,
		"failed":     0,
		"running":    0,
		"avg_duration": 0,
	}
	
	var totalDuration time.Duration
	since := time.Now().AddDate(0, 0, -days)
	
	for _, log := range s.logs {
		if log.StartedAt.Before(since) {
			continue
		}
		
		switch log.Status {
		case "success":
			stats["success"] = stats["success"].(int) + 1
		case "failed":
			stats["failed"] = stats["failed"].(int) + 1
		case "running":
			stats["running"] = stats["running"].(int) + 1
		}
		
		if log.EndedAt != nil {
			totalDuration += log.EndedAt.Sub(log.StartedAt)
		}
	}
	
	total := stats["success"].(int) + stats["failed"].(int)
	if total > 0 {
		stats["avg_duration"] = totalDuration.Seconds() / float64(total)
	}
	
	return stats
}

func generateID() string {
	return fmt.Sprintf("%d-%s", time.Now().Unix(), randomString(8))
}

func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyz0123456789"
	b := make([]byte, n)
	for i := range b {
		b[i] = letters[time.Now().UnixNano()%int64(len(letters))]
	}
	return string(b)
}
