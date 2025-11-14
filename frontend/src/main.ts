import { mount } from 'svelte'
import App from './App.svelte'
import { theme } from './stores/theme'

// 初始化主题
const storedTheme = localStorage.getItem('theme') as 'light' | 'dark' | null
if (storedTheme) {
  theme.set(storedTheme)
  document.documentElement.setAttribute('data-theme', storedTheme)
} else {
  document.documentElement.setAttribute('data-theme', 'light')
}

// @ts-ignore - Svelte 5 component type compatibility issue
const app = mount(App, {
  target: document.getElementById('app')!,
})

export default app
