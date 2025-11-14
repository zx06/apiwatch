package notification

import (
	"context"

	"github.com/gen2brain/beeep"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

// WailsNotifier Wails通知器实现
type WailsNotifier struct {
	ctx context.Context
}

// NewWailsNotifier 创建Wails通知器
func NewWailsNotifier(ctx context.Context) *WailsNotifier {
	return &WailsNotifier{
		ctx: ctx,
	}
}

// Notify 实现Notifier接口，发送系统通知
func (n *WailsNotifier) Notify(title, message string) error {
	if n.ctx == nil {
		return nil // 如果context为nil，静默失败
	}

	// 限制消息长度
	if len(message) > 200 {
		message = message[:197] + "..."
	}

	// 使用Wails runtime发送消息对话框
	_, err := runtime.MessageDialog(n.ctx, runtime.MessageDialogOptions{
		Type:    runtime.InfoDialog,
		Title:   title,
		Message: message,
	})

	// 用beeep发送通知
	err = beeep.Notify(title, message, "")

	return err
}
