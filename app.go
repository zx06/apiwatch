package main

import (
	"context"

	"github.com/wailsapp/wails/v2/pkg/runtime"
	"github.com/zx06/apiwatch/core"
	"github.com/zx06/apiwatch/models"
)

// App Wails应用结构
type App struct {
	ctx     context.Context
	coreAPI core.CoreAPI
}

// NewApp 创建应用实例
func NewApp(coreAPI core.CoreAPI) *App {
	return &App{
		coreAPI: coreAPI,
	}
}

// startup 应用启动时调用
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// shutdown 应用关闭时调用
func (a *App) shutdown(ctx context.Context) {
	// 优雅关闭核心引擎
	if a.coreAPI != nil {
		a.coreAPI.Shutdown()
	}
}

// domReady DOM加载完成时调用
func (a *App) domReady(ctx context.Context) {
	// 可以在这里执行DOM就绪后的操作
}

// GetRules 获取所有规则
func (a *App) GetRules() ([]*models.MonitorRule, error) {
	return a.coreAPI.GetRules()
}

// GetRule 获取单个规则
func (a *App) GetRule(id string) (*models.MonitorRule, error) {
	return a.coreAPI.GetRule(id)
}

// AddRule 添加规则
func (a *App) AddRule(rule *models.MonitorRule) error {
	return a.coreAPI.AddRule(rule)
}

// UpdateRule 更新规则
func (a *App) UpdateRule(rule *models.MonitorRule) error {
	return a.coreAPI.UpdateRule(rule)
}

// DeleteRule 删除规则
func (a *App) DeleteRule(id string) error {
	return a.coreAPI.DeleteRule(id)
}

// StartMonitoring 启动监控
func (a *App) StartMonitoring(ruleID string) error {
	return a.coreAPI.StartMonitoring(ruleID)
}

// StopMonitoring 停止监控
func (a *App) StopMonitoring(ruleID string) error {
	return a.coreAPI.StopMonitoring(ruleID)
}

// CheckNow 立即检查
func (a *App) CheckNow(ruleID string) error {
	return a.coreAPI.CheckNow(ruleID)
}

// EventListener 实现事件监听器接口
type EventListener struct {
	ctx context.Context
}

// OnEvent 处理核心层事件
func (e *EventListener) OnEvent(event core.Event) {
	// 只有在context可用时才推送事件
	if e.ctx != nil {
		runtime.EventsEmit(e.ctx, string(event.Type), event)
	}
}
