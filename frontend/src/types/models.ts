export type ExtractorType = 'css' | 'regex' | 'json'

export type RuleStatus = 'running' | 'paused' | 'error' | 'idle'

export interface MonitorRule {
  id: string
  name: string
  description: string
  url: string
  method: string
  headers?: Record<string, string>
  body?: string
  interval: string
  extractor_type: ExtractorType
  extractor_expr: string
  notify_enabled: boolean
  enabled: boolean
  last_content: string
  last_checked: string
  status: RuleStatus
  error_message?: string
}

export interface Event {
  type: string
  rule_id: string
  rule?: MonitorRule
  timestamp: string
  data?: any
}
