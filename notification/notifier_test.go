package notification

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNoOpNotifier(t *testing.T) {
	notifier := NewNoOpNotifier()
	assert.NotNil(t, notifier)

	t.Run("发送通知不报错", func(t *testing.T) {
		err := notifier.Notify("Test Title", "Test Message")
		require.NoError(t, err)
	})

	t.Run("发送空标题和消息", func(t *testing.T) {
		err := notifier.Notify("", "")
		require.NoError(t, err)
	})

	t.Run("发送长消息", func(t *testing.T) {
		longMessage := string(make([]byte, 1000))
		err := notifier.Notify("Title", longMessage)
		require.NoError(t, err)
	})
}

// MockNotifier 用于测试的模拟通知器
type MockNotifier struct {
	notifications []Notification
	shouldError   bool
	errorMessage  string
}

// Notification 记录的通知
type Notification struct {
	Title   string
	Message string
}

// NewMockNotifier 创建模拟通知器
func NewMockNotifier() *MockNotifier {
	return &MockNotifier{
		notifications: make([]Notification, 0),
	}
}

// Notify 记录通知
func (m *MockNotifier) Notify(title, message string) error {
	if m.shouldError {
		return assert.AnError
	}
	m.notifications = append(m.notifications, Notification{
		Title:   title,
		Message: message,
	})
	return nil
}

// GetNotifications 获取所有通知
func (m *MockNotifier) GetNotifications() []Notification {
	return m.notifications
}

// SetError 设置是否返回错误
func (m *MockNotifier) SetError(shouldError bool) {
	m.shouldError = shouldError
}

// Clear 清空通知记录
func (m *MockNotifier) Clear() {
	m.notifications = make([]Notification, 0)
}

func TestMockNotifier(t *testing.T) {
	t.Run("记录通知", func(t *testing.T) {
		mock := NewMockNotifier()

		err := mock.Notify("Title 1", "Message 1")
		require.NoError(t, err)

		err = mock.Notify("Title 2", "Message 2")
		require.NoError(t, err)

		notifications := mock.GetNotifications()
		assert.Len(t, notifications, 2)
		assert.Equal(t, "Title 1", notifications[0].Title)
		assert.Equal(t, "Message 1", notifications[0].Message)
		assert.Equal(t, "Title 2", notifications[1].Title)
		assert.Equal(t, "Message 2", notifications[1].Message)
	})

	t.Run("模拟错误", func(t *testing.T) {
		mock := NewMockNotifier()
		mock.SetError(true)

		err := mock.Notify("Title", "Message")
		require.Error(t, err)

		// 错误时不应记录通知
		assert.Empty(t, mock.GetNotifications())
	})

	t.Run("清空通知", func(t *testing.T) {
		mock := NewMockNotifier()
		mock.Notify("Title", "Message")
		assert.Len(t, mock.GetNotifications(), 1)

		mock.Clear()
		assert.Empty(t, mock.GetNotifications())
	})
}
