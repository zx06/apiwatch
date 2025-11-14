package monitor

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/zx06/apiwatch/extractor"
	"github.com/zx06/apiwatch/fetcher"
	"github.com/zx06/apiwatch/models"
	"github.com/zx06/apiwatch/notification"
)

// Task 监控任务
type Task struct {
	rule      *models.MonitorRule
	fetcher   fetcher.Fetcher
	extractor extractor.Extractor
	notifier  notification.Notifier

	ticker  *time.Ticker
	stopCh  chan struct{}
	mu      sync.RWMutex
	running bool

	// 回调函数，用于通知状态变化
	onUpdate func(*models.MonitorRule)
}

// NewTask 创建监控任务
func NewTask(
	rule *models.MonitorRule,
	fetcher fetcher.Fetcher,
	extractorFactory *extractor.Factory,
	notifier notification.Notifier,
	onUpdate func(*models.MonitorRule),
) (*Task, error) {
	// 创建提取器
	ext, err := extractorFactory.Create(rule.ExtractorType, rule.ExtractorExpr)
	if err != nil {
		return nil, fmt.Errorf("创建提取器失败: %w", err)
	}

	return &Task{
		rule:      rule,
		fetcher:   fetcher,
		extractor: ext,
		notifier:  notifier,
		stopCh:    make(chan struct{}),
		onUpdate:  onUpdate,
	}, nil
}

// Start 启动监控任务
func (t *Task) Start() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	if t.running {
		return fmt.Errorf("任务已在运行中")
	}

	t.running = true
	t.ticker = time.NewTicker(time.Duration(t.rule.Interval))

	// 启动监控goroutine
	go t.run()

	slog.Info("监控任务已启动",
		"rule_id", t.rule.ID,
		"rule_name", t.rule.Name,
		"interval", t.rule.Interval,
	)

	return nil
}

// Stop 停止监控任务
func (t *Task) Stop() {
	t.mu.Lock()
	defer t.mu.Unlock()

	if !t.running {
		return
	}

	t.running = false
	close(t.stopCh)

	if t.ticker != nil {
		t.ticker.Stop()
	}

	slog.Info("监控任务已停止",
		"rule_id", t.rule.ID,
		"rule_name", t.rule.Name,
	)
}

// run 监控循环
func (t *Task) run() {
	// 立即执行一次检查
	if err := t.RunOnce(); err != nil {
		slog.Error("首次检查失败",
			"rule_id", t.rule.ID,
			"error", err,
		)
	}

	// 定期检查
	for {
		select {
		case <-t.stopCh:
			return
		case <-t.ticker.C:
			if err := t.RunOnce(); err != nil {
				slog.Error("定期检查失败",
					"rule_id", t.rule.ID,
					"error", err,
				)
			}
		}
	}
}

// RunOnce 执行一次检查
func (t *Task) RunOnce() error {
	t.mu.Lock()
	defer t.mu.Unlock()

	slog.Debug("开始检查",
		"rule_id", t.rule.ID,
		"rule_name", t.rule.Name,
		"url", t.rule.URL,
	)

	// 更新状态为运行中
	t.rule.Status = models.StatusRunning
	t.rule.ErrorMessage = ""
	t.notifyUpdate()

	// 发送HTTP请求
	req := &fetcher.Request{
		URL:     t.rule.URL,
		Method:  t.rule.Method,
		Headers: t.rule.Headers,
		Body:    t.rule.Body,
	}

	resp, err := t.fetcher.Fetch(req)
	if err != nil {
		t.handleError(fmt.Errorf("HTTP请求失败: %w", err))
		return err
	}

	// 提取内容
	content, err := t.extractor.Extract(resp.Body, resp.ContentType)
	if err != nil {
		t.handleError(fmt.Errorf("内容提取失败: %w", err))
		return err
	}

	// 更新最后检查时间
	t.rule.LastChecked = time.Now()

	// 检测内容变化
	if t.rule.LastContent != "" && t.rule.LastContent != content {
		// 内容发生变化
		slog.Info("检测到内容变化",
			"rule_id", t.rule.ID,
			"rule_name", t.rule.Name,
		)

		// 发送通知
		if t.rule.NotifyEnabled {
			if err := t.sendNotification(content); err != nil {
				slog.Warn("发送通知失败",
					"rule_id", t.rule.ID,
					"error", err,
				)
			}
		}
	}

	// 更新内容
	t.rule.LastContent = content
	t.rule.Status = models.StatusRunning
	t.notifyUpdate()

	slog.Debug("检查完成",
		"rule_id", t.rule.ID,
		"content_length", len(content),
	)

	return nil
}

// Update 更新任务配置
func (t *Task) Update(rule *models.MonitorRule) error {
	t.mu.Lock()
	defer t.mu.Unlock()

	// 检查是否需要重新创建提取器
	if rule.ExtractorType != t.rule.ExtractorType || rule.ExtractorExpr != t.rule.ExtractorExpr {
		factory := extractor.NewFactory()
		ext, err := factory.Create(rule.ExtractorType, rule.ExtractorExpr)
		if err != nil {
			return fmt.Errorf("创建提取器失败: %w", err)
		}
		t.extractor = ext
	}

	// 检查是否需要重新创建ticker
	if rule.Interval != t.rule.Interval && t.running {
		if t.ticker != nil {
			t.ticker.Stop()
		}
		t.ticker = time.NewTicker(time.Duration(rule.Interval))
	}

	// 更新规则
	t.rule = rule

	slog.Info("任务配置已更新",
		"rule_id", t.rule.ID,
		"rule_name", t.rule.Name,
	)

	return nil
}

// IsRunning 检查任务是否在运行
func (t *Task) IsRunning() bool {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.running
}

// GetRule 获取规则
func (t *Task) GetRule() *models.MonitorRule {
	t.mu.RLock()
	defer t.mu.RUnlock()
	return t.rule
}

// handleError 处理错误
func (t *Task) handleError(err error) {
	t.rule.Status = models.StatusError
	t.rule.ErrorMessage = err.Error()
	t.notifyUpdate()

	slog.Error("任务执行错误",
		"rule_id", t.rule.ID,
		"rule_name", t.rule.Name,
		"error", err,
	)
}

// sendNotification 发送通知
func (t *Task) sendNotification(newContent string) error {
	title := fmt.Sprintf("内容变化: %s", t.rule.Name)

	// 限制消息长度
	message := newContent
	if len(message) > 200 {
		message = message[:200] + "..."
	}

	if t.rule.Description != "" {
		message = t.rule.Description + "\n\n" + message
	}

	return t.notifier.Notify(title, message)
}

// notifyUpdate 通知规则更新
func (t *Task) notifyUpdate() {
	if t.onUpdate != nil {
		t.onUpdate(t.rule)
	}
}

// UpdateNotifier 更新通知器
func (t *Task) UpdateNotifier(notifier notification.Notifier) {
	t.mu.Lock()
	defer t.mu.Unlock()
	t.notifier = notifier
}
