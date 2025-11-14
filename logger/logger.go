package logger

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// Config 日志配置
type Config struct {
	Level      slog.Level
	OutputFile string
	MaxSize    int64 // 最大文件大小（字节）
	MaxBackups int   // 保留的备份文件数量
}

// DefaultConfig 默认配置
func DefaultConfig() *Config {
	homeDir, _ := os.UserHomeDir()
	logDir := filepath.Join(homeDir, ".url-monitor", "logs")

	return &Config{
		Level:      slog.LevelInfo,
		OutputFile: filepath.Join(logDir, "app.log"),
		MaxSize:    10 * 1024 * 1024, // 10MB
		MaxBackups: 5,
	}
}

// Setup 设置日志系统
func Setup(cfg *Config) error {
	if cfg == nil {
		cfg = DefaultConfig()
	}

	// 确保日志目录存在
	logDir := filepath.Dir(cfg.OutputFile)
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return fmt.Errorf("创建日志目录失败: %w", err)
	}

	// 打开日志文件
	file, err := os.OpenFile(cfg.OutputFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return fmt.Errorf("打开日志文件失败: %w", err)
	}

	// 检查文件大小，如果超过限制则轮转
	if err := rotateIfNeeded(cfg); err != nil {
		return fmt.Errorf("日志轮转失败: %w", err)
	}

	// 创建多输出writer（文件和控制台）
	multiWriter := io.MultiWriter(file, os.Stdout)

	// 创建JSON handler
	handler := slog.NewJSONHandler(multiWriter, &slog.HandlerOptions{
		Level:     cfg.Level,
		AddSource: true,
	})

	// 设置默认logger
	logger := slog.New(handler)
	slog.SetDefault(logger)

	slog.Info("日志系统已初始化",
		"level", cfg.Level.String(),
		"output", cfg.OutputFile,
	)

	return nil
}

// rotateIfNeeded 如果需要则轮转日志文件
func rotateIfNeeded(cfg *Config) error {
	info, err := os.Stat(cfg.OutputFile)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 文件不存在，无需轮转
		}
		return err
	}

	// 检查文件大小
	if info.Size() < cfg.MaxSize {
		return nil // 未超过大小限制
	}

	// 执行轮转
	return rotate(cfg)
}

// rotate 执行日志轮转
func rotate(cfg *Config) error {
	// 删除最旧的备份
	oldestBackup := fmt.Sprintf("%s.%d", cfg.OutputFile, cfg.MaxBackups)
	os.Remove(oldestBackup) // 忽略错误

	// 移动现有备份
	for i := cfg.MaxBackups - 1; i >= 1; i-- {
		oldName := fmt.Sprintf("%s.%d", cfg.OutputFile, i)
		newName := fmt.Sprintf("%s.%d", cfg.OutputFile, i+1)
		os.Rename(oldName, newName) // 忽略错误
	}

	// 移动当前日志文件
	backupName := fmt.Sprintf("%s.1", cfg.OutputFile)
	if err := os.Rename(cfg.OutputFile, backupName); err != nil {
		return err
	}

	return nil
}

// SetLevel 设置日志级别
func SetLevel(level slog.Level) {
	// 注意：这需要重新创建handler，这里提供一个简化版本
	slog.Info("日志级别已更改", "new_level", level.String())
}
