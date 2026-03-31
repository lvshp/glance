# Glance

<img src="https://github.com/TimothyYe/glance/blob/master/demo/glance.png?raw=true" width="120" />

一款基于命令行跨平台文本小说阅读工具，996与10107程序员摸鱼划水必备神器。

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

## 功能亮点

* 使用Go开发，无需额外运行时和依赖库。
* 软件运行于命令行，对Vimer友好，支持Vim方式的Key Binding进行翻页和跳转。
* 支持Boss Key，方便紧急情况下对界面隐藏和伪装。
* 支持自动定时翻页模式

## 安装步骤

如果你使用的是这个维护版本，推荐直接从源码构建：

```bash
git clone <your-repo-url>
cd glance
go build -o glance ./cmd
./glance -n 8 /path/to/book.epub
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
