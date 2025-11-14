<script lang="ts">
  import { onMount } from 'svelte'
  import Toolbar from './components/Toolbar.svelte'
  import RuleList from './components/RuleList.svelte'
  import RuleDetail from './components/RuleDetail.svelte'
  import RuleDialog from './components/RuleDialog.svelte'
  import { rulesStore, selectedRule } from './stores/rules'
  import type { MonitorRule } from './types/models'
  
  let showDialog = $state(false)
  let editingRule = $state<MonitorRule | null>(null)
  
  function showAddDialog() {
    editingRule = null
    showDialog = true
  }
  
  function editRule(event: CustomEvent<MonitorRule>) {
    editingRule = event.detail
    showDialog = true
  }
  
  async function saveRule(event: CustomEvent<MonitorRule>) {
    try {
      const rule = event.detail
      if (editingRule?.id) {
        await rulesStore.updateRule(rule)
      } else {
        await rulesStore.addRule(rule)
      }
      showDialog = false
    } catch (error) {
      console.error('Failed to save rule:', error)
      alert('保存规则失败: ' + error)
    }
  }
  
  function closeDialog() {
    showDialog = false
    editingRule = null
  }
  
  async function handleCheckNow(event: CustomEvent<string>) {
    try {
      await rulesStore.checkNow(event.detail)
      alert('检查完成！')
    } catch (error) {
      console.error('Failed to check now:', error)
      alert('立即检查失败: ' + error + '\n\n提示：请确保规则已启动监控')
    }
  }
  
  onMount(() => {
    rulesStore.loadRules()
    rulesStore.setupEventListeners()
  })
  
  let rules = $derived($rulesStore.rules)
  let selectedRuleId = $derived($rulesStore.selectedRuleId)
</script>

<div class="app-container">
  <Toolbar 
    on:add-rule={showAddDialog} 
    on:refresh={() => rulesStore.loadRules()} 
  />
  <div class="main-content">
    <RuleList 
      {rules}
      selectedId={selectedRuleId}
      on:select={(e) => rulesStore.selectRule(e.detail)}
      on:edit={editRule}
      on:delete={(e) => rulesStore.deleteRule(e.detail)}
      on:toggle={(e) => rulesStore.toggleRule(e.detail)}
    />
    {#if $selectedRule}
      <RuleDetail 
        rule={$selectedRule}
        on:check-now={handleCheckNow}
      />
    {:else}
      <div class="no-selection">
        <p>请选择一个规则查看详情</p>
      </div>
    {/if}
  </div>
  {#if showDialog}
    <RuleDialog 
      rule={editingRule}
      on:save={saveRule}
      on:cancel={closeDialog}
    />
  {/if}
</div>

<style>
  /* 主题变量 */
  :global(:root[data-theme="light"]) {
    --bg-primary: #ffffff;
    --bg-secondary: #f5f5f5;
    --bg-tertiary: #f9f9f9;
    --bg-hover: #f0f0f0;
    --bg-selected: #e3f2fd;
    --text-primary: #333333;
    --text-secondary: #666666;
    --text-tertiary: #999999;
    --text-inverse: #ffffff;
    --border-primary: #dddddd;
    --border-secondary: #eeeeee;
    --color-primary: #007bff;
    --color-primary-hover: #0056b3;
    --color-secondary: #6c757d;
    --color-secondary-hover: #545b62;
    --color-success: #28a745;
    --color-success-bg: #d4edda;
    --color-success-text: #155724;
    --color-warning: #ffc107;
    --color-warning-bg: #fff3cd;
    --color-warning-text: #856404;
    --color-error: #dc3545;
    --color-error-bg: #f8d7da;
    --color-error-text: #721c24;
    --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.1);
    --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
    --shadow-lg: 0 10px 20px rgba(0, 0, 0, 0.15);
  }

  :global(:root[data-theme="dark"]) {
    --bg-primary: #1e1e1e;
    --bg-secondary: #252525;
    --bg-tertiary: #2d2d2d;
    --bg-hover: #333333;
    --bg-selected: #1a3a52;
    --text-primary: #e0e0e0;
    --text-secondary: #b0b0b0;
    --text-tertiary: #808080;
    --text-inverse: #1e1e1e;
    --border-primary: #404040;
    --border-secondary: #333333;
    --color-primary: #4a9eff;
    --color-primary-hover: #6bb0ff;
    --color-secondary: #8a9099;
    --color-secondary-hover: #9fa5ad;
    --color-success: #3dbd5d;
    --color-success-bg: #1e3a28;
    --color-success-text: #7ce89d;
    --color-warning: #f0ad4e;
    --color-warning-bg: #3d3020;
    --color-warning-text: #f5d89f;
    --color-error: #e74c3c;
    --color-error-bg: #3d2020;
    --color-error-text: #f5a9a0;
    --shadow-sm: 0 1px 3px rgba(0, 0, 0, 0.3);
    --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.3);
    --shadow-lg: 0 10px 20px rgba(0, 0, 0, 0.4);
  }

  :global(body) {
    margin: 0;
    padding: 0;
    font-family: -apple-system, BlinkMacSystemFont, 'Segoe UI', 'Roboto', 'Oxygen',
      'Ubuntu', 'Cantarell', 'Fira Sans', 'Droid Sans', 'Helvetica Neue',
      sans-serif;
    -webkit-font-smoothing: antialiased;
    -moz-osx-font-smoothing: grayscale;
    background-color: var(--bg-secondary);
    color: var(--text-primary);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  :global(*) {
    box-sizing: border-box;
  }
  
  .app-container {
    display: flex;
    flex-direction: column;
    height: 100vh;
    background: var(--bg-secondary);
    transition: background-color 0.3s ease;
  }
  
  .main-content {
    display: flex;
    flex: 1;
    overflow: hidden;
  }
  
  .main-content > :global(*) {
    flex: 1;
  }
  
  .no-selection {
    display: flex;
    align-items: center;
    justify-content: center;
    background: var(--bg-primary);
    color: var(--text-tertiary);
    transition: background-color 0.3s ease, color 0.3s ease;
  }
  
  .no-selection p {
    font-size: 16px;
  }
</style>
