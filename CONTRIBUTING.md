# Contributing to MemorAI

感谢你考虑为 MemorAI 贡献！

## 🐛 报告 Bug

发现 bug 请到 [Issues](https://github.com/owouwla2/MemorAI/issues) 提交，包含：
- 你的 Windows 版本（10 / 11 / 哪个 build）
- MemorAI 版本（在 footer 显示）
- 复现步骤
- 截图（如有）

## 💡 功能建议

直接在 Issues 提交，加上 `enhancement` 标签。先搜索一下避免重复。

## 🔧 开发环境

```bash
# 前置
go install github.com/wailsapp/wails/v2/cmd/wails@latest

# 克隆 + 开发
git clone https://github.com/owouwla2/MemorAI
cd MemorAI
wails dev   # 热重载
```

## 📝 提交 Pull Request

1. Fork 仓库
2. 创建分支：`git checkout -b feature/xxx` 或 `fix/xxx`
3. 修改 + 测试 (`wails build` 能成功)
4. 提交：commit message 用约定式格式
   - `feat: 添加 X 功能`
   - `fix: 修复 Y 问题`
   - `docs: 更新文档`
   - `refactor: 重构 Z`
5. Push 并发起 PR

## 🎯 优先需要的贡献

- **进程归类规则补全**：编辑 `memory/classifier.go`，添加更多软件识别规则（特别是国内软件）
- **国际化**：把界面文案抽出来支持多语言
- **Linux/macOS 支持**：替换 Windows 特定 API（`startup/` 包）
- **网络/磁盘 IO 监控**

## 📜 代码风格

- Go：跑 `go fmt ./...` 和 `go vet ./...`
- 前端：保持简洁，避免引入大型框架
- 中文注释 OK，但函数文档注释建议英文

## License

提交即表示你同意你的贡献以 MIT 协议开源。
