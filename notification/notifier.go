package notification

// Notifier 通知接口（抽象，不依赖具体UI）
type Notifier interface {
	// Notify 发送通知
	Notify(title, message string) error
}

// NoOpNotifier 空实现（用于测试或无通知场景）
type NoOpNotifier struct{}

// NewNoOpNotifier 创建空通知器
func NewNoOpNotifier() *NoOpNotifier {
	return &NoOpNotifier{}
}

// Notify 空实现，不做任何事
func (n *NoOpNotifier) Notify(title, message string) error {
	// 不发送通知
	return nil
}
