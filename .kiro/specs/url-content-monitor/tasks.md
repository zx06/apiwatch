# 实现计划

- [x] 1. 初始化Wails项目结构和依赖


  - 使用wails init创建项目基础结构
  - 创建核心层目录：models、config、extractor、fetcher、monitor、notification、core
  - 保留现有核心层代码（已实现）
  - 清理Fyne相关文件和依赖
  - 添加Wails依赖（github.com/wailsapp/wails/v2）
  - 初始化前端项目（Svelte + TypeScript + Vite）
  - 配置wails.json项目配置文件
  - _需求: 1.1, 1.5_



- [x] 2. 实现数据模型和基础类型

  - 在models/rule.go中定义MonitorRule结构体
  - 添加Description字段（规则描述）
  - 将Interval从int改为time.Duration类型
  - 添加Method字段（HTTP方法）
  - 添加Headers字段（map[string]string，自定义请求头）
  - 添加Body字段（请求体）
  - 定义ExtractorType和RuleStatus枚举类型
  - 实现Validate方法（URL格式、间隔范围、方法有效性等）
  - 添加YAML和JSON标签支持




  - _需求: 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 2.4_


- [ ] 3. 实现配置管理器
  - [x] 3.1 创建配置管理器接口和实现

    - 在config/manager.go中定义Manager接口
    - 实现YAMLManager结构体，支持Load和Save方法
    - 实现配置文件路径管理（~/.url-monitor/config.yaml）
    - _需求: 1.4, 1.5_



  

  - [ ] 3.2 实现规则的CRUD操作
    - 实现AddRule、UpdateRule、DeleteRule、GetRule方法
    - 添加文件锁机制防止并发写入
    - 实现配置文件的自动创建和初始化

    - _需求: 1.1, 6.1, 6.2, 6.3, 6.4_

- [ ] 4. 实现HTTP客户端
  - 在fetcher/fetcher.go中定义Request、Response和Fetcher接口

  - Request结构包含URL、Method、Headers、Body字段
  - 实现HTTPFetcher，配置30秒超时
  - 支持多种HTTP方法（GET、POST、PUT、DELETE、PATCH等）
  - 支持自定义请求头
  - 支持请求体（用于POST/PUT等）
  - 设置默认User-Agent
  - 实现重试机制（最多3次，指数退避）
  - 处理HTTP重定向和常见错误
  - _需求: 1.5, 1.6, 1.7, 3.1, 3.3_

- [x] 5. 实现内容提取器

  - [x] 5.1 创建提取器接口和工厂


    - 在extractor/extractor.go中定义Extractor接口
    - 实现Factory工厂类，根据类型创建提取器
    - _需求: 2.1, 2.2, 2.3, 2.4_
  
  - [x] 5.2 实现CSS选择器提取器


    - 创建CSSExtractor，使用goquery解析HTML
    - 实现CSS选择器匹配和文本提取
    - 处理选择器无匹配的情况
    - _需求: 2.1, 2.5_
  
  - [x] 5.3 实现正则表达式提取器


    - 创建RegexExtractor，使用regexp包
    - 预编译正则表达式以提高性能
    - 设置匹配超时防止ReDoS攻击
    - _需求: 2.2, 2.5_
  
  - [x] 5.4 实现JSON路径提取器


    - 创建JSONExtractor，使用gjson库
    - 实现JSON路径查询
    - 处理无效JSON和路径不存在的情况
    - _需求: 2.3, 2.5_



- [ ] 6. 实现Wails通知服务
  - 在notification/notifier.go中保留Notifier接口（已存在）
  - 移除FyneNotifier实现
  - 创建WailsNotifier，使用Wails runtime发送系统通知
  - 限制通知消息长度（最多200字符）
  - 处理通知发送失败的情况
  - _需求: 5.1, 5.2, 5.3_

- [x] 7. 实现监控任务


  - [x] 7.1 创建Task结构和基础方法


    - 在monitor/task.go中定义Task结构体
    - 实现Start方法，创建ticker和goroutine
    - 实现Stop方法，停止ticker和goroutine
    - 使用channel (stopCh) 实现优雅退出
    - _需求: 3.1, 3.2, 3.4_
  
  - [x] 7.2 实现任务执行逻辑

    - 实现RunOnce方法，执行单次检查
    - 集成Fetcher获取URL内容
    - 集成Extractor提取内容
    - 实现内容变化检测（与LastContent比较）
    - _需求: 3.1, 4.1, 4.5, 5.1_
  
  - [x] 7.3 实现通知和状态更新

    - 在内容变化时调用Notifier发送通知
    - 实现首次提取不通知的逻辑
    - 更新规则的LastContent、LastChecked和Status
    - 处理错误并更新ErrorMessage
    - _需求: 5.1, 5.4, 5.5, 4.2, 4.3_
  
  - [x] 7.4 实现任务更新和错误处理

    - 实现Update方法，支持运行时更新配置
    - 实现网络错误的自动重试逻辑
    - 记录错误日志
    - _需求: 3.3, 6.5_

