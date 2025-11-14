# API Watch - 主题切换功能

## 🎨 功能概述

API Watch现在支持**亮色主题**和**暗色主题**两种显示模式，用户可以根据个人喜好和使用环境自由切换。

## ✨ 主题特性

### 亮色主题（Light Theme）
- **背景色**: 白色和浅灰色系
- **文字色**: 深灰色和黑色系
- **适用场景**: 明亮环境、白天使用
- **视觉特点**: 清新明亮，对比度适中

### 暗色主题（Dark Theme）
- **背景色**: 深灰色和黑色系
- **文字色**: 浅灰色和白色系
- **适用场景**: 暗光环境、夜间使用
- **视觉特点**: 护眼舒适，减少眼睛疲劳

## 🔄 如何切换主题

### 方法1：使用工具栏按钮
1. 打开API Watch应用
2. 在顶部工具栏右侧找到主题切换按钮
3. 点击按钮即可切换：
   - 🌙 暗色 - 切换到暗色主题
   - ☀️ 亮色 - 切换到亮色主题

### 主题持久化
- 主题选择会自动保存到浏览器localStorage
- 下次打开应用时会自动应用上次选择的主题
- 无需重复设置

## 🎨 主题设计细节

### 颜色变量系统

应用使用CSS变量实现主题系统，所有颜色都通过变量定义：

#### 背景色
- `--bg-primary`: 主背景色
- `--bg-secondary`: 次要背景色
- `--bg-tertiary`: 第三级背景色
- `--bg-hover`: 悬停背景色
- `--bg-selected`: 选中背景色

#### 文字色
- `--text-primary`: 主文字色
- `--text-secondary`: 次要文字色
- `--text-tertiary`: 第三级文字色
- `--text-inverse`: 反色文字（用于按钮等）

#### 边框色
- `--border-primary`: 主边框色
- `--border-secondary`: 次要边框色

#### 品牌色
- `--color-primary`: 主品牌色（蓝色）
- `--color-secondary`: 次要品牌色（灰色）

#### 状态色
- `--color-success`: 成功状态（绿色）
- `--color-warning`: 警告状态（黄色）
- `--color-error`: 错误状态（红色）

### 过渡动画

所有主题切换都包含平滑的过渡动画：
- 背景色过渡：0.3秒
- 文字色过渡：0.3秒
- 边框色过渡：0.3秒

## 📱 组件适配

所有UI组件都已完全适配主题系统：

### ✅ 已适配组件
- **工具栏（Toolbar）**: 包含主题切换按钮
- **规则列表（RuleList）**: 列表项、状态徽章
- **规则详情（RuleDetail）**: 详情面板、代码块
- **规则对话框（RuleDialog）**: 表单、输入框、按钮
- **主应用（App）**: 整体布局、空状态

### 适配特点
- 所有文字在两种主题下都清晰可读
- 状态徽章在暗色主题下使用更柔和的颜色
- 边框和分隔线在暗色主题下更加明显
- 代码块和输入框在两种主题下都有良好的对比度

## 🛠️ 技术实现

### 前端架构
```
frontend/src/
├── stores/
│   └── theme.ts          # 主题状态管理
├── App.svelte            # 主题变量定义
└── components/
    ├── Toolbar.svelte    # 主题切换按钮
    ├── RuleList.svelte   # 使用主题变量
    ├── RuleDetail.svelte # 使用主题变量
    └── RuleDialog.svelte # 使用主题变量
```

### 状态管理
使用Svelte stores管理主题状态：
```typescript
export const theme = writable<Theme>('light')
export function toggleTheme() {
  theme.update(current => current === 'light' ? 'dark' : 'light')
}
```

### 持久化
主题选择保存在localStorage：
```typescript
localStorage.setItem('theme', value)
```

### DOM属性
通过data-theme属性控制主题：
```typescript
document.documentElement.setAttribute('data-theme', value)
```

## 🎯 最佳实践

### 使用建议
1. **白天使用**: 推荐使用亮色主题
2. **夜间使用**: 推荐使用暗色主题
3. **长时间使用**: 暗色主题可减少眼睛疲劳
4. **演示展示**: 亮色主题对比度更高，更适合演示

### 可访问性
- 两种主题都符合WCAG 2.1 AA级对比度标准
- 文字和背景对比度 ≥ 4.5:1
- 状态色在两种主题下都清晰可辨

## 🚀 未来计划

### 可能的增强功能
1. **自动切换**: 根据系统时间自动切换主题
2. **跟随系统**: 跟随操作系统的主题设置
3. **自定义主题**: 允许用户自定义颜色
4. **更多主题**: 添加更多预设主题（如高对比度主题）

## 📝 开发者指南

### 添加新组件时的主题适配

1. **使用CSS变量**：
```css
.my-component {
  background: var(--bg-primary);
  color: var(--text-primary);
  border: 1px solid var(--border-primary);
}
```

2. **添加过渡动画**：
```css
.my-component {
  transition: background-color 0.3s ease, color 0.3s ease;
}
```

3. **测试两种主题**：
- 在亮色主题下测试
- 在暗色主题下测试
- 确保文字清晰可读
- 确保交互元素明显可见

### 可用的CSS变量

参考`App.svelte`中的`:global(:root[data-theme="light"])`和`:global(:root[data-theme="dark"])`部分，查看所有可用的CSS变量。

---

**更新时间**: 2025-11-13  
**版本**: 1.0.0  
**状态**: ✅ 已完成并测试
