# 应用图标说明

## 可用图标

项目包含两个SVG图标文件：

### 1. icon.svg
- 标准版本图标
- 包含监控屏幕、链接符号和信号波
- 带有红色通知点
- 适合用作应用图标

### 2. icon-simple.svg
- 简化版本图标
- 包含地球仪（代表URL/网络）和心跳线（代表监控）
- 带有绿色活动指示器和脉冲动画
- 更现代的渐变设计

## 如何使用

### 方法1：使用Fyne打包工具

```bash
# 安装fyne工具
go install fyne.io/fyne/v2/cmd/fyne@latest

# 使用图标打包应用
fyne package -os windows -icon icon.svg
fyne package -os darwin -icon icon.svg
fyne package -os linux -icon icon.svg
```

### 方法2：转换为PNG

如果需要PNG格式，可以使用以下工具转换：

#### 使用ImageMagick
```bash
# 安装ImageMagick
# Windows: choco install imagemagick
# macOS: brew install imagemagick
# Linux: apt-get install imagemagick

# 转换为不同尺寸
convert icon.svg -resize 512x512 icon-512.png
convert icon.svg -resize 256x256 icon-256.png
convert icon.svg -resize 128x128 icon-128.png
convert icon.svg -resize 64x64 icon-64.png
convert icon.svg -resize 32x32 icon-32.png
convert icon.svg -resize 16x16 icon-16.png
```

#### 使用在线工具
- https://cloudconvert.com/svg-to-png
- https://convertio.co/svg-png/
- https://www.aconvert.com/image/svg-to-png/

### 方法3：在代码中使用

在Fyne应用中设置图标：

```go
import (
    "fyne.io/fyne/v2/app"
)

func main() {
    myApp := app.NewWithID("com.github.zx06.apiwatch")
    
    // 设置应用图标
    icon, err := fyne.LoadResourceFromPath("icon.png")
    if err == nil {
        myApp.SetIcon(icon)
    }
    
    // ... 其他代码
}
```

## 图标设计说明

### icon.svg 设计元素
- **蓝色圆形背景**: 代表应用的主题色
- **白色屏幕**: 代表监控界面
- **链接符号**: 代表URL/网络连接
- **信号波**: 代表持续监控
- **红色通知点**: 代表内容变化提醒

### icon-simple.svg 设计元素
- **渐变背景**: 现代化的视觉效果
- **地球仪**: 代表互联网/URL
- **心跳线**: 代表实时监控
- **绿色脉冲点**: 代表活动状态
- **动画效果**: SVG包含脉冲动画

## 自定义图标

如果你想自定义图标，可以：

1. 使用矢量图形编辑器（如Inkscape、Adobe Illustrator）编辑SVG文件
2. 修改颜色：
   - 主色：`#4A90E2` (蓝色)
   - 深色：`#2E5C8A` (深蓝)
   - 通知色：`#FF5252` (红色)
   - 活动色：`#48BB78` (绿色)

3. 调整尺寸：SVG是矢量格式，可以无损缩放到任意大小

## Windows图标（.ico）

Windows应用需要.ico格式的图标，可以使用以下工具转换：

```bash
# 使用ImageMagick创建多尺寸ico
convert icon.svg -define icon:auto-resize=256,128,64,48,32,16 icon.ico
```

或使用在线工具：
- https://convertio.co/svg-ico/
- https://www.icoconverter.com/

## macOS图标（.icns）

macOS应用需要.icns格式，Fyne打包工具会自动处理，或手动创建：

```bash
# 创建iconset目录
mkdir icon.iconset

# 生成不同尺寸
sips -z 16 16     icon.svg --out icon.iconset/icon_16x16.png
sips -z 32 32     icon.svg --out icon.iconset/icon_16x16@2x.png
sips -z 32 32     icon.svg --out icon.iconset/icon_32x32.png
sips -z 64 64     icon.svg --out icon.iconset/icon_32x32@2x.png
sips -z 128 128   icon.svg --out icon.iconset/icon_128x128.png
sips -z 256 256   icon.svg --out icon.iconset/icon_128x128@2x.png
sips -z 256 256   icon.svg --out icon.iconset/icon_256x256.png
sips -z 512 512   icon.svg --out icon.iconset/icon_256x256@2x.png
sips -z 512 512   icon.svg --out icon.iconset/icon_512x512.png
sips -z 1024 1024 icon.svg --out icon.iconset/icon_512x512@2x.png

# 转换为icns
iconutil -c icns icon.iconset
```

## 推荐使用

- **开发阶段**: 直接使用SVG文件
- **发布阶段**: 使用fyne package工具自动处理图标
- **跨平台**: SVG格式最灵活，支持所有平台

## 许可证

这些图标是为本项目创建的，可以自由使用和修改。
