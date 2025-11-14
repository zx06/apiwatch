<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { MonitorRule } from '../types/models'
  
  export let rule: MonitorRule
  
  const dispatch = createEventDispatcher()
  
  function formatTime(time: string): string {
    if (!time || time === '0001-01-01T00:00:00Z') {
      return '从未检查'
    }
    return new Date(time).toLocaleString('zh-CN')
  }
  
  function getExtractorTypeName(type: string): string {
    const typeMap: Record<string, string> = {
      css: 'CSS选择器',
      regex: '正则表达式',
      json: 'JSON路径'
    }
    return typeMap[type] || type
  }
</script>

<div class="rule-detail">
  <div class="detail-header">
    <h2>{rule.name}</h2>
    <button class="btn-check-now" on:click={() => dispatch('check-now', rule.id)}>
      🔍 立即检查
    </button>
  </div>
  
  <div class="detail-section">
    <h3>基本信息</h3>
    <div class="info-grid">
      <div class="info-item">
        <label>描述</label>
        <span>{rule.description || '无'}</span>
      </div>
      <div class="info-item">
        <label>URL</label>
        <span class="url-text">{rule.url}</span>
      </div>
      <div class="info-item">
        <label>HTTP方法</label>
        <span>{rule.method}</span>
      </div>
      <div class="info-item">
        <label>检查间隔</label>
        <span>{rule.interval}</span>
      </div>
    </div>
  </div>
  
  {#if rule.headers && Object.keys(rule.headers).length > 0}
    <div class="detail-section">
      <h3>请求头</h3>
      <div class="headers-list">
        {#each Object.entries(rule.headers) as [key, value]}
          <div class="header-item">
            <span class="header-key">{key}:</span>
            <span class="header-value">{value}</span>
          </div>
        {/each}
      </div>
    </div>
  {/if}
  
  {#if rule.body}
    <div class="detail-section">
      <h3>请求体</h3>
      <pre class="code-block">{rule.body}</pre>
    </div>
  {/if}
  
  <div class="detail-section">
    <h3>内容提取</h3>
    <div class="info-grid">
      <div class="info-item">
        <label>提取器类型</label>
        <span>{getExtractorTypeName(rule.extractor_type)}</span>
      </div>
      <div class="info-item">
        <label>提取表达式</label>
        <span class="code-text">{rule.extractor_expr}</span>
      </div>
    </div>
  </div>
  
  <div class="detail-section">
    <h3>状态信息</h3>
    <div class="info-grid">
      <div class="info-item">
        <label>当前状态</label>
        <span class="status-badge {rule.status}">{rule.status}</span>
      </div>
      <div class="info-item">
        <label>最后检查时间</label>
        <span>{formatTime(rule.last_checked)}</span>
      </div>
      <div class="info-item">
        <label>通知</label>
        <span>{rule.notify_enabled ? '✅ 已启用' : '❌ 已禁用'}</span>
      </div>
      <div class="info-item">
        <label>规则状态</label>
        <span>{rule.enabled ? '✅ 已启用' : '❌ 已禁用'}</span>
      </div>
    </div>
  </div>
  
  {#if rule.error_message}
    <div class="detail-section error-section">
      <h3>错误信息</h3>
      <div class="error-message">{rule.error_message}</div>
    </div>
  {/if}
  
  <div class="detail-section">
    <h3>最新提取内容</h3>
    <div class="content-box">
      {#if rule.last_content}
        <pre>{rule.last_content}</pre>
      {:else}
        <p class="no-content">暂无内容</p>
      {/if}
    </div>
  </div>
</div>

<style>
  .rule-detail {
    flex: 1;
    overflow-y: auto;
    padding: 20px;
    background: var(--bg-primary);
    transition: background-color 0.3s ease;
  }
  
  .detail-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    padding-bottom: 15px;
    border-bottom: 2px solid var(--color-primary);
  }
  
  .detail-header h2 {
    margin: 0;
    color: var(--text-primary);
  }
  
  .btn-check-now {
    padding: 8px 16px;
    background: var(--color-success);
    color: var(--text-inverse);
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s;
  }
  
  .btn-check-now:hover {
    opacity: 0.9;
  }
  
  .detail-section {
    margin-bottom: 25px;
  }
  
  .detail-section h3 {
    margin: 0 0 12px 0;
    font-size: 16px;
    color: var(--text-secondary);
    border-bottom: 1px solid var(--border-secondary);
    padding-bottom: 8px;
  }
  
  .info-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
  }
  
  .info-item {
    display: flex;
    flex-direction: column;
    gap: 5px;
  }
  
  .info-item label {
    font-size: 12px;
    color: var(--text-secondary);
    font-weight: 500;
  }
  
  .info-item span {
    font-size: 14px;
    color: var(--text-primary);
  }
  
  .url-text {
    word-break: break-all;
    color: var(--color-primary);
  }
  
  .code-text {
    font-family: 'Courier New', monospace;
    background: var(--bg-tertiary);
    color: var(--text-primary);
    padding: 4px 8px;
    border-radius: 3px;
    font-size: 13px;
  }
  
  .status-badge {
    display: inline-block;
    padding: 4px 12px;
    border-radius: 12px;
    font-size: 12px;
    font-weight: 500;
  }
  
  .status-badge.running {
    background: var(--color-success-bg);
    color: var(--color-success-text);
  }
  
  .status-badge.paused {
    background: var(--color-warning-bg);
    color: var(--color-warning-text);
  }
  
  .status-badge.error {
    background: var(--color-error-bg);
    color: var(--color-error-text);
  }
  
  .status-badge.idle {
    background: var(--bg-tertiary);
    color: var(--text-secondary);
  }
  
  .headers-list {
    background: var(--bg-tertiary);
    padding: 12px;
    border-radius: 4px;
  }
  
  .header-item {
    margin-bottom: 8px;
    font-size: 13px;
  }
  
  .header-key {
    font-weight: 600;
    color: var(--text-secondary);
  }
  
  .header-value {
    color: var(--text-primary);
    margin-left: 8px;
  }
  
  .code-block {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    padding: 12px;
    border-radius: 4px;
    overflow-x: auto;
    font-size: 13px;
    font-family: 'Courier New', monospace;
    margin: 0;
  }
  
  .error-section {
    background: var(--color-warning-bg);
    padding: 15px;
    border-radius: 4px;
    border-left: 4px solid var(--color-warning);
  }
  
  .error-message {
    color: var(--color-warning-text);
    font-size: 14px;
  }
  
  .content-box {
    background: var(--bg-tertiary);
    padding: 15px;
    border-radius: 4px;
    max-height: 300px;
    overflow-y: auto;
  }
  
  .content-box pre {
    margin: 0;
    white-space: pre-wrap;
    word-wrap: break-word;
    font-size: 13px;
    font-family: 'Courier New', monospace;
    color: var(--text-primary);
  }
  
  .no-content {
    color: var(--text-tertiary);
    text-align: center;
    margin: 20px 0;
  }
</style>
