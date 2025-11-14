package extractor

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewRegexExtractor(t *testing.T) {
	tests := []struct {
		name    string
		pattern string
		wantErr bool
	}{
		{
			name:    "有效的正则表达式",
			pattern: `\d+`,
			wantErr: false,
		},
		{
			name:    "带捕获组的正则",
			pattern: `(\w+)@(\w+)\.com`,
			wantErr: false,
		},
		{
			name:    "无效的正则表达式",
			pattern: `[invalid`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, err := NewRegexExtractor(tt.pattern)
			if tt.wantErr {
				require.Error(t, err)
				assert.Nil(t, extractor)
			} else {
				require.NoError(t, err)
				assert.NotNil(t, extractor)
			}
		})
	}
}

func TestRegexExtractor_Extract(t *testing.T) {
	tests := []struct {
		name        string
		pattern     string
		text        string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "提取数字",
			pattern: `\d+`,
			text:    "The price is 100 dollars and 50 cents",
			want:    "100\n50",
			wantErr: false,
		},
		{
			name:    "提取邮箱",
			pattern: `[\w.]+@[\w.]+`,
			text:    "Contact us at info@example.com or support@test.org",
			want:    "info@example.com\nsupport@test.org",
			wantErr: false,
		},
		{
			name:    "使用捕获组",
			pattern: `name:\s*(\w+)`,
			text:    "name: John, age: 30, name: Jane",
			want:    "John\nJane",
			wantErr: false,
		},
		{
			name:    "提取URL",
			pattern: `https?://[^\s]+`,
			text:    "Visit https://example.com or http://test.org for more info",
			want:    "https://example.com\nhttp://test.org",
			wantErr: false,
		},
		{
			name:        "未匹配到内容",
			pattern:     `\d+`,
			text:        "No numbers here",
			wantErr:     true,
			errContains: "未匹配到任何内容",
		},
		{
			name:    "匹配单个结果",
			pattern: `error:\s*(.+)`,
			text:    "error: Something went wrong",
			want:    "Something went wrong",
			wantErr: false,
		},
		{
			name:    "匹配多行文本",
			pattern: `(?m)^Line \d+`,
			text:    "Line 1\nSome text\nLine 2\nMore text\nLine 3",
			want:    "Line 1\nLine 2\nLine 3",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, err := NewRegexExtractor(tt.pattern)
			require.NoError(t, err)

			result, err := extractor.Extract([]byte(tt.text), "text/plain")

			if tt.wantErr {
				require.Error(t, err)
				if tt.errContains != "" {
					assert.Contains(t, err.Error(), tt.errContains)
				}
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.want, result)
			}
		})
	}
}

func TestRegexExtractor_Timeout(t *testing.T) {
	// 创建一个可能导致ReDoS的正则表达式
	// 注意：这个测试可能需要较长时间，但应该在超时时间内完成
	pattern := `(a+)+b`
	extractor, err := NewRegexExtractor(pattern)
	require.NoError(t, err)

	// 创建一个会导致大量回溯的输入
	text := strings.Repeat("a", 25) + "c" // 故意不匹配

	result, err := extractor.Extract([]byte(text), "text/plain")

	// 应该返回未匹配错误，而不是超时
	// 因为Go的regexp引擎使用了优化算法，不会出现灾难性回溯
	if err != nil {
		assert.Contains(t, err.Error(), "未匹配到任何内容")
	} else {
		// 如果匹配成功，验证结果
		assert.NotEmpty(t, result)
	}
}
