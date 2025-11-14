package core

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/zx06/apiwatch/config"
	"github.com/zx06/apiwatch/models"
	"github.com/zx06/apiwatch/monitor"
	"github.com/zx06/apiwatch/notification"
)

// Engine Monitor引擎，实现CoreAPI接口
type Engine struct {
	configMgr  config.Manager
	monitorSvc monitor.Service
	eventBus   EventBus
	notifier   notification.Notifier

	rules []*models.MonitorRule
	mu    sync.RWMutex
}

// NewEngine 创建Monitor引擎
func NewEngine(
	configMgr config.Manager,
	monitorSvc monitor.Service,
	notifier notification.Notifier,
) *Engine {
	return &Engine{
		configMgr:  configMgr,
		monitorSvc: monitorSvc,
		eventBus:   NewEventBus(),
		notifier:   notifier,
		rules:      make([]*models.MonitorRule, 0),
	}
}

// Initialize 初始化引擎
func (e *Engine) Initialize() error {
	slog.Info("初始化Monitor引擎")

	// 加载配置
	rules, err := e.configMgr.Load()
	if err != nil {
		return fmt.Errorf("加载配置失败: %w", err)
	}

	e.mu.Lock()
	e.rules = rules
	e.mu.Unlock()

	// 启动已启用的监控任务
	for _, rule := range rules {
		if rule.Enabled {
			if err := e.monitorSvc.StartTask(rule); err != nil {
				slog.Error("启动监控任务失败",
					"rule_id", rule.ID,
					"error", err,
				)
			}
		}
	}

	slog.Info("Monitor引擎初始化完成", "rules_count", len(rules))
	return nil
}

// Shutdown 关闭引擎
func (e *Engine) Shutdown() error {
	slog.Info("关闭Monitor引擎")

	// 停止所有监控任务
	e.monitorSvc.StopAll()

	// 保存配置
	e.mu.RLock()
	rules := e.rules
	e.mu.RUnlock()

	if err := e.configMgr.Save(rules); err != nil {
		slog.Error("保存配置失败", "error", err)
		return err
	}

	slog.Info("Monitor引擎已关闭")
	return nil
}

// GetRules 获取所有规则
func (e *Engine) GetRules() ([]*models.MonitorRule, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	// 返回副本
	rules := make([]*models.MonitorRule, len(e.rules))
	copy(rules, e.rules)
	return rules, nil
}

// GetRule 获取单个规则
func (e *Engine) GetRule(id string) (*models.MonitorRule, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()

	for _, rule := range e.rules {
		if rule.ID == id {
			return rule, nil
		}
	}

	return nil, fmt.Errorf("规则不存在: %s", id)
}

// AddRule 添加规则
func (e *Engine) AddRule(rule *models.MonitorRule) error {
	// 验证规则
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("规则验证失败: %w", err)
	}

	// 生成ID
	if rule.ID == "" {
		rule.ID = uuid.New().String()
	}

	// 设置初始状态
	if rule.Status == "" {
		rule.Status = models.StatusIdle
	}

	e.mu.Lock()
	e.rules = append(e.rules, rule)
	e.mu.Unlock()

	// 保存到配置
	if err := e.configMgr.AddRule(rule); err != nil {
		// 回滚
		e.mu.Lock()
		for i, r := range e.rules {
			if r.ID == rule.ID {
				e.rules = append(e.rules[:i], e.rules[i+1:]...)
				break
			}
		}
		e.mu.Unlock()
		return fmt.Errorf("保存规则失败: %w", err)
	}

	// 如果启用，启动监控
	if rule.Enabled {
		if err := e.monitorSvc.StartTask(rule); err != nil {
			slog.Error("启动监控任务失败", "rule_id", rule.ID, "error", err)
		}
	}

	// 发布事件
	e.eventBus.Publish(Event{
		Type:      EventRuleAdded,
		RuleID:    rule.ID,
		Rule:      rule,
		Timestamp: time.Now(),
	})

	slog.Info("规则已添加", "rule_id", rule.ID, "rule_name", rule.Name)
	return nil
}

// UpdateRule 更新规则
func (e *Engine) UpdateRule(rule *models.MonitorRule) error {
	// 验证规则
	if err := rule.Validate(); err != nil {
		return fmt.Errorf("规则验证失败: %w", err)
	}

	e.mu.Lock()
	found := false
	for i, r := range e.rules {
		if r.ID == rule.ID {
			e.rules[i] = rule
			found = true
			break
		}
	}
	e.mu.Unlock()

	if !found {
		return fmt.Errorf("规则不存在: %s", rule.ID)
	}

	// 保存到配置
	if err := e.configMgr.UpdateRule(rule); err != nil {
		return fmt.Errorf("保存规则失败: %w", err)
	}

	// 更新监控任务
	if rule.Enabled {
		if e.monitorSvc.IsTaskRunning(rule.ID) {
			if err := e.monitorSvc.UpdateTask(rule); err != nil {
				slog.Error("更新监控任务失败", "rule_id", rule.ID, "error", err)
			}
		} else {
			if err := e.monitorSvc.StartTask(rule); err != nil {
				slog.Error("启动监控任务失败", "rule_id", rule.ID, "error", err)
			}
		}
	} else {
		if e.monitorSvc.IsTaskRunning(rule.ID) {
			if err := e.monitorSvc.StopTask(rule.ID); err != nil {
				slog.Error("停止监控任务失败", "rule_id", rule.ID, "error", err)
			}
		}
	}

	// 发布事件
	e.eventBus.Publish(Event{
		Type:      EventRuleUpdated,
		RuleID:    rule.ID,
		Rule:      rule,
		Timestamp: time.Now(),
	})

	slog.Info("规则已更新", "rule_id", rule.ID, "rule_name", rule.Name)
	return nil
}

