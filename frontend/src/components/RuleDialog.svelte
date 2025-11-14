<script lang="ts">
  import { createEventDispatcher } from 'svelte'
  import type { MonitorRule } from '../types/models'
  import { v4 as uuidv4 } from 'uuid'
  
  export let rule: MonitorRule | null
  
  const dispatch = createEventDispatcher()
  
  let formData = {
    id: rule?.id || uuidv4(),
    name: rule?.name || '',
    description: rule?.description || '',
    url: rule?.url || '',
    method: rule?.method || 'GET',
    headers: rule?.headers || {},
    body: rule?.body || '',
    interval: rule?.interval || '5m',
    extractor_type: rule?.extractor_type || 'css',
    extractor_expr: rule?.extractor_expr || '',
    notify_enabled: rule?.notify_enabled ?? true,
    enabled: rule?.enabled ?? true,
    last_content: rule?.last_content || '',
    last_checked: rule?.last_checked || '0001-01-01T00:00:00Z',
    status: rule?.status || 'idle',
    error_message: rule?.error_message || ''
  }
  
  let headerKey = ''
  let headerValue = ''
  let showBodyField = formData.method !== 'GET'
  
  $: showBodyField = formData.method !== 'GET'
  
  function addHeader() {
    if (headerKey && headerValue) {
      formData.headers = { ...formData.headers, [headerKey]: headerValue }
      headerKey = ''
      headerValue = ''
    }
  }
  
  function removeHeader(key: string) {
    const { [key]: _, ...rest } = formData.headers
    formData.headers = rest
  }
  
  function handleSubmit() {
    dispatch('save', formData)
  }
</script>

<div 
  class="dialog-overlay" 
  role="button" 
  tabindex="0"
  on:click={() => dispatch('cancel')}
  on:keydown={(e) => e.key === 'Escape' && dispatch('cancel')}
