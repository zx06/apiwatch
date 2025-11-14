package extractor

import (
	"fmt"

	"github.com/tidwall/gjson"
)

// JSONExtractor JSON路径提取器
type JSONExtractor struct {
	path string
}

// NewJSONExtractor 创建JSON路径提取器
func NewJSONExtractor(path string) *JSONExtractor {
	return &JSONExtractor{
		path: path,
	}
}

// Extract 使用JSON路径提取内容
func (e *JSONExtractor) Extract(body []byte, contentType string) (string, error) {
	// 验证JSON格式
	if !gjson.ValidBytes(body) {
		return "", fmt.Errorf("无效的JSON格式")
	}

	// 使用gjson查询
	result := gjson.GetBytes(body, e.path)

	// 检查是否存在
	if !result.Exists() {
		return "", fmt.Errorf("JSON路径未找到: %s", e.path)
	}

	// 返回结果的字符串表示
	// 如果是数组或对象，返回JSON字符串
	// 如果是基本类型，返回其值
	return result.String(), nil
}
