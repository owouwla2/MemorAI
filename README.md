# 🧠 MemorAI

**AI-native Windows 系统资源分析工具** | AI-native System Resource Analyzer for Windows

让 AI 帮你解读电脑为什么这么慢。
*Let AI explain why your PC is sluggish.*

[![License: MIT](https://img.shields.io/badge/License-MIT-green.svg)](LICENSE)
[![Built with Wails](https://img.shields.io/badge/Built%20with-Wails-red)](https://wails.io)
[![Go](https://img.shields.io/badge/Go-1.23+-00ADD8)](https://go.dev)
[![Platform](https://img.shields.io/badge/Platform-Windows-blue)]()
[![GitHub release](https://img.shields.io/github/v/release/owouwla2/MemorAI?include_prereleases&label=release)](https://github.com/owouwla2/MemorAI/releases/latest)
[![Downloads](https://img.shields.io/github/downloads/owouwla2/MemorAI/total.svg)](https://github.com/owouwla2/MemorAI/releases)
[![Build](https://github.com/owouwla2/MemorAI/actions/workflows/build.yml/badge.svg)](https://github.com/owouwla2/MemorAI/actions/workflows/build.yml)

---

## 🌐 Languages

[简体中文](#简体中文) | [English](#english)

---

## 简体中文

### 🤔 这是什么？

**MemorAI** 把任务管理器、资源监视器、Autoruns 和一个**会读懂你电脑的 AI 助手**捏成了一个工具。

不止是看数字——它告诉你：
- ✅ Chrome 占了 950 MB **但你只开了一个标签页**，是后台扩展在偷跑
- ✅ vmmem 吃掉 2.1 GB —— **你装过 Docker 但今天没用，这部分可以省**
- ✅ 卡巴斯基占 1.8 GB **是正常的**，但你可以调防护级别
- ✅ 这 3 个开机自启**完全没必要**，建议立即禁用

### ✨ 核心功能

| 功能 | 说明 |
|------|------|
| 🧮 **智能软件分组** | 把 `chrome.exe × 12` 自动合并成 "Chrome 浏览器 950MB"，普通人也看得懂 |
| 🤖 **AI 智能解读** | 接入任意 OpenAI 兼容 API（DeepSeek / Kimi / 通义 / Ollama），自动结合软件分组 + 自启项给出建议 |
| 💬 **AI 对话空间** | 不只是一次性分析，还能追问、讨论："vmmem 是什么？""我能禁用 SysMain 吗？" |
| 💉 **自定义提示词** | 告诉 AI "这些软件我必须保留"、"用更口语化的方式回答"等 |
| 📊 **30 秒采样图表** | 折线图（CPU + 内存 + Top 软件变化）+ 饼图（软件占比），可导出 HTML 报告 |
| 🚀 **自启项管理** | 一键启用/禁用开机启动项（注册表 Run + 启动文件夹），与任务管理器同款机制（StartupApproved），**不删除原始项可逆** |
| 🎨 **现代化 UI** | Tokyo Night 配色 / 暗亮主题切换 / Lucide 图标 / 等宽数字 |
| 📦 **单 exe** | 双击即用，约 12MB，无需安装运行时 |

### 🚀 快速开始

#### 直接使用（推荐）

1. 下载 [Releases](https://github.com/owouwla2/MemorAI/releases) 中的最新 `MemorAI.exe`
2. 双击运行（首次会自动用 Windows 自带的 WebView2 打开）
3. 进入「设置」配置你的 AI API（推荐 [DeepSeek](https://platform.deepseek.com)，便宜且国内速度快）
4. 点击「刷新快照」开始分析

> ⚠️ 想管理"系统级"自启项？右键 exe → "以管理员身份运行"

#### 从源码构建

前置：[Go 1.23+](https://go.dev/dl/) · [Node.js 18+](https://nodejs.org/) · [Wails CLI](https://wails.io/docs/gettingstarted/installation)

```bash
# 安装 Wails CLI
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 克隆项目
git clone https://github.com/owouwla2/MemorAI
cd MemorAI

# 开发模式（热重载）
wails dev

# 打包生产版
wails build
# 输出: build/bin/MemorAI.exe
```

### ⚙️ AI 配置

支持任何兼容 OpenAI Chat Completions 协议的接口。

| 服务商 | Base URL | 推荐模型 |
|--------|----------|----------|
| **DeepSeek** ⭐ | `https://api.deepseek.com/v1` | `deepseek-chat` |
| OpenAI | `https://api.openai.com/v1` | `gpt-4o-mini` |
| Kimi (Moonshot) | `https://api.moonshot.cn/v1` | `moonshot-v1-8k` |
| 通义千问 | `https://dashscope.aliyuncs.com/compatible-mode/v1` | `qwen-turbo` |
| 本地 Ollama | `http://localhost:11434/v1` | `llama3` 等 |

> 🔒 你的 API Key 只保存在本地 `config.json`，**不会上传到任何服务器**。

### 💡 自定义提示词示例

进入「设置 → 提示词注入」可以追加你的指令，让 AI 按你的偏好回答：

```text
- 我的电脑必须保留：微信、Chrome、卡巴斯基，请不要建议禁用
- 我习惯开 Docker，vmmem 占内存属于正常
- 用更直接犀利的口吻，不要客套
- 给出操作步骤时，使用 Win11 的界面术语
```

### 🛠 技术栈

- **后端**：Go 1.23+ · [Wails v2](https://wails.io) · [gopsutil](https://github.com/shirou/gopsutil)
- **前端**：Vanilla JS · [Chart.js](https://www.chartjs.org/) · 自制 SVG 图标系统
- **设计**：Tokyo Night 配色 · Inter + JetBrains Mono 字体

### 📂 项目结构

```
MemorAI/
├── main.go                # Wails 入口
├── app.go                 # 暴露给前端的方法
├── memory/                # 内存采集 + 软件分组
├── sysstat/               # CPU 等系统状态
├── monitor/               # 30 秒采样器
├── startup/               # 自启项枚举与开关
├── ai/                    # OpenAI 协议客户端
├── config/                # 本地配置管理
└── frontend/              # 前端 (HTML/CSS/JS)
```

### ❓ 为什么做这个？

我电脑开机就 60% 内存占用，找不到原因。试过 Process Explorer、Autoruns，但那些工具只给数字、不解读。

我想要：**给 AI 一份完整的系统快照，让它告诉我哪些可以关、哪些不能动**。

市面上没有这种工具，所以做了一个。

### 📋 RoadMap

- [x] 内存采集 + 软件分组 + 进程详情
- [x] AI 一键分析 + 自定义提示词
- [x] 30 秒采样 + 折线图/饼图 + 报告导出
- [x] 自启项管理（启用/禁用）
- [x] CPU 监控
- [x] 现代化 UI（暗亮主题切换）
- [ ] 网络流量监控（哪些进程在偷传数据）
- [ ] 磁盘 IO 监控
- [ ] 历史数据追踪（最近 7 天趋势）
- [ ] Linux / macOS 支持

### 🤝 贡献

欢迎 PR / Issue。特别需要：
- 进程归类规则补全（`memory/classifier.go`）
- 多语言支持

### 📜 License

MIT © MemorAI Contributors

---

## English

### 🤔 What is this?

**MemorAI** combines a task manager, resource monitor, Autoruns, and **an AI assistant that understands your PC** into a single tool.

More than numbers — it tells you:
- ✅ Chrome is using 950 MB **but you only have one tab open** — background extensions
- ✅ vmmem is eating 2.1 GB — **you installed Docker but aren't using it today, free this up**
- ✅ Kaspersky at 1.8 GB **is normal**, but you can lower the protection level
- ✅ These 3 startup items are **completely unnecessary** — disable them now

### ✨ Features

| Feature | Description |
|---------|-------------|
| 🧮 **Smart software grouping** | Auto-merges `chrome.exe × 12` into "Chrome Browser 950MB" |
| 🤖 **AI insights** | Plug in any OpenAI-compatible API (DeepSeek / Kimi / Ollama / etc.) |
| 💬 **AI chat space** | Not a one-shot — follow up: "What is vmmem?" "Can I disable SysMain?" |
| 💉 **Custom prompts** | Tell the AI your preferences and constraints |
| 📊 **30-sec sampling charts** | Line + Pie charts of CPU/Memory/Top software, exportable as HTML |
| 🚀 **Startup management** | One-click enable/disable startup entries (Registry + Folders), reversible (uses StartupApproved like Task Manager) |
| 🎨 **Modern UI** | Tokyo Night theme, dark/light mode, Lucide icons |
| 📦 **Single exe** | ~12 MB, no installer needed |

### 🚀 Quick Start

#### Use the binary

1. Download `MemorAI.exe` from [Releases](https://github.com/owouwla2/MemorAI/releases)
2. Double-click to run (uses Windows' built-in WebView2)
3. Open Settings, configure your AI API (recommended: DeepSeek)
4. Click "Refresh" to start analyzing

> ⚠️ To manage system-wide startup items, run as Administrator.

#### Build from source

Requirements: Go 1.23+ · Node.js 18+ · Wails CLI

```bash
go install github.com/wailsapp/wails/v2/cmd/wails@latest
git clone https://github.com/owouwla2/MemorAI
cd MemorAI
wails build
# Output: build/bin/MemorAI.exe
```

### ⚙️ AI Configuration

Works with any OpenAI Chat Completions-compatible endpoint:

| Provider | Base URL | Recommended Model |
|----------|----------|-------------------|
| **DeepSeek** ⭐ | `https://api.deepseek.com/v1` | `deepseek-chat` |
| OpenAI | `https://api.openai.com/v1` | `gpt-4o-mini` |
| Kimi | `https://api.moonshot.cn/v1` | `moonshot-v1-8k` |
| Local Ollama | `http://localhost:11434/v1` | `llama3` etc. |

> 🔒 Your API key is stored only in local `config.json`. **Never uploaded anywhere.**

### 🛠 Tech Stack

- **Backend**: Go 1.23+, Wails v2, gopsutil
- **Frontend**: Vanilla JS, Chart.js, custom SVG icon set
- **Design**: Tokyo Night palette, Inter + JetBrains Mono

### 📋 Roadmap

- [x] Memory analysis + Software grouping
- [x] AI one-click analysis + Custom prompts
- [x] 30-sec sampling + Charts + Report export
- [x] Startup item management
- [x] CPU monitoring
- [x] Modern UI with theme toggle
- [ ] Network traffic monitoring
- [ ] Disk I/O monitoring
- [ ] Historical trend tracking
- [ ] Linux / macOS support

### 📜 License

MIT