>
  <div 
    class="dialog" 
    role="dialog"
    aria-modal="true"
    tabindex="-1"
    on:click|stopPropagation
    on:keydown|stopPropagation
  >
    <div class="dialog-header">
      <h2>{rule?.id ? '编辑规则' : '添加规则'}</h2>
      <button class="btn-close" on:click={() => dispatch('cancel')}>✕</button>
    </div>
    
    <form on:submit|preventDefault={handleSubmit}>
      <div class="form-content">
        <div class="form-group">
          <label for="name">规则名称 *</label>
          <input 
            id="name"
            type="text"
            bind:value={formData.name} 
            required 
            placeholder="例如：监控API响应"
          />
        </div>
        
        <div class="form-group">
          <label for="description">描述</label>
          <textarea 
            id="description"
            bind:value={formData.description}
            rows="2"
            placeholder="规则的详细描述（可选）"
          ></textarea>
        </div>
        
        <div class="form-group">
          <label for="url">URL *</label>
          <input 
            id="url"
            type="url"
            bind:value={formData.url} 
            required 
            placeholder="https://example.com/api/data"
          />
        </div>
        
        <div class="form-row">
          <div class="form-group">
            <label for="method">HTTP方法</label>
            <select id="method" bind:value={formData.method}>
              <option value="GET">GET</option>
              <option value="POST">POST</option>
              <option value="PUT">PUT</option>
              <option value="DELETE">DELETE</option>
              <option value="PATCH">PATCH</option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="interval">检查间隔 *</label>
            <input 
              id="interval"
              type="text"
              bind:value={formData.interval} 
              required 
              placeholder="例如: 5m, 1h, 30s"
            />
          </div>
        </div>
        
        <div class="form-group">
          <div class="form-label">自定义请求头</div>
          <div class="headers-editor">
            {#if Object.keys(formData.headers).length > 0}
              <div class="headers-list">
                {#each Object.entries(formData.headers) as [key, value]}
                  <div class="header-item">
                    <span class="header-text">{key}: {value}</span>
                    <button 
                      type="button" 
                      class="btn-remove" 
                      on:click={() => removeHeader(key)}
                    >
                      ✕
                    </button>
                  </div>
                {/each}
              </div>
            {/if}
            <div class="header-input">
              <input 
                type="text" 
                bind:value={headerKey} 
                placeholder="键"
                class="header-key-input"
              />
              <input 
                type="text" 
                bind:value={headerValue} 
                placeholder="值"
                class="header-value-input"
              />
              <button type="button" class="btn-add" on:click={addHeader}>
                ➕
              </button>
            </div>
          </div>
        </div>
        
        {#if showBodyField}
          <div class="form-group">
            <label for="body">请求体</label>
            <textarea 
              id="body"
              bind:value={formData.body}
              rows="4"
              placeholder="{`{\"key\": \"value\"}`}"
            ></textarea>
          </div>
        {/if}
        
        <div class="form-row">
          <div class="form-group">
            <label for="extractor_type">提取器类型</label>
            <select id="extractor_type" bind:value={formData.extractor_type}>
              <option value="css">CSS选择器</option>
              <option value="regex">正则表达式</option>
              <option value="json">JSON路径</option>
            </select>
          </div>
          
          <div class="form-group">
            <label for="extractor_expr">提取表达式 *</label>
            <input 
              id="extractor_expr"
              type="text"
              bind:value={formData.extractor_expr} 
              required 
              placeholder={
                formData.extractor_type === 'css' ? 'h1.title' :
                formData.extractor_type === 'regex' ? '\\d+' :
                'data.items[0].title'
              }
            />
          </div>
        </div>
        
        <div class="form-group checkbox-group">
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={formData.notify_enabled} />
            <span>启用通知</span>
          </label>
          
          <label class="checkbox-label">
            <input type="checkbox" bind:checked={formData.enabled} />
            <span>启用规则</span>
          </label>
        </div>
      </div>
      
      <div class="dialog-actions">
        <button type="submit" class="btn btn-primary">保存</button>
        <button type="button" class="btn btn-secondary" on:click={() => dispatch('cancel')}>
          取消
        </button>
      </div>
    </form>
  </div>
</div>

<style>
  .dialog-overlay {
    position: fixed;
    top: 0;
    left: 0;
    right: 0;
    bottom: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    z-index: 1000;
  }
  
  .dialog {
    background: var(--bg-primary);
    border-radius: 8px;
    width: 90%;
    max-width: 600px;
    max-height: 90vh;
    display: flex;
    flex-direction: column;
    box-shadow: var(--shadow-lg);
    transition: background-color 0.3s ease;
  }
  
  .dialog-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 20px;
    border-bottom: 1px solid var(--border-secondary);
  }
  
  .dialog-header h2 {
    margin: 0;
    font-size: 20px;
    color: var(--text-primary);
  }
  
  .btn-close {
    background: none;
    border: none;
    font-size: 24px;
    cursor: pointer;
    color: var(--text-tertiary);
    padding: 0;
    width: 30px;
    height: 30px;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  
  .btn-close:hover {
    color: var(--text-primary);
  }
  
  .form-content {
    padding: 20px;
    overflow-y: auto;
    flex: 1;
  }
  
  .form-group {
    margin-bottom: 20px;
  }
  
  .form-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 15px;
  }
  
  label,
  .form-label {
    display: block;
    margin-bottom: 6px;
    font-size: 14px;
    font-weight: 500;
    color: var(--text-secondary);
  }
  
  input[type="text"],
  input[type="url"],
  select,
  textarea {
    width: 100%;
    padding: 8px 12px;
    border: 1px solid var(--border-primary);
    background: var(--bg-primary);
    color: var(--text-primary);
    border-radius: 4px;
    font-size: 14px;
    font-family: inherit;
    transition: border-color 0.2s, background-color 0.3s ease;
  }
  
  input:focus,
  select:focus,
  textarea:focus {
    outline: none;
    border-color: var(--color-primary);
  }
  
  textarea {
    resize: vertical;
    font-family: 'Courier New', monospace;
  }
  
  .headers-editor {
    border: 1px solid var(--border-primary);
    border-radius: 4px;
    padding: 10px;
    background: var(--bg-tertiary);
  }
  
  .headers-list {
    margin-bottom: 10px;
  }
  
  .header-item {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding: 6px 10px;
    background: var(--bg-primary);
    border-radius: 3px;
    margin-bottom: 5px;
  }
  
  .header-text {
    font-size: 13px;
    color: var(--text-primary);
  }
  
  .btn-remove {
    background: none;
    border: none;
    color: var(--color-error);
    cursor: pointer;
    font-size: 16px;
    padding: 0 5px;
  }
  
  .header-input {
    display: flex;
    gap: 5px;
  }
  
  .header-key-input,
  .header-value-input {
    flex: 1;
    padding: 6px 10px;
    font-size: 13px;
  }
  
  .btn-add {
    padding: 6px 12px;
    background: var(--color-success);
    color: var(--text-inverse);
    border: none;
    border-radius: 3px;
    cursor: pointer;
    font-size: 14px;
  }
  
  .btn-add:hover {
    opacity: 0.9;
  }
  
  .checkbox-group {
    display: flex;
    gap: 20px;
  }
  
  .checkbox-label {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    margin-bottom: 0;
  }
  
  .checkbox-label input[type="checkbox"] {
    width: auto;
    cursor: pointer;
  }
  
  .checkbox-label span {
    font-size: 14px;
    color: var(--text-secondary);
  }
  
  .dialog-actions {
    display: flex;
    justify-content: flex-end;
    gap: 10px;
    padding: 15px 20px;
    border-top: 1px solid var(--border-secondary);
  }
  
  .btn {
    padding: 8px 20px;
    border: none;
    border-radius: 4px;
    cursor: pointer;
    font-size: 14px;
    transition: all 0.2s;
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
</style>
