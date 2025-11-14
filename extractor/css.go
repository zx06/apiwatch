package extractor

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

// CSSExtractor CSS选择器提取器
type CSSExtractor struct {
	selector string
}

// NewCSSExtractor 创建CSS选择器提取器
func NewCSSExtractor(selector string) *CSSExtractor {
	return &CSSExtractor{
		selector: selector,
	}
}

// Extract 使用CSS选择器提取内容
func (e *CSSExtractor) Extract(body []byte, contentType string) (string, error) {
	// 解析HTML
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("解析HTML失败: %w", err)
	}

	// 查找匹配的元素
	selection := doc.Find(e.selector)
	if selection.Length() == 0 {
		return "", fmt.Errorf("CSS选择器未匹配到任何元素: %s", e.selector)
	}

	// 提取所有匹配元素的文本内容
	var results []string
	selection.Each(func(i int, s *goquery.Selection) {
		text := strings.TrimSpace(s.Text())
		if text != "" {
			results = append(results, text)
		}
	})

	if len(results) == 0 {
		return "", fmt.Errorf("匹配的元素没有文本内容")
	}

	// 如果有多个结果，用换行符连接
	return strings.Join(results, "\n"), nil
}