- [x] 8. 实现监控服务


  - 在monitor/service.go中定义Service接口
  - 实现MonitorService，使用map存储活跃任务
  - 使用sync.RWMutex保护并发访问
  - 实现StartTask、StopTask、StopAll方法
  - 实现UpdateTask和RunTaskOnce方法
  - 实现GetTaskStatus方法
  - _需求: 3.1, 3.2, 3.4, 3.5, 4.5_






- [ ] 9. 实现核心引擎和API接口
  - [ ] 9.1 定义核心API接口
    - 在core/api.go中定义CoreAPI接口
    - 定义规则管理方法（GetRules、AddRule、UpdateRule、DeleteRule）
    - 定义监控控制方法（StartMonitoring、StopMonitoring、CheckNow）



    - 定义事件订阅方法（Subscribe、Unsubscribe）
    - _需求: 1.5, 6.1, 6.2_
  
  - [ ] 9.2 实现事件总线
    - 在core/event.go中定义Event、EventType、EventListener接口

    - 实现EventBus，支持事件发布和订阅
    - 使用goroutine异步分发事件，避免阻塞核心逻辑
    - 使用sync.RWMutex保护监听器列表
    - _需求: 4.1, 4.2, 5.1_
  
  - [x] 9.3 实现Monitor Engine

    - 在core/engine.go中实现Engine结构体
    - 实现CoreAPI接口的所有方法
    - 集成配置管理器、监控服务、事件总线
    - 实现Initialize和Shutdown方法
    - 在状态变化时发布事件到EventBus


    - _需求: 1.5, 3.5, 6.5_

- [ ] 10. 实现Wails应用结构
  - [x] 10.1 创建Wails应用入口

    - 创建app.go，定义App结构体
    - 实现startup、shutdown、domReady生命周期方法
    - 持有CoreAPI和EventBus引用
    - _需求: 4.1, 4.2_
  
  - [ ] 10.2 实现Go绑定方法
    - 创建bindings.go，实现前端可调用的方法

    - 实现GetRules、GetRule、AddRule、UpdateRule、DeleteRule
    - 实现StartMonitoring、StopMonitoring、CheckNow
    - 实现事件推送到前端的机制（使用runtime.EventsEmit）
    - _需求: 1.1, 6.1, 6.2, 4.5_

- [ ] 11. 实现前端基础结构
  - [x] 11.1 初始化Svelte项目

    - 配置Vite和TypeScript
    - 安装Svelte和相关插件
    - 配置Wails runtime绑定
    - 创建基础目录结构（components、stores、types）
    - 配置svelte.config.js和tsconfig.json


    - _需求: 4.1_
  
  - [ ] 11.2 定义TypeScript类型
    - 在types/models.ts中定义MonitorRule接口
    - 定义ExtractorType和RuleStatus枚举
    - 定义Event接口
    - _需求: 1.2, 1.3_

- [x] 12. 实现前端状态管理


  - 创建stores/rules.ts，使用Svelte stores管理规则状态
  - 实现writable store存储rules和selectedRuleId
  - 实现loadRules、addRule、updateRule、deleteRule方法
  - 实现toggleRule和checkNow方法


  - 实现selectRule方法
  - 实现setupEventListeners，监听后端事件
  - 创建derived store获取selectedRule
  - _需求: 4.1, 4.2, 6.3, 6.4_

