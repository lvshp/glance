# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目简介

ReadCLI 是一个基于 Go 的终端电子书阅读器，支持 TXT 和 EPUB 格式。基于 [glance](https://github.com/TimothyYe/glance) 项目开发，使用 tview（基于 tcell v2）构建 TUI 界面。

## 常用命令

```bash
# 构建（在 cmd/ 目录下执行）
cd cmd && make build VERSION=v0.2

# 直接构建
cd cmd && go build -o readcli

# 运行测试
go test ./...

# 运行单个包的测试
go test ./reader/...
go test ./core/...
go test ./lib/...

# 运行单个测试函数
go test ./reader/ -run TestEpubReader
go test ./core/ -run TestImport

# 清理构建产物
cd cmd && make clean

# 构建发布包
cd cmd && make release VERSION=v0.2
```

## 架构概览

### 三层架构

```
cmd/main.go          入口：解析参数，启动 core.Run()
  └─ core/           UI 层：状态管理、事件循环、界面渲染（基于 tview/tcell v2）
       └─ reader/    阅读层：TXT/EPUB 格式解析，统一 Reader 接口
  └─ lib/            基础设施层：存储、进度、更新、工具函数
```

### 核心包说明

**core/** — 应用核心，所有 UI 逻辑集中于此
- `ui.go` — `appState` 结构体持有全部运行时状态，`Run()` 是主事件循环入口。支持 9 种界面模式（`mode` 类型），通过状态机切换
- `event.go` — 键盘事件分发，处理 Vim 风格和方向键两套键位映射
- `consts.go` — 常量和帮助文本

**reader/** — 电子书解析，实现统一的 `Reader` 接口
- `Reader` 接口定义在 `reader.go`，包含 Load/Reflow/Search/翻页/章节跳转等方法
- `txt_reader.go` / `epub_reader.go` — 两种格式的具体实现
- `content_reader.go` — 内容读取器，处理分页逻辑
- `cache.go` / `progress_anchor.go` — 内容缓存和进度锚点

**lib/** — 数据持久化和工具
- `storage.go` — JSON 文件存储，管理配置(`config.json`)、书架(`bookshelf.json`)、书签(`bookmarks.json`)、进度(`progress.json`)
- 数据目录：`~/.readcli/`
- `update.go` — GitHub Releases 自动更新
- `progress.go` — 阅读进度计算

### 关键依赖

- `github.com/rivo/tview` — 终端 UI 框架
- `github.com/gdamore/tcell/v2` — 终端抽象层（tview 的依赖）
- `github.com/mattn/go-runewidth` — CJK 字符宽度计算

### 版本注入

版本号通过 `-ldflags "-X main.Version=xxx"` 在构建时注入 `cmd/main.go` 的 `Version` 变量。

## 开发注意事项

- 构建命令需在 `cmd/` 目录下执行（Makefile 位于该目录）
- 需要 Go 1.24+（go.mod 声明为 1.24）
- 新增阅读格式需实现 `reader.Reader` 接口
- 界面模式切换通过 `appState.mode` 状态机管理，新增模式需在 `ui.go` 和 `event.go` 中同步处理
