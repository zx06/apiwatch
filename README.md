# API Watch

一个基于Go语言和Wails开发的跨平台桌面应用程序，用于监控URL内容变化并发送通知。

## 功能特性

- ✅ 支持多个监控规则同时运行
- ✅ 支持多种HTTP方法（GET、POST、PUT、DELETE等）
- ✅ 自定义HTTP请求头和请求体
- ✅ 三种内容提取方式：CSS选择器、正则表达式、JSON路径
- ✅ 灵活的检查间隔配置（支持duration格式：5m、1h等）
- ✅ 内容变化时系统通知
- ✅ 完整的规则管理（增删改查）
- ✅ 实时状态监控
- ✅ 核心逻辑与UI完全解耦
- ✅ 基于事件驱动的架构

## 技术栈

- **语言**: Go 1.21+
- **UI框架**: Wails v2
- **前端**: Svelte + TypeScript + Vite
- **配置**: YAML
- **日志**: log/slog
- **测试**: testify

## 项目结构

```
.
├── models/          # 数据模型
├── config/          # 配置管理
├── fetcher/         # HTTP客户端
├── extractor/       # 内容提取器
├── monitor/         # 监控任务和服务
├── notification/    # 通知服务
├── core/            # 核心引擎和API
├── logger/          # 日志系统
├── frontend/        # Svelte前端（开发中）
└── main.go          # 主程序入口
```

## 构建和运行

### 前置要求

#### 安装Wails CLI
```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

#### 平台要求

**Windows**
- Go 1.21+
- Node.js 16+
- WebView2 (Windows 10/11自带)

**macOS**
- Go 1.21+
- Node.js 16+
- Xcode Command Line Tools

**Linux**
- Go 1.21+
- Node.js 16+
- WebKit2GTK

### 开发模式

```bash
# 启动开发服务器（热重载）
wails dev
```

### 编译

```bash
# 编译生产版本
wails build

# 指定平台编译
wails build -platform windows/amd64
wails build -platform darwin/universal
wails build -platform linux/amd64
```

### 运行应用

```bash
# 直接运行编译好的应用
./build/bin/apiwatch.exe  # Windows
./build/bin/apiwatch       # macOS/Linux
```

## 配置文件

配置文件位置：`~/.url-monitor/config.yaml`

示例配置：

```yaml
version: "1.0"
rules:
  - id: uuid-1
    name: 示例监控
    description: 监控API响应
    url: https://api.example.com/data
    method: POST
    headers:
      Authorization: Bearer token123
      Content-Type: application/json
    body: '{"query": "latest"}'
    interval: 5m
    extractor_type: json
    extractor_expr: data.items[0].title
    notify_enabled: true
    enabled: true
```

## 日志

日志文件位置：`~/.url-monitor/logs/app.log`

- 格式：JSON
- 级别：Debug、Info、Warn、Error
- 自动轮转：每个文件最大10MB，保留5个备份

## 测试

```bash
# 运行所有测试
go test ./...

# 运行特定包的测试
go test ./models
go test ./extractor

# 查看测试覆盖率
go test ./... -cover
```

## 开发

### 添加新的提取器

1. 在`extractor/`目录创建新文件
2. 实现`Extractor`接口
3. 在`Factory.Create()`中添加新类型
4. 添加测试

### 添加新的UI实现

核心层与UI完全解耦，可以轻松添加新的UI实现：

1. 实现Go绑定方法（类似Wails的App结构）
2. 通过`CoreAPI`与核心层交互
3. 监听`EventBus`获取状态更新
4. 可选择Web、CLI、移动端等任何UI技术栈

## 架构设计

### 核心原则

1. **核心逻辑独立**: 所有业务逻辑不依赖UI框架
2. **接口驱动**: 通过清晰的接口实现松耦合
3. **事件驱动**: 使用事件总线实现核心层向UI层的通知
4. **可扩展性**: 支持多种UI实现和部署模式

### 数据流

```
UI层 → CoreAPI → Engine → Monitor Service → Task
                    ↓
                EventBus → UI层（状态更新）
```

## 许可证

MIT License

## 贡献

欢迎提交Issue和Pull Request！

## 作者

zx06
