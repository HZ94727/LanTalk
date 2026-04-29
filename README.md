# LanTalk

简体中文 | [English](README.en-US.md)

LanTalk 是一个基于 `Wails + Go + Vite` 构建的局域网桌面通讯工具。它面向同一局域网内的电脑，提供自动发现、一对一聊天、图片发送、文件发送、本地历史记录和系统托盘等能力。

## 功能特性

- UDP 广播自动发现同一局域网内的 LanTalk 用户
- TCP 点对点发送文本、图片和文件消息
- 支持消息发送状态、失败重试、删除、转发和清空会话
- 支持图片预览、缩放、复制、另存为和打开所在位置
- 支持文件保存、打开、定位和可用性检查
- 使用 SQLite 保存联系人、设置和聊天记录
- 图片和文件落盘到本地媒体目录，并基于哈希去重
- 本地敏感字段使用 AES-GCM 加密保存
- 支持简体中文和英文界面
- 支持多主题切换
- 支持系统托盘、隐藏到托盘、未读数提示和单实例唤起
- 支持本机 Echo Bot 和手动添加联系人，方便调试

## 技术栈

- 桌面框架：Wails v2
- 后端语言：Go
- 前端构建：Vite
- 前端实现：原生 JavaScript、HTML、CSS
- 本地数据库：SQLite，使用 `modernc.org/sqlite`
- 局域网通信：UDP 广播发现，TCP 点对点传输
- 系统托盘：`github.com/getlantern/systray`

## 项目结构

```text
.
├── app.go                    # Wails 绑定层
├── main.go                   # 应用入口和窗口配置
├── tray.go                   # 系统托盘
├── internal
│   ├── chat                  # 聊天、发现、存储、加密和媒体管理
│   ├── clipimg               # 图片复制到剪贴板
│   └── mediautil             # 媒体类型、Data URL 和文件名工具
├── frontend
│   ├── src                   # 前端界面和交互
│   └── wailsjs               # Wails 生成绑定
├── build                     # 图标、清单和安装器模板
├── PROJECT_OVERVIEW.md       # 更详细的项目进度文档
└── wails.json                # Wails 配置
```

## 开发运行

安装前端依赖：

```bash
cd frontend
npm install
```

运行开发模式：

```bash
wails dev
```

如果本机还没有安装 Wails CLI：

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
```

## 构建

```bash
wails build
```

Windows 构建产物通常位于：

```text
build/bin/LanTalk.exe
```

## 本地数据

LanTalk 默认把数据保存在用户配置目录下：

```text
%AppData%\LanTalk
```

主要数据包括：

- `lantalk.db`：SQLite 数据库
- `master.key`：本地加密密钥
- `media/`：图片和文件等媒体内容

应用设置面板会显示当前保存目录，并支持迁移到新的目录。

## 测试

```bash
go test ./...
```

## 当前限制

- UDP 发现端口固定为 `48555`，同一台机器运行多个真实实例并不适合作为主要测试方式。
- 局域网发现依赖广播能力，防火墙、虚拟网卡或网络隔离可能影响发现结果。
- 文件传输目前适合小型附件，尚未实现大文件分片、断点续传和身份配对机制。

更多细节见 [PROJECT_OVERVIEW.md](PROJECT_OVERVIEW.md)。
