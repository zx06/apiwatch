package core

import (
	"sync"
	"time"

	"github.com/zx06/apiwatch/models"
)

// EventType 事件类型
type EventType string

const (
	EventRuleAdded         EventType = "rule_added"
	EventRuleUpdated       EventType = "rule_updated"
	EventRuleDeleted       EventType = "rule_deleted"
	EventRuleStatusChanged EventType = "rule_status_changed"
	EventContentChanged    EventType = "content_changed"
	EventMonitorError      EventType = "monitor_error"
)

// Event 事件
type Event struct {
	Type      EventType
	RuleID    string
	Rule      *models.MonitorRule
	Timestamp time.Time
	Data      interface{}
}

// EventListener 事件监听器接口
type EventListener interface {
	OnEvent(event Event)
}

// EventBus 事件总线接口
type EventBus interface {
	Publish(event Event)
	Subscribe(listener EventListener)
	Unsubscribe(listener EventListener)
}

// DefaultEventBus 默认事件总线实现
type DefaultEventBus struct {
	listeners []EventListener
	mu        sync.RWMutex
}

// NewEventBus 创建事件总线
func NewEventBus() *DefaultEventBus {
	return &DefaultEventBus{
		listeners: make([]EventListener, 0),
	}
}

// Publish 发布事件
func (b *DefaultEventBus) Publish(event Event) {
	b.mu.RLock()
	listeners := make([]EventListener, len(b.listeners))
	copy(listeners, b.listeners)
	b.mu.RUnlock()

	// 异步分发事件，避免阻塞
	for _, listener := range listeners {
		go func(l EventListener) {
			defer func() {
				if r := recover(); r != nil {
					// 捕获监听器中的panic，避免影响其他监听器
				}
			}()
			l.OnEvent(event)
		}(listener)
	}
}

// Subscribe 订阅事件
func (b *DefaultEventBus) Subscribe(listener EventListener) {
	b.mu.Lock()
	defer b.mu.Unlock()
	b.listeners = append(b.listeners, listener)
}

// Unsubscribe 取消订阅
func (b *DefaultEventBus) Unsubscribe(listener EventListener) {
	b.mu.Lock()
	defer b.mu.Unlock()

	for i, l := range b.listeners {
		// 使用指针比较
		if l == listener {
			b.listeners = append(b.listeners[:i], b.listeners[i+1:]...)
			return
		}
	}
}
