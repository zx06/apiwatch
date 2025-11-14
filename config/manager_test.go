package config

import (
	"net/http"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zx06/apiwatch/models"
)

func TestNewYAMLManager(t *testing.T) {
	t.Run("使用默认路径", func(t *testing.T) {
		manager, err := NewYAMLManager("")
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.Contains(t, manager.configPath, ".url-monitor")
	})

	t.Run("使用自定义路径", func(t *testing.T) {
		tempDir := t.TempDir()
		configPath := filepath.Join(tempDir, "test-config.yaml")

		manager, err := NewYAMLManager(configPath)
		require.NoError(t, err)
		assert.Equal(t, configPath, manager.configPath)
	})
}

func TestYAMLManager_LoadAndSave(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	manager, err := NewYAMLManager(configPath)
	require.NoError(t, err)

	t.Run("加载不存在的配置文件返回空列表", func(t *testing.T) {
		rules, err := manager.Load()
		require.NoError(t, err)
		assert.Empty(t, rules)
	})

	t.Run("保存和加载规则", func(t *testing.T) {
		testRules := []*models.MonitorRule{
			{
				ID:            uuid.New().String(),
				Name:          "测试规则1",
				Description:   "这是一个测试规则",
				URL:           "https://example.com",
				Method:        http.MethodGet,
				Interval:      models.Duration(5 * time.Minute),
				ExtractorType: models.ExtractorCSS,
				ExtractorExpr: ".content",
				NotifyEnabled: true,
				Enabled:       true,
			},
			{
				ID:     uuid.New().String(),
				Name:   "测试规则2",
				URL:    "https://api.example.com",
				Method: http.MethodPost,
				Headers: map[string]string{
					"Authorization": "Bearer token",
				},
				Body:          `{"query": "test"}`,
				Interval:      models.Duration(10 * time.Minute),
				ExtractorType: models.ExtractorJSON,
				ExtractorExpr: "data.result",
				NotifyEnabled: false,
				Enabled:       false,
			},
		}

		// 保存规则
		err := manager.Save(testRules)
		require.NoError(t, err)

		// 验证文件存在
		_, err = os.Stat(configPath)
		require.NoError(t, err)

		// 加载规则
		loadedRules, err := manager.Load()
		require.NoError(t, err)
		require.Len(t, loadedRules, 2)

		// 验证规则内容
		assert.Equal(t, testRules[0].ID, loadedRules[0].ID)
		assert.Equal(t, testRules[0].Name, loadedRules[0].Name)
		assert.Equal(t, testRules[0].Description, loadedRules[0].Description)
		assert.Equal(t, testRules[0].URL, loadedRules[0].URL)
		assert.Equal(t, testRules[0].Method, loadedRules[0].Method)
		assert.Equal(t, testRules[0].Interval, loadedRules[0].Interval)

		assert.Equal(t, testRules[1].Headers, loadedRules[1].Headers)
		assert.Equal(t, testRules[1].Body, loadedRules[1].Body)
	})
}

func TestYAMLManager_AddRule(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	manager, err := NewYAMLManager(configPath)
	require.NoError(t, err)

	rule := &models.MonitorRule{
		ID:            uuid.New().String(),
		Name:          "新规则",
		URL:           "https://example.com",
		Method:        http.MethodGet,
		Interval:      models.Duration(5 * time.Minute),
		ExtractorType: models.ExtractorCSS,
		ExtractorExpr: ".content",
	}

	t.Run("添加新规则", func(t *testing.T) {
		err := manager.AddRule(rule)
		require.NoError(t, err)

		rules, err := manager.Load()
		require.NoError(t, err)
		require.Len(t, rules, 1)
		assert.Equal(t, rule.ID, rules[0].ID)
	})

	t.Run("添加重复ID的规则应失败", func(t *testing.T) {
		err := manager.AddRule(rule)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "规则ID已存在")
	})
}

