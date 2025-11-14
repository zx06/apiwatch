package extractor

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"
)

// RegexExtractor 正则表达式提取器
type RegexExtractor struct {
	pattern *regexp.Regexp
	timeout time.Duration
}

// NewRegexExtractor 创建正则表达式提取器
func NewRegexExtractor(pattern string) (*RegexExtractor, error) {
	// 预编译正则表达式
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("无效的正则表达式: %w", err)
	}

	return &RegexExtractor{
		pattern: re,
		timeout: 5 * time.Second, // 设置5秒超时防止ReDoS攻击
	}, nil
}

// Extract 使用正则表达式提取内容
func (e *RegexExtractor) Extract(body []byte, contentType string) (string, error) {
	// 使用context实现超时控制
	ctx, cancel := context.WithTimeout(context.Background(), e.timeout)
	defer cancel()

	// 在goroutine中执行正则匹配
	resultCh := make(chan []string, 1)
	errCh := make(chan error, 1)

	go func() {
		// 查找所有匹配
		matches := e.pattern.FindAllStringSubmatch(string(body), -1)
		if len(matches) == 0 {
			errCh <- fmt.Errorf("正则表达式未匹配到任何内容: %s", e.pattern.String())
			return
		}

		// 提取匹配结果
		var results []string
		for _, match := range matches {
			if len(match) > 1 {
				// 如果有捕获组，使用第一个捕获组
				results = append(results, match[1])
			} else if len(match) > 0 {
				// 否则使用整个匹配
				results = append(results, match[0])
			}
		}

		resultCh <- results
	}()

	// 等待结果或超时
	select {
	case <-ctx.Done():
		return "", fmt.Errorf("正则表达式匹配超时（可能存在ReDoS风险）")
	case err := <-errCh:
		return "", err
	case results := <-resultCh:
		// 用换行符连接所有匹配结果
		return strings.Join(results, "\n"), nil
	}
}
