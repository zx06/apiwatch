package models

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMonitorRule_Validate(t *testing.T) {
	tests := []struct {
		name    string
		rule    *MonitorRule
		wantErr bool
		errMsg  string
	}{
		{
			name: "有效的规则",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com",
				Method:        http.MethodGet,
				Interval:      Duration(Duration(5 * time.Minute)),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: false,
		},
		{
			name: "空名称",
			rule: &MonitorRule{
				URL:           "https://example.com",
				Method:        http.MethodGet,
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: true,
			errMsg:  "规则名称不能为空",
		},
		{
			name: "空URL",
			rule: &MonitorRule{
				Name:          "测试规则",
				Method:        http.MethodGet,
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: true,
			errMsg:  "URL不能为空",
		},
		{
			name: "无效的URL",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "://invalid-url",
				Method:        http.MethodGet,
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: true,
			errMsg:  "无效的URL",
		},
		{
			name: "空方法应设置为默认GET",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com",
				Method:        "",
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: false,
		},
		{
			name: "无效的HTTP方法",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com",
				Method:        "INVALID",
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: true,
			errMsg:  "无效的HTTP方法",
		},
		{
			name: "间隔小于1秒",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com",
				Method:        http.MethodGet,
				Interval:      Duration(500 * time.Millisecond),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: true,
			errMsg:  "检查间隔不能小于1秒",
		},
		{
			name: "空提取表达式",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com",
				Method:        http.MethodGet,
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: "",
			},
			wantErr: true,
			errMsg:  "提取表达式不能为空",
		},
		{
			name: "无效的提取器类型",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com",
				Method:        http.MethodGet,
				Interval:      Duration(5 * time.Minute),
				ExtractorType: "invalid",
				ExtractorExpr: ".content",
			},
			wantErr: true,
			errMsg:  "无效的提取器类型",
		},
		{
			name: "POST方法带请求体",
			rule: &MonitorRule{
				Name:          "测试规则",
				URL:           "https://example.com/api",
				Method:        http.MethodPost,
				Body:          `{"key": "value"}`,
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorJSON,
				ExtractorExpr: "data.result",
			},
			wantErr: false,
		},
		{
			name: "带自定义请求头",
			rule: &MonitorRule{
				Name:   "测试规则",
				URL:    "https://example.com",
				Method: http.MethodGet,
				Headers: map[string]string{
					"Authorization": "Bearer token",
					"Content-Type":  "application/json",
				},
				Interval:      Duration(5 * time.Minute),
				ExtractorType: ExtractorCSS,
				ExtractorExpr: ".content",
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.rule.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
				// 验证空方法被设置为默认值
				if tt.rule.Method == "" {
					assert.Equal(t, DefaultMethod, tt.rule.Method)
				}
			}
		})
	}
}

func TestExtractorType_Constants(t *testing.T) {
	assert.Equal(t, ExtractorType("css"), ExtractorCSS)
	assert.Equal(t, ExtractorType("regex"), ExtractorRegex)
	assert.Equal(t, ExtractorType("json"), ExtractorJSON)
}

func TestRuleStatus_Constants(t *testing.T) {
	assert.Equal(t, RuleStatus("running"), StatusRunning)
	assert.Equal(t, RuleStatus("paused"), StatusPaused)
	assert.Equal(t, RuleStatus("error"), StatusError)
	assert.Equal(t, RuleStatus("idle"), StatusIdle)
}

func TestDefaultMethod(t *testing.T) {
	assert.Equal(t, http.MethodGet, DefaultMethod)
}
