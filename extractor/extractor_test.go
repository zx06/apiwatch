package extractor

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zx06/apiwatch/models"
)

func TestFactory_Create(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name          string
		extractorType models.ExtractorType
		expr          string
		wantErr       bool
	}{
		{
			name:          "创建CSS提取器",
			extractorType: models.ExtractorCSS,
			expr:          ".content",
			wantErr:       false,
		},
		{
			name:          "创建正则提取器",
			extractorType: models.ExtractorRegex,
			expr:          `\d+`,
			wantErr:       false,
		},
		{
			name:          "创建JSON提取器",
			extractorType: models.ExtractorJSON,
			expr:          "data.value",
			wantErr:       false,
		},
		{
			name:          "无效的正则表达式",
			extractorType: models.ExtractorRegex,
			expr:          `[invalid`,
			wantErr:       true,
		},
		{
			name:          "不支持的提取器类型",
			extractorType: "invalid",
			expr:          "test",
			wantErr:       true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			extractor, err := factory.Create(tt.extractorType, tt.expr)
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
