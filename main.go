package main

import (
	"context"
	"embed"
	"log/slog"
	"os"
	"time"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/zx06/apiwatch/config"
	"github.com/zx06/apiwatch/core"
	"github.com/zx06/apiwatch/fetcher"
	"github.com/zx06/apiwatch/logger"
	"github.com/zx06/apiwatch/models"
	"github.com/zx06/apiwatch/monitor"
	"github.com/zx06/apiwatch/notification"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// 初始化日志系统
	if err := logger.Setup(logger.DefaultConfig()); err != nil {
		panic("初始化日志系统失败: " + err.Error())
	}

	slog.Info("API Watch 启动中...")

	// 创建配置管理器
	configMgr, err := config.NewYAMLManager("")
	if err != nil {
		slog.Error("创建配置管理器失败", "error", err)
		os.Exit(1)
	}

	// 创建HTTP客户端
	httpFetcher := fetcher.NewHTTPFetcher()

	// 创建临时通知器（将在startup中替换为Wails通知器）
	notifier := notification.NewNoOpNotifier()

	// 创建核心引擎（先创建，以便设置回调）
	var engine *core.Engine

	// 创建规则更新回调函数
	onRuleUpdate := func(rule *models.MonitorRule) {
		if engine != nil {
			// 更新Engine中的规则
			engine.UpdateRuleInMemory(rule)

			// 发布规则状态变化事件
			engine.PublishEvent(core.Event{
				Type:      core.EventRuleStatusChanged,
				RuleID:    rule.ID,
				Rule:      rule,
				Timestamp: time.Now(),
			})
		}
	}

	// 创建监控服务
	monitorSvc := monitor.NewMonitorService(httpFetcher, notifier, onRuleUpdate)

	// 初始化引擎
	engine = core.NewEngine(configMgr, monitorSvc, notifier)

	// 初始化引擎
	if err := engine.Initialize(); err != nil {
		slog.Error("初始化引擎失败", "error", err)
		os.Exit(1)
	}

	// 创建Wails应用
	app := NewApp(engine)

	// 创建事件监听器
	eventListener := &EventListener{}

	// 订阅核心层事件
	engine.Subscribe(eventListener)

	// 运行Wails应用
	err = wails.Run(&options.App{
		Title:  "API Watch - URL内容监控器",
		Width:  1200,
		Height: 900,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup: func(ctx context.Context) {
			app.startup(ctx)
			eventListener.ctx = ctx

			// 创建Wails通知器并更新引擎
			wailsNotifier := notification.NewWailsNotifier(ctx)
			monitorSvc.UpdateNotifier(wailsNotifier)

			slog.Info("API Watch 已启动")
		},
		OnShutdown: app.shutdown,
		OnDomReady: app.domReady,
		Bind: []interface{}{
			app,
		},
	})

	if err != nil {
		slog.Error("启动应用失败", "error", err)
		os.Exit(1)
	}
}
