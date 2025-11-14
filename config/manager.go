package config

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/zx06/apiwatch/models"
	"gopkg.in/yaml.v3"
)

// Manager 配置管理器接口
type Manager interface {
	// Load 加载所有规则
	Load() ([]*models.MonitorRule, error)

	// Save 保存所有规则
	Save(rules []*models.MonitorRule) error

	// AddRule 添加新规则
	AddRule(rule *models.MonitorRule) error

	// UpdateRule 更新规则
	UpdateRule(rule *models.MonitorRule) error

	// DeleteRule 删除规则
	DeleteRule(id string) error

	// GetRule 获取单个规则
	GetRule(id string) (*models.MonitorRule, error)
}

// Config 配置文件结构
type Config struct {
	Version string                `yaml:"version"`
	Rules   []*models.MonitorRule `yaml:"rules"`
}

// YAMLManager YAML配置管理器
type YAMLManager struct {
	configPath string
	mu         sync.RWMutex
}

// NewYAMLManager 创建YAML配置管理器
func NewYAMLManager(configPath string) (*YAMLManager, error) {
	if configPath == "" {
		// 默认配置路径
		homeDir, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("获取用户目录失败: %w", err)
		}
		configPath = filepath.Join(homeDir, ".url-monitor", "config.yaml")
	}

	// 确保配置目录存在
	configDir := filepath.Dir(configPath)
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return nil, fmt.Errorf("创建配置目录失败: %w", err)
	}

	return &YAMLManager{
		configPath: configPath,
	}, nil
}

// Load 加载所有规则
func (m *YAMLManager) Load() ([]*models.MonitorRule, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	// 如果配置文件不存在，返回空列表
	if _, err := os.Stat(m.configPath); os.IsNotExist(err) {
		return []*models.MonitorRule{}, nil
	}

	data, err := os.ReadFile(m.configPath)
	if err != nil {
		return nil, fmt.Errorf("读取配置文件失败: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("解析配置文件失败: %w", err)
	}

	return config.Rules, nil
}

// Save 保存所有规则
func (m *YAMLManager) Save(rules []*models.MonitorRule) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	config := Config{
		Version: "1.0",
		Rules:   rules,
	}

	data, err := yaml.Marshal(&config)
	if err != nil {
		return fmt.Errorf("序列化配置失败: %w", err)
	}

	// 先写入临时文件，然后重命名（原子操作）
	tempPath := m.configPath + ".tmp"
	if err := os.WriteFile(tempPath, data, 0600); err != nil {
		return fmt.Errorf("写入临时配置文件失败: %w", err)
	}

	if err := os.Rename(tempPath, m.configPath); err != nil {
		os.Remove(tempPath) // 清理临时文件
		return fmt.Errorf("更新配置文件失败: %w", err)
	}

	return nil
}

// AddRule 添加新规则
func (m *YAMLManager) AddRule(rule *models.MonitorRule) error {
	rules, err := m.Load()
	if err != nil {
		return err
	}

	// 检查ID是否已存在
	for _, r := range rules {
		if r.ID == rule.ID {
			return fmt.Errorf("规则ID已存在: %s", rule.ID)
		}
	}

	rules = append(rules, rule)
	return m.Save(rules)
}

// UpdateRule 更新规则
func (m *YAMLManager) UpdateRule(rule *models.MonitorRule) error {
	rules, err := m.Load()
	if err != nil {
		return err
	}

	found := false
	for i, r := range rules {
		if r.ID == rule.ID {
			rules[i] = rule
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("规则不存在: %s", rule.ID)
	}

	return m.Save(rules)
}

// DeleteRule 删除规则
func (m *YAMLManager) DeleteRule(id string) error {
	rules, err := m.Load()
	if err != nil {
		return err
	}

	newRules := make([]*models.MonitorRule, 0, len(rules))
	found := false
	for _, r := range rules {
		if r.ID != id {
			newRules = append(newRules, r)
		} else {
			found = true
		}
	}

	if !found {
		return fmt.Errorf("规则不存在: %s", id)
	}

	return m.Save(newRules)
}

// GetRule 获取单个规则
func (m *YAMLManager) GetRule(id string) (*models.MonitorRule, error) {
	rules, err := m.Load()
	if err != nil {
		return nil, err
	}

	for _, r := range rules {
		if r.ID == id {
			return r, nil
		}
	}

	return nil, fmt.Errorf("规则不存在: %s", id)
}