// DeleteRule 删除规则
func (e *Engine) DeleteRule(id string) error {
	e.mu.Lock()
	found := false
	for i, r := range e.rules {
		if r.ID == id {
			e.rules = append(e.rules[:i], e.rules[i+1:]...)
			found = true
			break
		}
	}
	e.mu.Unlock()

	if !found {
		return fmt.Errorf("规则不存在: %s", id)
	}

	// 停止监控任务
	if e.monitorSvc.IsTaskRunning(id) {
		if err := e.monitorSvc.StopTask(id); err != nil {
			slog.Error("停止监控任务失败", "rule_id", id, "error", err)
		}
	}

	// 从配置中删除
	if err := e.configMgr.DeleteRule(id); err != nil {
		return fmt.Errorf("删除规则失败: %w", err)
	}

	// 发布事件
	e.eventBus.Publish(Event{
		Type:      EventRuleDeleted,
		RuleID:    id,
		Timestamp: time.Now(),
	})

	slog.Info("规则已删除", "rule_id", id)
	return nil
}

// StartMonitoring 启动监控
func (e *Engine) StartMonitoring(ruleID string) error {
	rule, err := e.GetRule(ruleID)
	if err != nil {
		return err
	}

	if err := e.monitorSvc.StartTask(rule); err != nil {
		return fmt.Errorf("启动监控失败: %w", err)
	}

	// 更新规则状态
	rule.Enabled = true
	return e.UpdateRule(rule)
}

// StopMonitoring 停止监控
func (e *Engine) StopMonitoring(ruleID string) error {
	if err := e.monitorSvc.StopTask(ruleID); err != nil {
		return fmt.Errorf("停止监控失败: %w", err)
	}

	// 更新规则状态
	rule, err := e.GetRule(ruleID)
	if err != nil {
		return err
	}

	rule.Enabled = false
	rule.Status = models.StatusPaused
	return e.UpdateRule(rule)
}

// StopAllMonitoring 停止所有监控
func (e *Engine) StopAllMonitoring() error {
	e.monitorSvc.StopAll()

	e.mu.Lock()
	for _, rule := range e.rules {
		rule.Enabled = false
		rule.Status = models.StatusPaused
	}
	e.mu.Unlock()

	slog.Info("已停止所有监控")
	return nil
}

// CheckNow 立即检查
func (e *Engine) CheckNow(ruleID string) error {
	return e.monitorSvc.RunTaskOnce(ruleID)
}

// Subscribe 订阅事件
func (e *Engine) Subscribe(listener EventListener) {
	e.eventBus.Subscribe(listener)
}

// Unsubscribe 取消订阅
func (e *Engine) Unsubscribe(listener EventListener) {
	e.eventBus.Unsubscribe(listener)
}

// PublishEvent 发布事件（供外部调用）
func (e *Engine) PublishEvent(event Event) {
	e.eventBus.Publish(event)
}

// UpdateRuleInMemory 更新内存中的规则（不保存到配置文件）
func (e *Engine) UpdateRuleInMemory(rule *models.MonitorRule) {
	e.mu.Lock()
	defer e.mu.Unlock()

	// 查找并更新规则
	for i, r := range e.rules {
		if r.ID == rule.ID {
			e.rules[i] = rule
			slog.Debug("已更新内存中的规则",
				"rule_id", rule.ID,
				"status", rule.Status,
				"last_checked", rule.LastChecked,
			)
			return
		}
	}
}

// handleRuleUpdate 处理规则更新（由监控任务回调）
func (e *Engine) handleRuleUpdate(rule *models.MonitorRule) {
	e.mu.Lock()
	for i, r := range e.rules {
		if r.ID == rule.ID {
			e.rules[i] = rule
			break
		}
	}
	e.mu.Unlock()

	// 保存到配置
	if err := e.configMgr.UpdateRule(rule); err != nil {
		slog.Error("保存规则更新失败", "rule_id", rule.ID, "error", err)
	}

	// 发布事件
	e.eventBus.Publish(Event{
		Type:      EventRuleStatusChanged,
		RuleID:    rule.ID,
		Rule:      rule,
		Timestamp: time.Now(),
	})
}
