import { writable, derived } from 'svelte/store'
import type { MonitorRule } from '../types/models'
import { GetRules, AddRule, UpdateRule, DeleteRule, StartMonitoring, StopMonitoring, CheckNow } from '../wailsjs/wailsjs/go/main/App'
import { EventsOn } from '../wailsjs/wailsjs/runtime/runtime'

interface RulesState {
  rules: MonitorRule[]
  selectedRuleId: string | null
}

function createRulesStore() {
  const { subscribe, set, update } = writable<RulesState>({
    rules: [],
    selectedRuleId: null
  })
  
  return {
    subscribe,
    
    async loadRules() {
      try {
        const rules = await GetRules()
        update(state => ({ ...state, rules }))
      } catch (error) {
        console.error('Failed to load rules:', error)
      }
    },
    
    async addRule(rule: MonitorRule) {
      try {
        await AddRule(rule)
        await this.loadRules()
      } catch (error) {
        console.error('Failed to add rule:', error)
        throw error
      }
    },
    
    async updateRule(rule: MonitorRule) {
      try {
        await UpdateRule(rule)
        await this.loadRules()
      } catch (error) {
        console.error('Failed to update rule:', error)
        throw error
      }
    },
    
    async deleteRule(id: string) {
      try {
        await DeleteRule(id)
        update(state => ({
          ...state,
          selectedRuleId: state.selectedRuleId === id ? null : state.selectedRuleId
        }))
        await this.loadRules()
      } catch (error) {
        console.error('Failed to delete rule:', error)
        throw error
      }
    },
    
    async toggleRule(id: string) {
      try {
        let currentState: RulesState | null = null
        const unsubscribe = subscribe(s => {
          currentState = s
          unsubscribe()
        })
        
        if (!currentState) return
        
        const rule = currentState.rules.find((r: MonitorRule) => r.id === id)
        if (rule) {
          if (rule.enabled) {
            await StopMonitoring(id)
          } else {
            await StartMonitoring(id)
          }
          await this.loadRules()
        }
      } catch (error) {
        console.error('Failed to toggle rule:', error)
        throw error
      }
    },
    
    async checkNow(id: string) {
      try {
        await CheckNow(id)
        // 检查成功后刷新规则列表
        await this.loadRules()
      } catch (error) {
        console.error('Failed to check now:', error)
        throw error
      }
    },
    
    selectRule(id: string) {
      update(state => ({ ...state, selectedRuleId: id }))
    },
    
    setupEventListeners() {
      // 监听后端事件
      EventsOn('rule_updated', () => {
        this.loadRules()
      })
      
      EventsOn('content_changed', (event: any) => {
        console.log('Content changed:', event)
        this.loadRules()
      })
      
      EventsOn('rule_status_changed', () => {
        this.loadRules()
      })
    }
  }
}

export const rulesStore = createRulesStore()

// 派生store：获取当前选中的规则
export const selectedRule = derived(
  rulesStore,
  $rulesStore => $rulesStore.rules.find(r => r.id === $rulesStore.selectedRuleId) || null
)
