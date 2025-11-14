package models

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"time"
)

// Duration 自定义Duration类型，支持JSON序列化/反序列化
type Duration time.Duration

// MarshalJSON 实现JSON序列化
func (d Duration) MarshalJSON() ([]byte, error) {
	return json.Marshal(time.Duration(d).String())
}

// UnmarshalJSON 实现JSON反序列化
func (d *Duration) UnmarshalJSON(b []byte) error {
	var v interface{}
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	switch value := v.(type) {
	case float64:
		*d = Duration(time.Duration(value))
		return nil
	case string:
		tmp, err := time.ParseDuration(value)
		if err != nil {
			return err
		}
		*d = Duration(tmp)
		return nil
	default:
		return errors.New("invalid duration")
	}
}

// MarshalYAML 实现YAML序列化
func (d Duration) MarshalYAML() (interface{}, error) {
	return time.Duration(d).String(), nil
}

// UnmarshalYAML 实现YAML反序列化
func (d *Duration) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var s string
	if err := unmarshal(&s); err != nil {
		return err
	}
	tmp, err := time.ParseDuration(s)
	if err != nil {
		return err
	}
	*d = Duration(tmp)
	return nil
}

// ExtractorType 提取器类型
type ExtractorType string

const (
	ExtractorCSS   ExtractorType = "css"
	ExtractorRegex ExtractorType = "regex"
	ExtractorJSON  ExtractorType = "json"
)

// RuleStatus 规则状态
type RuleStatus string

const (
	StatusRunning RuleStatus = "running"
	StatusPaused  RuleStatus = "paused"
	StatusError   RuleStatus = "error"
	StatusIdle    RuleStatus = "idle"
)

// DefaultMethod 默认HTTP方法
const DefaultMethod = http.MethodGet

// MonitorRule 监控规则
type MonitorRule struct {
	ID            string            `json:"id" yaml:"id"`
	Name          string            `json:"name" yaml:"name"`
	Description   string            `json:"description" yaml:"description"`
	URL           string            `json:"url" yaml:"url"`
	Method        string            `json:"method" yaml:"method"`
	Headers       map[string]string `json:"headers,omitempty" yaml:"headers,omitempty"`
	Body          string            `json:"body,omitempty" yaml:"body,omitempty"`
	Interval      Duration          `json:"interval" yaml:"interval"`
	ExtractorType ExtractorType     `json:"extractor_type" yaml:"extractor_type"`
	ExtractorExpr string            `json:"extractor_expr" yaml:"extractor_expr"`
	NotifyEnabled bool              `json:"notify_enabled" yaml:"notify_enabled"`
	Enabled       bool              `json:"enabled" yaml:"enabled"`
	LastContent   string            `json:"last_content" yaml:"last_content"`
	LastChecked   string            `json:"last_checked" yaml:"last_checked"` // RFC3339 格式的时间字符串
	Status        RuleStatus        `json:"status" yaml:"status"`
	ErrorMessage  string            `json:"error_message,omitempty" yaml:"error_message,omitempty"`
}

// Validate 验证规则的有效性
func (r *MonitorRule) Validate() error {
	if r.Name == "" {
		return errors.New("规则名称不能为空")
	}

	if r.URL == "" {
		return errors.New("URL不能为空")
	}

	if _, err := url.Parse(r.URL); err != nil {
		return fmt.Errorf("无效的URL: %w", err)
	}

	if r.Method == "" {
		r.Method = DefaultMethod
	}

	// 验证HTTP方法
	validMethods := map[string]bool{
		http.MethodGet:     true,
		http.MethodPost:    true,
		http.MethodPut:     true,
		http.MethodDelete:  true,
		http.MethodPatch:   true,
		http.MethodHead:    true,
		http.MethodOptions: true,
	}
	if !validMethods[r.Method] {
		return fmt.Errorf("无效的HTTP方法: %s", r.Method)
	}

	if time.Duration(r.Interval) < time.Second {
		return errors.New("检查间隔不能小于1秒")
	}

	if r.ExtractorExpr == "" {
		return errors.New("提取表达式不能为空")
	}

	// 验证提取器类型
	validExtractors := map[ExtractorType]bool{
		ExtractorCSS: true, ExtractorRegex: true, ExtractorJSON: true,
	}
	if !validExtractors[r.ExtractorType] {
		return fmt.Errorf("无效的提取器类型: %s", r.ExtractorType)
	}

	return nil
}