- [ ] 13. 实现前端UI组件
  - [ ] 13.1 实现工具栏组件
    - 创建components/Toolbar.svelte
    - 实现添加规则按钮（dispatch add-rule事件）
    - 实现刷新按钮（dispatch refresh事件）
    - _需求: 1.1, 6.1_


  
  - [ ] 13.2 实现规则列表组件
    - 创建components/RuleList.svelte
    - 接收rules和selectedId作为props
    - 使用{#each}循环显示规则列表
    - 显示规则名称、状态徽章、URL、最后检查时间
    - 实现列表项选择（dispatch select事件）
    - 实现启动/暂停按钮（dispatch toggle事件）
    - 实现编辑按钮（dispatch edit事件）


    - 实现删除按钮（dispatch delete事件）
    - 添加状态样式（running、paused、error、idle）
    - 使用class:selected实现选中样式
    - _需求: 4.1, 4.2, 4.3_
  
  - [ ] 13.3 实现规则详情组件
    - 创建components/RuleDetail.svelte
    - 接收rule作为prop
    - 显示规则完整信息（名称、描述、URL、方法、请求头、请求体）
    - 显示提取器配置（类型、表达式）
    - 显示最新提取内容（可滚动文本区域）
    - 显示最后检查时间和状态


    - 使用{#if}条件显示错误信息
    - 实现"立即检查"按钮（dispatch check-now事件）
    - _需求: 1.3, 1.5, 1.6, 1.7, 4.1, 4.2, 4.3, 4.4, 4.5_
  
  - [ ] 13.4 实现规则编辑对话框组件
    - 创建components/RuleDialog.svelte
    - 接收rule作为prop（null表示新建）
    - 使用bind:value实现双向绑定

    - 实现表单字段：名称、描述、URL、HTTP方法、请求头、请求体
    - 实现检查间隔输入（支持duration格式）
    - 实现提取器配置（类型选择、表达式输入）
    - 使用bind:checked实现复选框（启用通知、启用规则）
    - 实现表单验证（URL格式、duration格式、必填字段）
    - 使用响应式语句根据HTTP方法显示/隐藏请求体字段
    - 实现保存按钮（dispatch save事件）
    - 实现取消按钮（dispatch cancel事件）
    - _需求: 1.1, 1.2, 1.3, 1.4, 1.5, 1.6, 1.7, 2.4, 5.4, 6.1_

- [ ] 14. 实现主应用组件
  - 创建App.svelte，实现主布局
  - 导入并集成Toolbar、RuleList、RuleDetail、RuleDialog组件
  - 使用$rulesStore访问store状态
  - 实现左右分栏布局（规则列表和详情面板）
  - 使用{#if}条件渲染RuleDetail和RuleDialog
  - 实现对话框显示/隐藏逻辑
  - 在onMount时加载规则和设置事件监听
  - _需求: 4.1, 4.4_

- [x] 15. 实现样式和主题


  - 创建全局CSS样式
  - 实现响应式布局
  - 实现状态徽章样式（running、paused、error、idle）
  - 实现按钮和表单样式
  - 实现对话框遮罩和动画
  - _需求: 4.1_

- [x] 16. 实现日志系统


  - 创建logger包，封装标准库log/slog
  - 配置slog使用JSON格式输出

  - 实现日志文件管理（~/.url-monitor/logs/app.log）
  - 配置日志级别（Debug、Info、Warn、Error）
  - 实现日志轮转（每个文件最大10MB，保留5个文件）
  - 配置多输出（文件和控制台）
  - 在核心层关键位置添加结构化日志记录
  - _需求: 2.5, 3.3_



- [ ] 16. 实现Wails主程序入口
  - 在main.go中初始化核心层组件
  - 创建配置管理器、监控服务
  - 创建并初始化Monitor Engine
  - 创建Wails App实例并绑定方法
  - 配置Wails应用选项（窗口大小、标题、前端资源路径）


  - 使用wails.Run启动应用
  - 在startup中加载配置并启动已启用的监控任务
  - 在shutdown中实现优雅退出（调用Engine.Shutdown）
  - _需求: 1.5, 3.5, 7.1, 7.2_

- [x] 17. 清理Fyne相关文件



  - 删除ui/fyne目录及所有文件
  - 删除FyneApp.toml配置文件
  - 删除fyne_metadata_init.go
  - 从go.mod中移除Fyne依赖
  - 更新README和文档，移除Fyne相关说明
  - 删除Fyne相关的构建脚本和配置
  - _需求: 1.1_

- [ ] 18. 完善错误处理和边界情况
  - 处理配置文件不存在或损坏的情况
  - 处理网络不可用的情况
  - 处理无效的提取表达式
  - 处理并发修改冲突
  - 在前端添加用户友好的错误提示（Toast或Dialog）
  - 确保核心层错误通过事件传递到前端
  - _需求: 2.5, 3.3_

- [ ] 19. 测试和优化
  - 测试Wails应用的构建和打包
  - 测试前后端通信的稳定性
  - 测试事件推送机制
  - 优化前端性能（虚拟滚动、懒加载）
  - 测试多任务并发性能
  - 测试跨平台兼容性（Windows、macOS、Linux）
  - _需求: 3.2_
