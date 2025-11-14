package monitor

import (
	"fmt"
	"log/slog"
	"sync"

	"github.com/zx06/apiwatch/extractor"
	"github.com/zx06/apiwatch/fetcher"
	"github.com/zx06/apiwatch/models"
	"github.com/zx06/apiwatch/notification"
)

// Service 监控服务接口
type Service interface {
	// StartTask 启动监控任务
	StartTask(rule *models.MonitorRule) error

	// StopTask 停止监控任务
	StopTask(ruleID string) error

	// StopAll 停止所有任务
	StopAll()

	// UpdateTask 更新任务配置
	UpdateTask(rule *models.MonitorRule) error

	// RunTaskOnce 手动执行一次检查
	RunTaskOnce(ruleID string) error

	// GetTaskStatus 获取任务状态
	GetTaskStatus(ruleID string) models.RuleStatus

	// IsTaskRunning 检查任务是否在运行
	IsTaskRunning(ruleID string) bool
}

// MonitorService 监控服务实现
type MonitorService struct {
	tasks            map[string]*Task
	fetcher          fetcher.Fetcher
	extractorFactory *extractor.Factory
	notifier         notification.Notifier
	mu               sync.RWMutex

	// 回调函数，用于通知规则更新
	onRuleUpdate func(*models.MonitorRule)
}

// NewMonitorService 创建监控服务
func NewMonitorService(
	fetcher fetcher.Fetcher,
	notifier notification.Notifier,
	onRuleUpdate func(*models.MonitorRule),
) *MonitorService {
	return &MonitorService{
		tasks:            make(map[string]*Task),
		fetcher:          fetcher,
		extractorFactory: extractor.NewFactory(),
		notifier:         notifier,
		onRuleUpdate:     onRuleUpdate,
	}
}

// StartTask 启动监控任务
func (s *MonitorService) StartTask(rule *models.MonitorRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	// 检查任务是否已存在
	if task, exists := s.tasks[rule.ID]; exists {
		if task.IsRunning() {
			return fmt.Errorf("任务已在运行中: %s", rule.ID)
		}
		// 停止旧任务
		task.Stop()
		delete(s.tasks, rule.ID)
	}

	// 创建新任务
	task, err := NewTask(rule, s.fetcher, s.extractorFactory, s.notifier, s.onRuleUpdate)
	if err != nil {
		return fmt.Errorf("创建任务失败: %w", err)
	}

	// 启动任务
	if err := task.Start(); err != nil {
		return fmt.Errorf("启动任务失败: %w", err)
	}

	s.tasks[rule.ID] = task

	slog.Info("监控服务已启动任务",
		"rule_id", rule.ID,
		"rule_name", rule.Name,
	)

	return nil
}

// StopTask 停止监控任务
func (s *MonitorService) StopTask(ruleID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[ruleID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", ruleID)
	}

	task.Stop()
	delete(s.tasks, ruleID)

	slog.Info("监控服务已停止任务", "rule_id", ruleID)

	return nil
}

// StopAll 停止所有任务
func (s *MonitorService) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()

	for id, task := range s.tasks {
		task.Stop()
		delete(s.tasks, id)
	}

	slog.Info("监控服务已停止所有任务", "count", len(s.tasks))
}

// UpdateTask 更新任务配置
func (s *MonitorService) UpdateTask(rule *models.MonitorRule) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	task, exists := s.tasks[rule.ID]
	if !exists {
		return fmt.Errorf("任务不存在: %s", rule.ID)
	}

	// 如果任务正在运行，需要重启
	wasRunning := task.IsRunning()
	if wasRunning {
		task.Stop()
	}

	// 更新任务配置
	if err := task.Update(rule); err != nil {
		return fmt.Errorf("更新任务失败: %w", err)
	}

	// 如果之前在运行，重新启动
	if wasRunning {
		if err := task.Start(); err != nil {
			return fmt.Errorf("重启任务失败: %w", err)
		}
	}

	slog.Info("监控服务已更新任务",
		"rule_id", rule.ID,
		"rule_name", rule.Name,
	)

	return nil
}

// RunTaskOnce 手动执行一次检查
func (s *MonitorService) RunTaskOnce(ruleID string) error {
	s.mu.RLock()
	task, exists := s.tasks[ruleID]
	s.mu.RUnlock()

	if !exists {
		return fmt.Errorf("任务不存在或未启动: %s，请先启动监控", ruleID)
	}

	slog.Info("手动执行任务检查", "rule_id", ruleID)

	return task.RunOnce()
}

// GetTaskStatus 获取任务状态
func (s *MonitorService) GetTaskStatus(ruleID string) models.RuleStatus {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[ruleID]
	if !exists {
		return models.StatusIdle
	}

	rule := task.GetRule()
	return rule.Status
}

// IsTaskRunning 检查任务是否在运行
func (s *MonitorService) IsTaskRunning(ruleID string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, exists := s.tasks[ruleID]
	if !exists {
		return false
	}

	return task.IsRunning()
}

// GetActiveTasks 获取所有活跃任务数量
func (s *MonitorService) GetActiveTasks() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.tasks)
}

// UpdateNotifier 更新通知器（用于Wails启动后替换通知器）
func (s *MonitorService) UpdateNotifier(notifier notification.Notifier) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.notifier = notifier

	// 更新所有活跃任务的通知器
	for _, task := range s.tasks {
		task.UpdateNotifier(notifier)
	}

	slog.Info("监控服务已更新通知器")
}
