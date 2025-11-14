<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import { theme, toggleTheme } from '../stores/theme'
  
  const dispatch = createEventDispatcher()
</script>

<div class="toolbar">
  <button class="btn btn-primary" on:click={() => dispatch('add-rule')}>
    ➕ 添加规则
  </button>
  <button class="btn btn-secondary" on:click={() => dispatch('refresh')}>
    🔄 刷新
  </button>
  <div class="spacer"></div>
  <button class="btn btn-theme" on:click={toggleTheme} title="切换主题">
    {$theme === 'light' ? '🌙 暗色' : '☀️ 亮色'}
  </button>
</div>

<style>
  .toolbar {
    display: flex;
    gap: 10px;
    padding: 15px;
    background: var(--bg-secondary);
    border-bottom: 1px solid var(--border-primary);
    transition: background-color 0.3s ease, border-color 0.3s ease;
  }
  
  .spacer {
    flex: 1;
  }
  
  .btn {
    padding: 8px 16px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s;
  }
  
  .btn:hover {
    opacity: 0.8;
  }
  
  .btn-primary {
    background: var(--color-primary);
    color: var(--text-inverse);
  }
  
  .btn-primary:hover {
    background: var(--color-primary-hover);
  }
  
  .btn-secondary {
    background: var(--color-secondary);
    color: var(--text-inverse);
  }
  
  .btn-secondary:hover {
    background: var(--color-secondary-hover);
  }
  
  .btn-theme {
    background: var(--bg-tertiary);
    color: var(--text-primary);
    border: 1px solid var(--border-primary);
  }
  
  .btn-theme:hover {
    background: var(--bg-hover);
  }
</style>
