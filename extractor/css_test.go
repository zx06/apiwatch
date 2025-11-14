package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCSSExtractor_Extract(t *testing.T) {
	tests := []struct {
		name        string
		html        string
		selector    string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:     "提取单个元素",
			html:     `<html><body><h1 class="title">Hello World</h1></body></html>`,
			selector: ".title",
			want:     "Hello World",
			wantErr:  false,
		},
		{
			name: "提取多个元素",
			html: `<html><body>
				<p class="item">Item 1</p>
				<p class="item">Item 2</p>
				<p class="item">Item 3</p>
			</body></html>`,
			selector: ".item",
			want:     "Item 1\nItem 2\nItem 3",
			wantErr:  false,
		},
		{
			name: "提取嵌套元素",
			html: `<html><body>
				<div id="content">
					<span>Nested Content</span>
				</div>
			</body></html>`,
			selector: "#content span",
			want:     "Nested Content",
			wantErr:  false,
		},
		{
			name:        "选择器未匹配",
			html:        `<html><body><p>Test</p></body></html>`,
			selector:    ".nonexistent",
			wantErr:     true,
			errContains: "未匹配到任何元素",
		},
		{
			name:        "匹配元素无文本",
			html:        `<html><body><div class="empty"></div></body></html>`,
			selector:    ".empty",
			wantErr:     true,
			errContains: "没有文本内容",
		},
		{
			name:        "无效的HTML",
			html:        `not valid html`,
			selector:    ".title",
			wantErr:     true, // 无法匹配到.title
			errContains: "未匹配到任何元素",
		},
		{
			name:     "提取属性值",
			html:     `<html><body><a href="https://example.com" class="link">Link</a></body></html>`,
			selector: ".link",
			want:     "Link",
			wantErr:  false,
		},
		{
			name: "忽略空白文本",
			html: `<html><body>
				<p class="text">   Content with spaces   </p>
			</body></html>`,
			selector: ".text",
			want:     "Content with spaces",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewCSSExtractor(tt.selector)
			result, err := extractor.Extract([]byte(tt.html), "text/html")

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
