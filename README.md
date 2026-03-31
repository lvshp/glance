# Glance 维护版

[![Release](https://img.shields.io/github/v/release/lvshp/glance?label=%E5%8F%91%E5%B8%83%E7%89%88)](https://github.com/lvshp/glance/releases)
[![CI](https://img.shields.io/github/actions/workflow/status/lvshp/glance/go.yml?branch=main&label=CI)](https://github.com/lvshp/glance/actions/workflows/go.yml)
[![License](https://img.shields.io/badge/license-Apache%202.0-blue.svg)](./LICENSE)

<img src="https://github.com/TimothyYe/glance/blob/master/demo/glance.png?raw=true" width="120" />

一个持续维护中的命令行小说阅读工具分支，主打 EPUB、进度恢复、终端自适应排版，以及更像 IDE 的摸鱼伪装界面。

适合这样的使用场景：

* 在终端里阅读 TXT / EPUB
* 在工作环境下低调摸鱼
* 想保留键盘流操作和轻量体验
* 想要一个还在持续维护的 Glance 分支

![](https://github.com/TimothyYe/glance/blob/master/demo/demo.gif?raw=true)

[English Version](#)

## Fork 说明

这个仓库基于原项目 [TimothyYe/glance](https://github.com/TimothyYe/glance) 修改而来。

我在保留原项目协议与致谢的基础上，继续维护并补充了以下能力：

* EPUB 支持
* EPUB 章节目录与目录跳转
* 自动记忆阅读进度并恢复
* 多行显示与按页翻页
* 按终端宽度自适应折行
* 更偏 IDE / 工作台风格的伪装界面

## 下载

发布版本与预编译包：

* [Releases 页面](https://github.com/lvshp/glance/releases)
* [v0.1.0 首个维护版发布](https://github.com/lvshp/glance/releases/tag/v0.1.0)

当前已提供：

* macOS arm64
* macOS amd64
* Linux amd64

## 功能亮点

* 使用Go开发，无需额外运行时和依赖库。
* 软件运行于命令行，对Vimer友好，支持Vim方式的Key Binding进行翻页和跳转。
* 支持Boss Key，方便紧急情况下对界面隐藏和伪装。
* 支持自动定时翻页模式

## 维护版增强

相比原项目，这个维护版本重点增强了下面几块：

* EPUB 支持与章节目录导航
* 记忆阅读进度与自动恢复
* 多行显示、按页滚动、数字跳章
* 按终端宽度与中文显示宽度自适应排版
* IDE 风格工作台界面，更适合伪装成开发环境
* 多平台 Release 附件与自动发布流程

## 安装步骤

如果你使用的是这个维护版本，可以直接下载 Release 附件，或者从源码构建：

```bash
git clone https://github.com/lvshp/glance.git
cd glance
go build -o glance ./cmd
./glance -n 8 /path/to/book.epub
```

如果你已经把可执行文件放进 `PATH`，也可以直接：

```bash
glance -n 8 /path/to/book.epub
```

原项目安装方式如下：

```bash
brew tap timothyye/tap
brew install timothyye/tap/glance
```

注: 也可以选择[直接下载](https://github.com/TimothyYe/glance/releases)可执行文件并运行。

## 支持平台
* Mac OS
* Linux
* Windows (计划中)

## 支持格式
* TXT (已支持)
* PDF (计划中)
* EPUB (已支持)

## 快捷键说明

* `?` 显示与隐藏帮助菜单
* `q` 或者 `ctrl+c` 退出程序
* `j`, `ctrl+n`, `<Space>` 或者 `<Enter>` 按当前显示行数向下翻一页
* `k` 或者 `ctrl+p` 按当前显示行数向上翻一页
* `+` / `-` 增加或减少同时显示的正文行数
* `[` 跳转到上一章（EPUB）
* `]` 跳转到下一章（EPUB）
* `m` 显示与隐藏目录（EPUB，目录中可用 `j/k` 选择、数字加 `Enter` 跳转，`*` 为当前章节）
* `p` 显示与隐藏当前阅读进度
* `b` Boss Key，隐藏当前内容并显示伪装Shell提示符
* `f` 显示与隐藏边框
* `c` 切换显示字体颜色

## 跳转命令

Glance支持与Vim相同的快捷跳转命令，方便在阅读时快速定位以及跳转到想要阅读的位置。例如:

* `G` 跳转到最后一行  
* `50G` 跳转到第50行
* `gg` 跳转到第一行
* `20j` 向下跳转20行
* `30k` 向上跳转30行

## 开发环境集成展示

Glance可以运行在任何支持Terminal的开发软件及环境中，包括并不仅限于JetBrains全家桶, Vim, Tmux, Emacs……

* GoLand
![](https://github.com/TimothyYe/glance/blob/master/demo/goland.png?raw=true)

* Spacemacs
![](https://github.com/TimothyYe/glance/blob/master/demo/spacemacs.png?raw=true)

* VSCode
![](https://github.com/TimothyYe/glance/blob/master/demo/vscode.png?raw=true)

* Tmux
![](https://github.com/TimothyYe/glance/blob/master/demo/tmux.png?raw=true)

## Issue 与 PR

欢迎基于这个维护版本提交 issue 与 PR。

## 协议

本项目基于原项目继续演进，并继续遵循 [Apache License 2.0](./LICENSE)。