func TestYAMLManager_UpdateRule(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	manager, err := NewYAMLManager(configPath)
	require.NoError(t, err)

	ruleID := uuid.New().String()
	rule := &models.MonitorRule{
		ID:            ruleID,
		Name:          "原始规则",
		URL:           "https://example.com",
		Method:        http.MethodGet,
		Interval:      models.Duration(5 * time.Minute),
		ExtractorType: models.ExtractorCSS,
		ExtractorExpr: ".content",
	}

	err = manager.AddRule(rule)
	require.NoError(t, err)

	t.Run("更新存在的规则", func(t *testing.T) {
		updatedRule := &models.MonitorRule{
			ID:            ruleID,
			Name:          "更新后的规则",
			URL:           "https://updated.example.com",
			Method:        http.MethodPost,
			Interval:      models.Duration(10 * time.Minute),
			ExtractorType: models.ExtractorJSON,
			ExtractorExpr: "data.value",
		}

		err := manager.UpdateRule(updatedRule)
		require.NoError(t, err)

		rules, err := manager.Load()
		require.NoError(t, err)
		require.Len(t, rules, 1)
		assert.Equal(t, "更新后的规则", rules[0].Name)
		assert.Equal(t, "https://updated.example.com", rules[0].URL)
		assert.Equal(t, http.MethodPost, rules[0].Method)
	})

	t.Run("更新不存在的规则应失败", func(t *testing.T) {
		nonExistentRule := &models.MonitorRule{
			ID:            uuid.New().String(),
			Name:          "不存在的规则",
			URL:           "https://example.com",
			Method:        http.MethodGet,
			Interval:      models.Duration(5 * time.Minute),
			ExtractorType: models.ExtractorCSS,
			ExtractorExpr: ".content",
		}

		err := manager.UpdateRule(nonExistentRule)
		require.Error(t, err)
		assert.Contains(t, err.Error(), "规则不存在")
	})
}

func TestYAMLManager_DeleteRule(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	manager, err := NewYAMLManager(configPath)
	require.NoError(t, err)

	rule1ID := uuid.New().String()
	rule2ID := uuid.New().String()

	rules := []*models.MonitorRule{
		{
			ID:            rule1ID,
			Name:          "规则1",
			URL:           "https://example1.com",
			Method:        http.MethodGet,
			Interval:      models.Duration(5 * time.Minute),
			ExtractorType: models.ExtractorCSS,
			ExtractorExpr: ".content",
		},
		{
			ID:            rule2ID,
			Name:          "规则2",
			URL:           "https://example2.com",
			Method:        http.MethodGet,
			Interval:      models.Duration(5 * time.Minute),
			ExtractorType: models.ExtractorCSS,
			ExtractorExpr: ".content",
		},
	}

	err = manager.Save(rules)
	require.NoError(t, err)

	t.Run("删除存在的规则", func(t *testing.T) {
		err := manager.DeleteRule(rule1ID)
		require.NoError(t, err)

		remainingRules, err := manager.Load()
		require.NoError(t, err)
		require.Len(t, remainingRules, 1)
		assert.Equal(t, rule2ID, remainingRules[0].ID)
	})

	t.Run("删除不存在的规则应失败", func(t *testing.T) {
		err := manager.DeleteRule(uuid.New().String())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "规则不存在")
	})
}

func TestYAMLManager_GetRule(t *testing.T) {
	tempDir := t.TempDir()
	configPath := filepath.Join(tempDir, "config.yaml")

	manager, err := NewYAMLManager(configPath)
	require.NoError(t, err)

	ruleID := uuid.New().String()
	rule := &models.MonitorRule{
		ID:            ruleID,
		Name:          "测试规则",
		URL:           "https://example.com",
		Method:        http.MethodGet,
		Interval:      models.Duration(5 * time.Minute),
		ExtractorType: models.ExtractorCSS,
		ExtractorExpr: ".content",
	}

	err = manager.AddRule(rule)
	require.NoError(t, err)

	t.Run("获取存在的规则", func(t *testing.T) {
		foundRule, err := manager.GetRule(ruleID)
		require.NoError(t, err)
		assert.Equal(t, rule.ID, foundRule.ID)
		assert.Equal(t, rule.Name, foundRule.Name)
	})

	t.Run("获取不存在的规则应失败", func(t *testing.T) {
		_, err := manager.GetRule(uuid.New().String())
		require.Error(t, err)
		assert.Contains(t, err.Error(), "规则不存在")
	})
}
