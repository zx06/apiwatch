package fetcher

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewHTTPFetcher(t *testing.T) {
	fetcher := NewHTTPFetcher()
	assert.NotNil(t, fetcher)
	assert.NotNil(t, fetcher.client)
	assert.Equal(t, 30*time.Second, fetcher.client.Timeout)
}

func TestHTTPFetcher_Fetch_Success(t *testing.T) {
	// 创建测试服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证请求
		assert.Equal(t, http.MethodGet, r.Method)
		assert.Equal(t, "URL-Monitor/1.0", r.Header.Get("User-Agent"))

		// 返回响应
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("<html><body>Test Content</body></html>"))
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
	}

	resp, err := fetcher.Fetch(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Equal(t, "text/html", resp.ContentType)
	assert.Contains(t, string(resp.Body), "Test Content")
}

func TestHTTPFetcher_Fetch_POST_WithBody(t *testing.T) {
	expectedBody := `{"key": "value"}`

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodPost, r.Method)

		// 读取请求体
		body := make([]byte, len(expectedBody))
		r.Body.Read(body)
		assert.Equal(t, expectedBody, string(body))

		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"result": "success"}`))
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodPost,
		Body:   expectedBody,
	}

	resp, err := fetcher.Fetch(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
	assert.Contains(t, string(resp.Body), "success")
}

func TestHTTPFetcher_Fetch_CustomHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 验证自定义请求头
		assert.Equal(t, "Bearer token123", r.Header.Get("Authorization"))
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "custom-value", r.Header.Get("X-Custom-Header"))

		w.WriteHeader(http.StatusOK)
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
		Headers: map[string]string{
			"Authorization":   "Bearer token123",
			"Content-Type":    "application/json",
			"X-Custom-Header": "custom-value",
		},
	}

	resp, err := fetcher.Fetch(req)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPFetcher_Fetch_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
	}

	_, err := fetcher.Fetch(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "HTTP错误: 404")
}

func TestHTTPFetcher_Fetch_Retry(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		if attemptCount < 3 {
			// 前两次请求返回错误
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			// 第三次请求成功
			w.WriteHeader(http.StatusOK)
			w.Write([]byte("Success"))
		}
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
	}

	resp, err := fetcher.Fetch(req)
	require.NoError(t, err)
	assert.Equal(t, 3, attemptCount)
	assert.Equal(t, http.StatusOK, resp.StatusCode)
}

func TestHTTPFetcher_Fetch_RetryExhausted(t *testing.T) {
	attemptCount := 0

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attemptCount++
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
	}

	_, err := fetcher.Fetch(req)
	require.Error(t, err)
	assert.Equal(t, 3, attemptCount)
	assert.Contains(t, err.Error(), "请求失败（已重试3次）")
}

func TestHTTPFetcher_Fetch_InvalidURL(t *testing.T) {
	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    "://invalid-url",
		Method: http.MethodGet,
	}

	_, err := fetcher.Fetch(req)
	require.Error(t, err)
}

func TestHTTPFetcher_Fetch_BodySizeLimit(t *testing.T) {
	// 创建一个返回大响应的服务器
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		// 写入超过10MB的数据
		largeData := strings.Repeat("x", 11*1024*1024)
		w.Write([]byte(largeData))
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
	}

	_, err := fetcher.Fetch(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "响应体过大")
}

func TestHTTPFetcher_Fetch_Redirect(t *testing.T) {
	// 创建目标服务器
	targetServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Final destination"))
	}))
	defer targetServer.Close()

	// 创建重定向服务器
	redirectServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, targetServer.URL, http.StatusFound)
	}))
	defer redirectServer.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    redirectServer.URL,
		Method: http.MethodGet,
	}

	resp, err := fetcher.Fetch(req)
	require.NoError(t, err)
	assert.Contains(t, string(resp.Body), "Final destination")
}

func TestHTTPFetcher_Fetch_TooManyRedirects(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 创建循环重定向到自己
		http.Redirect(w, r, r.URL.String(), http.StatusFound)
	}))
	defer server.Close()

	fetcher := NewHTTPFetcher()
	req := &Request{
		URL:    server.URL,
		Method: http.MethodGet,
	}

	_, err := fetcher.Fetch(req)
	require.Error(t, err)
	assert.Contains(t, err.Error(), "重定向")
}

func TestHTTPFetcher_Fetch_DifferentMethods(t *testing.T) {
	methods := []string{
		http.MethodGet,
		http.MethodPost,
		http.MethodPut,
		http.MethodDelete,
		http.MethodPatch,
	}

	for _, method := range methods {
		t.Run(method, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, method, r.Method)
				w.WriteHeader(http.StatusOK)
			}))
			defer server.Close()

			fetcher := NewHTTPFetcher()
			req := &Request{
				URL:    server.URL,
				Method: method,
			}

			resp, err := fetcher.Fetch(req)
			require.NoError(t, err)
			assert.Equal(t, http.StatusOK, resp.StatusCode)
		})
	}
}
