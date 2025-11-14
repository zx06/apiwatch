package extractor

import (
	"fmt"

	"github.com/zx06/apiwatch/models"
)

// Extractor 内容提取器接口
type Extractor interface {
	// Extract 从响应中提取内容
	Extract(body []byte, contentType string) (string, error)
}

// Factory 提取器工厂
type Factory struct{}

// NewFactory 创建提取器工厂
func NewFactory() *Factory {
	return &Factory{}
}

// Create 根据类型创建提取器
func (f *Factory) Create(extractorType models.ExtractorType, expr string) (Extractor, error) {
	switch extractorType {
	case models.ExtractorCSS:
		return NewCSSExtractor(expr), nil
	case models.ExtractorRegex:
		return NewRegexExtractor(expr)
	case models.ExtractorJSON:
		return NewJSONExtractor(expr), nil
	default:
		return nil, fmt.Errorf("不支持的提取器类型: %s", extractorType)
	}
}
