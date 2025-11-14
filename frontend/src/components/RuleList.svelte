<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { MonitorRule } from '../types/models'
  
  export let rules: MonitorRule[]
  export let selectedId: string | null
  
  const dispatch = createEventDispatcher()
  
  function formatTime(time: string): string {
    if (!time || time === '0001-01-01T00:00:00Z') {
      return '从未检查'
    }
    return new Date(time).toLocaleString('zh-CN')
  }
  
  function getStatusText(status: string): string {
    const statusMap: Record<string, string> = {
      running: '运行中',
      paused: '已暂停',
      error: '错误',
      idle: '空闲'
    }
    return statusMap[status] || status
  }
</script>

<div class="rule-list">
  {#if rules.length === 0}
    <div class="empty-state">
      <p>暂无监控规则</p>
      <p class="hint">点击"添加规则"开始创建</p>
    </div>
  {:else}
    {#each rules as rule (rule.id)}
      <div 
        class="rule-item"
        class:selected={rule.id === selectedId}
        on:click={() => dispatch('select', rule.id)}
        role="button"
        tabindex="0"
        on:keypress={(e) => e.key === 'Enter' && dispatch('select', rule.id)}
      >
        <div class="rule-header">
          <span class="rule-name">{rule.name}</span>
          <span class="status-badge {rule.status}">{getStatusText(rule.status)}</span>
        </div>
        <div class="rule-info">
          <span class="url" title={rule.url}>{rule.url}</span>
          <span class="last-checked">{formatTime(rule.last_checked)}</span>
        </div>
        <div class="rule-actions">
          <button 
            class="btn-action" 
            on:click|stopPropagation={() => dispatch('toggle', rule.id)}
            title={rule.enabled ? '暂停监控' : '启动监控'}
          >
            {rule.enabled ? '⏸️ 暂停' : '▶️ 启动'}
          </button>
          <button 
            class="btn-action" 
            on:click|stopPropagation={() => dispatch('edit', rule)}
            title="编辑规则"
          >
            ✏️ 编辑
          </button>
          <button 
            class="btn-action btn-danger" 
            on:click|stopPropagation={() => dispatch('delete', rule.id)}
            title="删除规则"
          >
            🗑️ 删除
          </button>
        </div>
      </div>
    {/each}
  {/if}
</div>

<style>
  .rule-list {
    flex: 1;
    overflow-y: auto;
    background: var(--bg-primary);
    transition: background-color 0.3s ease;
  }
  
  .empty-state {
    display: flex;
    flex-direction: column;
    align-items: center;
    justify-content: center;
    height: 100%;
    color: var(--text-tertiary);
  }
  
  .empty-state p {
    margin: 5px 0;
  }
  
  .hint {
    font-size: 12px;
  }
  
  .rule-item {
    padding: 15px;
    border-bottom: 1px solid var(--border-secondary);
    cursor: pointer;
    transition: background 0.2s;
  }
  
  .rule-item:hover {
    background: var(--bg-tertiary);
  }
  
  .rule-item.selected {
    background: var(--bg-selected);
    border-left: 3px solid var(--color-primary);
  }
  
  .rule-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }
  
  .rule-name {
    font-weight: 600;
    font-size: 14px;
    color: var(--text-primary);
  }
  
  .status-badge {
    padding: 2px 8px;
    border-radius: 12px;
    font-size: 11px;
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
  
  .rule-info {
    display: flex;
    flex-direction: column;
    gap: 4px;
    margin-bottom: 10px;
  }
  
  .url {
    font-size: 12px;
    color: var(--text-secondary);
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }
  
  .last-checked {
    font-size: 11px;
    color: var(--text-tertiary);
  }
  
  .rule-actions {
    display: flex;
    gap: 5px;
  }
  
  .btn-action {
    padding: 4px 8px;
    border: 1px solid var(--border-primary);
    background: var(--bg-primary);
    color: var(--text-primary);
    border-radius: 3px;
    cursor: pointer;
    font-size: 11px;
    transition: all 0.2s;
  }
  
  .btn-action:hover {
    background: var(--bg-secondary);
  }
  
  .btn-danger:hover {
    background: var(--color-error-bg);
    border-color: var(--color-error);
  }
</style>
