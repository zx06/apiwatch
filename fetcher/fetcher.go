package fetcher

import (
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Request HTTP请求参数
type Request struct {
	URL     string
	Method  string
	Headers map[string]string
	Body    string
}

// Response HTTP响应
type Response struct {
	Body        []byte
	ContentType string
	StatusCode  int
}

// Fetcher HTTP客户端接口
type Fetcher interface {
	// Fetch 发送HTTP请求并获取响应
	Fetch(req *Request) (*Response, error)
}

// HTTPFetcher HTTP客户端实现
type HTTPFetcher struct {
	client *http.Client
}

// NewHTTPFetcher 创建HTTP客户端
func NewHTTPFetcher() *HTTPFetcher {
	return &HTTPFetcher{
		client: &http.Client{
			Timeout: 30 * time.Second,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				// 最多允许10次重定向
				if len(via) >= 10 {
					return fmt.Errorf("重定向次数过多")
				}
				return nil
			},
		},
	}
}

// Fetch 发送HTTP请求并获取响应
func (f *HTTPFetcher) Fetch(req *Request) (*Response, error) {
	var lastErr error

	// 重试机制：最多3次，指数退避
	for attempt := 0; attempt < 3; attempt++ {
		if attempt > 0 {
			// 指数退避：1s, 2s, 4s
			backoff := time.Duration(1<<uint(attempt-1)) * time.Second
			time.Sleep(backoff)
		}

		resp, err := f.doRequest(req)
		if err == nil {
			return resp, nil
		}

		lastErr = err
	}

	return nil, fmt.Errorf("请求失败（已重试3次）: %w", lastErr)
}

// doRequest 执行单次HTTP请求
func (f *HTTPFetcher) doRequest(req *Request) (*Response, error) {
	// 创建HTTP请求
	var bodyReader io.Reader
	if req.Body != "" {
		bodyReader = strings.NewReader(req.Body)
	}

	httpReq, err := http.NewRequest(req.Method, req.URL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认User-Agent
	httpReq.Header.Set("User-Agent", "URL-Monitor/1.0")

	// 设置自定义请求头
	for key, value := range req.Headers {
		httpReq.Header.Set(key, value)
	}

	// 发送请求
	httpResp, err := f.client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("发送请求失败: %w", err)
	}
	defer httpResp.Body.Close()

	// 检查状态码
	if httpResp.StatusCode < 200 || httpResp.StatusCode >= 300 {
		return nil, fmt.Errorf("HTTP错误: %d %s", httpResp.StatusCode, httpResp.Status)
	}

	// 限制响应大小（最大10MB）
	const maxBodySize = 10 * 1024 * 1024
	limitedReader := io.LimitReader(httpResp.Body, maxBodySize)

	// 读取响应体
	body, err := io.ReadAll(limitedReader)
	if err != nil {
		return nil, fmt.Errorf("读取响应失败: %w", err)
	}

	// 检查是否超过大小限制
	if len(body) == maxBodySize {
		// 尝试再读一个字节，看是否还有数据
		var buf [1]byte
		if n, _ := httpResp.Body.Read(buf[:]); n > 0 {
			return nil, fmt.Errorf("响应体过大（超过10MB）")
		}
	}

	return &Response{
		Body:        body,
		ContentType: httpResp.Header.Get("Content-Type"),
		StatusCode:  httpResp.StatusCode,
	}, nil
}
