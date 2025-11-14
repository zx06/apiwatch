package core

import (
	"github.com/zx06/apiwatch/models"
)

// CoreAPI 核心API接口，UI层通过此接口与核心层交互
type CoreAPI interface {
	// 规则管理
	GetRules() ([]*models.MonitorRule, error)
	GetRule(id string) (*models.MonitorRule, error)
	AddRule(rule *models.MonitorRule) error
	UpdateRule(rule *models.MonitorRule) error
	DeleteRule(id string) error

	// 监控控制
	StartMonitoring(ruleID string) error
	StopMonitoring(ruleID string) error
	StopAllMonitoring() error
	CheckNow(ruleID string) error

	// 事件订阅
	Subscribe(listener EventListener)
	Unsubscribe(listener EventListener)

	// 生命周期
	Initialize() error
	Shutdown() error
}
