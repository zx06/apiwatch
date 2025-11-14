package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJSONExtractor_Extract(t *testing.T) {
	tests := []struct {
		name        string
		json        string
		path        string
		want        string
		wantErr     bool
		errContains string
	}{
		{
			name:    "提取简单字段",
			json:    `{"name": "John", "age": 30}`,
			path:    "name",
			want:    "John",
			wantErr: false,
		},
		{
			name:    "提取嵌套字段",
			json:    `{"user": {"name": "John", "email": "john@example.com"}}`,
			path:    "user.name",
			want:    "John",
			wantErr: false,
		},
		{
			name:    "提取数组元素",
			json:    `{"items": ["apple", "banana", "orange"]}`,
			path:    "items.1",
			want:    "banana",
			wantErr: false,
		},
		{
			name:    "提取数组中的对象",
			json:    `{"users": [{"name": "John"}, {"name": "Jane"}]}`,
			path:    "users.0.name",
			want:    "John",
			wantErr: false,
		},
		{
			name:    "提取数字",
			json:    `{"count": 42, "price": 19.99}`,
			path:    "count",
			want:    "42",
			wantErr: false,
		},
		{
			name:    "提取布尔值",
			json:    `{"active": true, "verified": false}`,
			path:    "active",
			want:    "true",
			wantErr: false,
		},
		{
			name:    "提取整个数组",
			json:    `{"tags": ["go", "json", "api"]}`,
			path:    "tags",
			want:    `["go", "json", "api"]`,
			wantErr: false,
		},
		{
			name:    "提取整个对象",
			json:    `{"user": {"name": "John", "age": 30}}`,
			path:    "user",
			want:    `{"name": "John", "age": 30}`,
			wantErr: false,
		},
		{
			name:        "路径不存在",
			json:        `{"name": "John"}`,
			path:        "nonexistent",
			wantErr:     true,
			errContains: "未找到",
		},
		{
			name:        "无效的JSON",
			json:        `{invalid json}`,
			path:        "name",
			wantErr:     true,
			errContains: "无效的JSON格式",
		},
		{
			name:    "使用通配符",
			json:    `{"users": [{"name": "John"}, {"name": "Jane"}]}`,
			path:    "users.#.name",
			want:    `["John","Jane"]`,
			wantErr: false,
		},
		{
			name:    "使用条件查询",
			json:    `{"users": [{"name": "John", "age": 30}, {"name": "Jane", "age": 25}]}`,
			path:    "users.#(age>26).name",
			want:    `John`,
			wantErr: false,
		},
		{
			name:    "提取null值",
			json:    `{"value": null}`,
			path:    "value",
			want:    "",
			wantErr: false,
		},
		{
			name:    "深层嵌套",
			json:    `{"a": {"b": {"c": {"d": {"e": "deep value"}}}}}`,
			path:    "a.b.c.d.e",
			want:    "deep value",
			wantErr: false,
		},
		{
			name:        "数组索引越界",
			json:        `{"items": ["a", "b"]}`,
			path:        "items.10",
			wantErr:     true,
			errContains: "未找到",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewJSONExtractor(tt.path)
			result, err := extractor.Extract([]byte(tt.json), "application/json")

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

func TestJSONExtractor_ComplexScenarios(t *testing.T) {
	complexJSON := `{
		"status": "success",
		"data": {
			"posts": [
				{
					"id": 1,
					"title": "First Post",
					"author": {"name": "John", "verified": true},
					"tags": ["go", "programming"]
				},
				{
					"id": 2,
					"title": "Second Post",
					"author": {"name": "Jane", "verified": false},
					"tags": ["api", "rest"]
				}
			],
			"total": 2
		}
	}`

	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "提取状态",
			path: "status",
			want: "success",
		},
		{
			name: "提取总数",
			path: "data.total",
			want: "2",
		},
		{
			name: "提取第一篇文章标题",
			path: "data.posts.0.title",
			want: "First Post",
		},
		{
			name: "提取所有标题",
			path: "data.posts.#.title",
			want: `["First Post","Second Post"]`,
		},
		{
			name: "提取已验证作者的文章",
			path: "data.posts.#(author.verified==true).title",
			want: `First Post`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor := NewJSONExtractor(tt.path)
			result, err := extractor.Extract([]byte(complexJSON), "application/json")
			require.NoError(t, err)
			assert.Equal(t, tt.want, result)
		})
	}
}
