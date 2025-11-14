# 构建说明

## Windows构建问题解决

Fyne需要CGO支持，在Windows上需要MinGW-w64。

### 方法1：使用MSYS2（推荐）

1. 下载并安装MSYS2: https://www.msys2.org/

2. 打开MSYS2 MINGW64终端，安装gcc：
   ```bash
   pacman -S mingw-w64-x86_64-gcc
   ```

3. 将MinGW64的bin目录添加到系统PATH：
   ```
   C:\msys64\mingw64\bin
   ```

4. 重新打开PowerShell或CMD，验证gcc：
   ```bash
   gcc --version
   ```

5. 编译项目：
   ```bash
   go build -o url-monitor.exe
   ```

### 方法2：使用TDM-GCC

1. 下载TDM-GCC: https://jmeubank.github.io/tdm-gcc/

2. 安装时选择64位版本

3. 添加到PATH并编译

### 方法3：使用fyne命令行工具

```bash
# 安装fyne工具
go install fyne.io/fyne/v2/cmd/fyne@latest

# 使用fyne打包
fyne package -os windows -icon icon.png
```

## 无GUI版本

如果不需要GUI，可以注释掉main.go中的UI相关代码，使用命令行版本：

```go
// 注释掉UI相关代码
// adapter := fyneui.NewFyneUIAdapter(engine)
// adapter.Run()

// 使用信号等待
sigCh := make(chan os.Signal, 1)
signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
<-sigCh
```

## 跨平台编译

### 编译Linux版本（在Windows上）

```bash
set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=1
go build -o url-monitor-linux
```

### 编译macOS版本（在Windows上）

需要OSXCross工具链，较为复杂，建议在目标平台上编译。

## 测试

测试不需要CGO，可以直接运行：

```bash
go test ./...
```

## 依赖管理

```bash
# 更新依赖
go get -u ./...

# 清理依赖
go mod tidy

# 查看依赖树
go mod graph
```

## 性能优化

### 减小二进制大小

```bash
go build -ldflags="-s -w" -o url-monitor.exe
```

### 使用UPX压缩

```bash
upx --best url-monitor.exe
```

## 调试

### 启用详细日志

修改`logger/logger.go`中的日志级别：

```go
Level: slog.LevelDebug,
```

### 使用delve调试

```bash
go install github.com/go-delve/delve/cmd/dlv@latest
dlv debug
```

## 常见问题

### Q: 编译时提示"build constraints exclude all Go files"

A: 这是因为缺少CGO支持，请按照上述方法安装MinGW-w64。

### Q: 运行时提示"找不到DLL"

A: 确保MinGW的bin目录在PATH中，或将相关DLL复制到程序目录。

### Q: 如何禁用GUI？

A: 修改main.go，注释掉Fyne相关代码，使用命令行模式。

### Q: 测试失败

A: 确保所有依赖已安装：`go mod download`

## 持续集成

### GitHub Actions示例

```yaml
name: Build

on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go test ./...

  build:
    runs-on: ${{ matrix.os }}
    strategy:
      matrix:
        os: [ubuntu-latest, windows-latest, macos-latest]
    steps:
      - uses: actions/checkout@v2
      - uses: actions/setup-go@v2
        with:
          go-version: 1.21
      - run: go build
```
