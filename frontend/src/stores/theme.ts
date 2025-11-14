import { writable } from 'svelte/store'

export type Theme = 'light' | 'dark'

// 从localStorage读取主题，默认为light
const storedTheme = (typeof localStorage !== 'undefined' 
  ? localStorage.getItem('theme') 
  : null) as Theme | null

export const theme = writable<Theme>(storedTheme || 'light')

// 订阅主题变化，保存到localStorage并应用到document
theme.subscribe(value => {
  if (typeof localStorage !== 'undefined') {
    localStorage.setItem('theme', value)
  }
  if (typeof document !== 'undefined') {
    document.documentElement.setAttribute('data-theme', value)
  }
})

export function toggleTheme() {
  theme.update(current => current === 'light' ? 'dark' : 'light')
}
