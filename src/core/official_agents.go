// SPDX-License-Identifier: MIT
// Copyright (c) 2026 EasyAI contributors
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the Software), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
// The above copyright notice and this permission notice shall be included in
// all copies or substantial portions of the Software.
// THE SOFTWARE IS PROVIDED AS IS, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN
// THE SOFTWARE.

// Package core 提供可复用的 Agent 官方命令与配置资料目录。
//
// 本文件独立编写，仅整理公开的官方文档，不是 EasyAI 生产安装执行器的源码。
// 公开范围是此目录文件，不代表 EasyAI 整个客户端、服务端或核心实现全部开源。
// MIT 授权仅适用于本文件；Agent 本体、官方安装脚本和商标遵循各自授权。
//
// 使用：把本文件复制进自己的 Go 模块，调用 OfficialAgents() 查看数据。
// 本包没有执行命令、写配置、读取密钥、请求网络或修改环境变量的功能。
//
// 命令中的 <SESSION_ID> 等是占位符，执行前应替换。powershell 与 bash
// 命令不能混用；ui 表示图形界面操作而非 shell 命令。安装链接会随上游变化。
// 配置示例为新编写的最小片段，不应覆盖用户整个配置文件。
package core

const VerifiedOn = "2026-11-06"

// Source 标识每项操作依据的官方文档。
type Source struct {
	ID    string
	Title string
	URL   string
}

// Command 是资料，不会自动执行。Status 为 documented 或 ui-only。
type Command struct {
	Action    string
	Platforms []string
	Shell     string
	Text      string
	Status    string
	Notes     string
	SourceIDs []string
}

// ConfigFile.Path 保留上游的路径语法；本包不会展开 ~ 或环境变量。
type ConfigFile struct {
	Platforms []string
	Path      string
	Purpose   string
	SourceIDs []string
}

// ConfigExample 仅展示官方服务的最小配置或官方设置向导。
type ConfigExample struct {
	Title     string
	Format    string
	Content   string
	Notes     string
	SourceIDs []string
}

type AgentReference struct {
	ID             string
	Name           string
	Surface        string
	Platforms      []string
	Prerequisites  []string
	Commands       []Command
	ConfigFiles    []ConfigFile
	ConfigExamples []ConfigExample
	Notes          []string
	Sources        []Source
}

// OfficialAgents 每次返回新的目录切片。目录内容需按 VerifiedOn 和来源核对，
// 不代表在所有机器上都已执行安装；模型可用性取决于官方账号权限。
func OfficialAgents() []AgentReference {
	return []AgentReference{
		{
			ID:        "codex-cli",
			Name:      "Codex CLI",
			Surface:   "cli",
			Platforms: []string{"windows", "macos", "linux"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"macos", "linux", "wsl"},
					Shell:     "bash",
					Text:      "curl -fsSL https://chatgpt.com/codex/install.sh | sh",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "install",
					Platforms: []string{"windows"},
					Shell:     "powershell",
					Text:      "irm https://chatgpt.com/codex/install.ps1 | iex",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "install-alternative",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "npm install -g @openai/codex",
					Status:    "documented",
					Notes:     "需 Node.js/npm；与独立安装器择一使用。",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "install-alternative",
					Platforms: []string{"macos"},
					Shell:     "bash",
					Text:      "brew install --cask codex",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "codex",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "login",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "codex login",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"auth"},
				},
				{
					Action:    "resume-picker",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "codex resume",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "codex resume <SESSION_ID>",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "continue",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "codex resume --last",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos", "linux", "wsl"},
					Path:      "~/.codex/config.toml",
					Purpose:   "用户默认配置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%USERPROFILE%\\.codex\\config.toml",
					Purpose:   "Windows 用户默认配置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      ".codex/config.toml",
					Purpose:   "项目配置；仅在项目信任规则允许时加载。",
					SourceIDs: []string{"config"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "OpenAI GPT 模型示例",
					Format:    "toml",
					Content:   "model = \"gpt-6.1-sol\"",
					Notes:     "核对时官方文档中的示例；实际可用模型以官方账号权限和 /model 为准。",
					SourceIDs: []string{"config"},
				},
				{
					Title:     "OpenAI API 密钥认证（bash）",
					Format:    "bash",
					Content:   "export OPENAI_API_KEY=\"YOUR_OPENAI_API_KEY\"\nprintenv OPENAI_API_KEY | codex login --with-api-key",
					Notes:     "",
					SourceIDs: []string{"auth"},
				},
				{
					Title:     "OpenAI API 密钥认证（PowerShell）",
					Format:    "powershell",
					Content:   "$env:OPENAI_API_KEY = \"YOUR_OPENAI_API_KEY\"\n$env:OPENAI_API_KEY | codex login --with-api-key",
					Notes:     "",
					SourceIDs: []string{"auth"},
				},
			},
			Notes: []string{"也可用 codex login 登录 ChatGPT。", "本目录采用官方默认 ~/.codex，不设置全局 CODEX_HOME。", "恢复时使用 Codex 自己的会话 ID。"},
			Sources: []Source{
				{
					ID:    "install",
					Title: "OpenAI Codex 官方仓库",
					URL:   "https://github.com/openai/codex",
				},
				{
					ID:    "cli",
					Title: "Codex CLI 命令参考",
					URL:   "https://developers.openai.com/codex/cli/reference",
				},
				{
					ID:    "config",
					Title: "Codex 配置基础",
					URL:   "https://developers.openai.com/codex/config-basic",
				},
				{
					ID:    "auth",
					Title: "Codex 官方认证",
					URL:   "https://developers.openai.com/codex/auth",
				},
			},
		},
		{
			ID:        "codex-desktop",
			Name:      "Codex 桌面体验",
			Surface:   "desktop",
			Platforms: []string{"windows", "macos"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "打开官方下载页，选择对应系统并安装。",
					Status:    "ui-only",
					Notes:     "按官方当前提供的桌面应用版本操作。",
					SourceIDs: []string{"app", "download"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos"},
					Shell:     "terminal",
					Text:      "codex app",
					Status:    "documented",
					Notes:     "需要支持该命令的 Codex CLI，按官方提示打开或获取桌面应用。",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "在桌面应用的 Codex 会话列表中打开已有会话。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"app"},
				},
			},
			ConfigFiles: nil,
			ConfigExamples: []ConfigExample{
				{
					Title:     "官方账号配置",
					Format:    "ui",
					Content:   "登录 ChatGPT 账号，在 Codex 界面选择可用模型。",
					Notes:     "",
					SourceIDs: []string{"app"},
				},
			},
			Notes: []string{"官方页面当前将 Codex 桌面体验整合到 ChatGPT；以官方下载页为准。", "CLI 配置见 codex-cli 条目；不将桌面缓存或数据库当作稳定配置接口。"},
			Sources: []Source{
				{
					ID:    "app",
					Title: "OpenAI 桌面应用文档",
					URL:   "https://developers.openai.com/codex/app/",
				},
				{
					ID:    "install",
					Title: "Codex 官方 README",
					URL:   "https://github.com/openai/codex",
				},
				{
					ID:    "download",
					Title: "官方桌面下载",
					URL:   "https://chatgpt.com/download/",
				},
			},
		},
		{
			ID:        "chatgpt-desktop",
			Name:      "ChatGPT 桌面版",
			Surface:   "desktop",
			Platforms: []string{"windows", "macos"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "打开 https://chatgpt.com/download/，下载对应平台版本并按向导安装。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"download"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "从 Windows 开始菜单或 macOS 应用程序中打开 ChatGPT。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"download"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "从应用内会话列表打开已有对话。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"app"},
				},
			},
			ConfigFiles: nil,
			ConfigExamples: []ConfigExample{
				{
					Title:     "官方配置方式",
					Format:    "ui",
					Content:   "登录 ChatGPT 账号，在应用设置和对话界面选择可用功能与模型。",
					Notes:     "",
					SourceIDs: []string{"app", "download"},
				},
			},
			Notes: []string{"系统和芯片要求以官方下载页为准。", "普通 ChatGPT 对话没有在本文所引用资料中公布独立终端 resume 命令。", "官方未在本文引用资料中提供稳定的手动模型配置文件接口，因此不填写猜测的缓存路径。"},
			Sources: []Source{
				{
					ID:    "download",
					Title: "ChatGPT 官方下载",
					URL:   "https://chatgpt.com/download/",
				},
				{
					ID:    "app",
					Title: "OpenAI 桌面应用说明",
					URL:   "https://developers.openai.com/codex/app/",
				},
			},
		},
		{
			ID:        "claude-code",
			Name:      "Claude Code",
			Surface:   "cli",
			Platforms: []string{"windows", "macos", "linux"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"macos", "linux", "wsl"},
					Shell:     "bash",
					Text:      "curl -fsSL https://claude.ai/install.sh | bash",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"setup"},
				},
				{
					Action:    "install",
					Platforms: []string{"windows"},
					Shell:     "powershell",
					Text:      "irm https://claude.ai/install.ps1 | iex",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"setup"},
				},
				{
					Action:    "install-alternative",
					Platforms: []string{"windows"},
					Shell:     "terminal",
					Text:      "winget install Anthropic.ClaudeCode",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"setup"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "claude",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"setup", "cli"},
				},
				{
					Action:    "login",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "claude auth login",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "login-api-account",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "claude auth login --console",
					Status:    "documented",
					Notes:     "登录官方 Anthropic Console，使用 API 计费。",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "resume-picker",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "claude --resume",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "claude --resume <SESSION_ID>",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "continue",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "claude --continue",
					Status:    "documented",
					Notes:     "继续当前目录的最近会话。",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "select-model",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "claude-interactive",
					Text:      "/model",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"models"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos", "linux", "wsl"},
					Path:      "~/.claude/settings.json",
					Purpose:   "用户设置。",
					SourceIDs: []string{"settings"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%USERPROFILE%\\.claude\\settings.json",
					Purpose:   "Windows 用户设置。",
					SourceIDs: []string{"settings"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      ".claude/settings.json",
					Purpose:   "项目共享设置。",
					SourceIDs: []string{"settings"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      ".claude/settings.local.json",
					Purpose:   "项目本机设置。",
					SourceIDs: []string{"settings"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "Anthropic 官方模型别名",
					Format:    "json",
					Content:   "{\n  \"model\": \"sonnet\"\n}",
					Notes:     "官方别名随上游更新，也可通过 /model 选择账号可用模型。",
					SourceIDs: []string{"models"},
				},
			},
			Notes: []string{"本条目只使用 Anthropic 官方服务，不引入第三方模型映射。", "原生安装器的系统和 Windows 依赖要求请核对 setup 文档。"},
			Sources: []Source{
				{
					ID:    "setup",
					Title: "Claude Code 官方安装",
					URL:   "https://code.claude.com/docs/en/setup",
				},
				{
					ID:    "cli",
					Title: "Claude Code CLI 参考",
					URL:   "https://code.claude.com/docs/en/cli-reference",
				},
				{
					ID:    "settings",
					Title: "Claude Code 设置文件",
					URL:   "https://code.claude.com/docs/en/settings",
				},
				{
					ID:    "models",
					Title: "Claude Code 模型配置",
					URL:   "https://code.claude.com/docs/en/model-config",
				},
			},
		},
		{
			ID:        "claude-desktop",
			Name:      "Claude Desktop",
			Surface:   "desktop",
			Platforms: []string{"windows", "macos"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "打开 https://claude.com/download，下载并安装对应平台的 Claude。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"download"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "从开始菜单或应用程序目录打开 Claude。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"download"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "在 Claude 应用内打开历史对话。",
					Status:    "ui-only",
					Notes:     "不提供未经官方确认的桌面版 resume 命令。",
					SourceIDs: []string{"download"},
				},
				{
					Action:    "edit-mcp-config",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "Settings → Developer → Edit Config；修改 MCP 配置后重启应用。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"mcp"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos"},
					Path:      "~/Library/Application Support/Claude/claude_desktop_config.json",
					Purpose:   "MCP 服务器配置，不是模型 API 替换文件。",
					SourceIDs: []string{"mcp"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%APPDATA%\\Claude\\claude_desktop_config.json",
					Purpose:   "MCP 服务器配置。",
					SourceIDs: []string{"mcp"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "官方账号配置",
					Format:    "ui",
					Content:   "登录 Claude 账号，从应用对话界面选择可用模型。",
					Notes:     "",
					SourceIDs: []string{"download"},
				},
				{
					Title:     "MCP 配置结构",
					Format:    "json",
					Content:   "{\n  \"mcpServers\": {}\n}",
					Notes:     "空对象仅示意结构，不会连接服务器；不要覆盖已有配置。",
					SourceIDs: []string{"mcp"},
				},
			},
			Notes: []string{"Claude Desktop 与 Claude Code 的配置不同。", "不包含客户端补丁或第三方模型转换。"},
			Sources: []Source{
				{
					ID:    "download",
					Title: "Claude 官方下载",
					URL:   "https://claude.com/download",
				},
				{
					ID:    "mcp",
					Title: "MCP 官方 Claude Desktop 配置指南",
					URL:   "https://modelcontextprotocol.io/docs/develop/connect-local-servers",
				},
			},
		},
		{
			ID:        "gemini-cli",
			Name:      "Gemini CLI",
			Surface:   "cli",
			Platforms: []string{"windows", "macos", "linux"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "npm install -g @google/gemini-cli",
					Status:    "documented",
					Notes:     "需满足当前包 engines 要求的 Node.js 和 npm。",
					SourceIDs: []string{"readme"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "gemini",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"readme"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "gemini --resume <SESSION_ID>",
					Status:    "documented",
					Notes:     "使用官方会话列表提供的 UUID 或索引。",
					SourceIDs: []string{"sessions"},
				},
				{
					Action:    "continue",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "gemini --resume",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"sessions"},
				},
				{
					Action:    "list-sessions",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "gemini --list-sessions",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"sessions"},
				},
				{
					Action:    "resume-picker",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "gemini-interactive",
					Text:      "/resume",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"sessions"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos", "linux", "wsl"},
					Path:      "~/.gemini/settings.json",
					Purpose:   "用户设置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%USERPROFILE%\\.gemini\\settings.json",
					Purpose:   "Windows 用户设置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      ".gemini/settings.json",
					Purpose:   "项目设置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      ".gemini/.env",
					Purpose:   "环境变量文件；可能含密钥，不应公开。",
					SourceIDs: []string{"config"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "Google 官方 API 密钥（bash）",
					Format:    "bash",
					Content:   "export GEMINI_API_KEY=\"YOUR_GEMINI_API_KEY\"\ngemini",
					Notes:     "",
					SourceIDs: []string{"readme"},
				},
				{
					Title:     "Google 官方 API 密钥（PowerShell）",
					Format:    "powershell",
					Content:   "$env:GEMINI_API_KEY = \"YOUR_GEMINI_API_KEY\"\ngemini",
					Notes:     "",
					SourceIDs: []string{"readme", "config"},
				},
				{
					Title:     "默认模型示例",
					Format:    "json",
					Content:   "{\n  \"model\": {\n    \"name\": \"gemini-3-flash-preview\"\n  }\n}",
					Notes:     "示例来自核对时配置文档，实际模型以账号可用范围为准。",
					SourceIDs: []string{"config"},
				},
			},
			Notes: []string{"官方支持的账号类型和认证政策可能调整，请核对当前认证说明。"},
			Sources: []Source{
				{
					ID:    "readme",
					Title: "Google Gemini CLI 官方仓库",
					URL:   "https://github.com/google-gemini/gemini-cli",
				},
				{
					ID:    "config",
					Title: "Gemini CLI 配置参考",
					URL:   "https://geminicli.com/docs/reference/configuration/",
				},
				{
					ID:    "sessions",
					Title: "Gemini CLI 会话管理",
					URL:   "https://geminicli.com/docs/cli/session-management/",
				},
				{
					ID:    "keys",
					Title: "Google AI Studio API 密钥",
					URL:   "https://aistudio.google.com/apikey",
				},
			},
		},
		{
			ID:        "grok",
			Name:      "Grok Build",
			Surface:   "cli",
			Platforms: []string{"windows", "macos", "linux"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"macos", "linux", "wsl"},
					Shell:     "bash",
					Text:      "curl -fsSL https://x.ai/cli/install.sh | bash",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"overview"},
				},
				{
					Action:    "install",
					Platforms: []string{"windows"},
					Shell:     "powershell",
					Text:      "irm https://x.ai/cli/install.ps1 | iex",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"overview"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "login",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok login",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok --resume <SESSION_ID>",
					Status:    "documented",
					Notes:     "--resume 是选项，不是 resume 子命令。",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "continue",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok --continue",
					Status:    "documented",
					Notes:     "继续当前目录最近的会话。",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "resume-latest",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok --resume",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "list-sessions",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok sessions list",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "list-models",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok models",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "inspect-config",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "grok inspect",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos", "linux", "wsl"},
					Path:      "~/.grok/config.toml",
					Purpose:   "用户级配置。",
					SourceIDs: []string{"overview", "settings"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%USERPROFILE%\\.grok\\config.toml",
					Purpose:   "Windows 用户级配置。",
					SourceIDs: []string{"overview"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "官方 API 密钥（bash）",
					Format:    "bash",
					Content:   "export XAI_API_KEY=\"YOUR_XAI_API_KEY\"\ngrok",
					Notes:     "",
					SourceIDs: []string{"overview"},
				},
				{
					Title:     "官方 API 密钥（PowerShell）",
					Format:    "powershell",
					Content:   "$env:XAI_API_KEY = \"YOUR_XAI_API_KEY\"\ngrok",
					Notes:     "",
					SourceIDs: []string{"overview"},
				},
				{
					Title:     "官方登录与模型查看",
					Format:    "terminal",
					Content:   "grok login\ngrok models",
					Notes:     "默认官方服务不需要另写自定义代理。",
					SourceIDs: []string{"cli"},
				},
			},
			Notes: []string{"使用官方 Grok Build 安装器；模型选择以 grok models 输出为准。"},
			Sources: []Source{
				{
					ID:    "overview",
					Title: "Grok Build 官方入门",
					URL:   "https://docs.x.ai/build/overview",
				},
				{
					ID:    "cli",
					Title: "Grok Build CLI 参考",
					URL:   "https://docs.x.ai/build/cli/reference",
				},
				{
					ID:    "settings",
					Title: "Grok Build 设置",
					URL:   "https://docs.x.ai/build/settings",
				},
			},
		},
		{
			ID:        "hermes",
			Name:      "Hermes Agent / Desktop",
			Surface:   "cli-and-desktop",
			Platforms: []string{"windows", "macos", "linux"},
			Commands: []Command{
				{
					Action:    "install",
					Platforms: []string{"macos", "linux", "wsl"},
					Shell:     "bash",
					Text:      "curl -fsSL https://hermes-agent.nousresearch.com/install.sh | bash",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"quickstart"},
				},
				{
					Action:    "install",
					Platforms: []string{"windows"},
					Shell:     "powershell",
					Text:      "iex (irm https://hermes-agent.nousresearch.com/install.ps1)",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"quickstart"},
				},
				{
					Action:    "install-desktop",
					Platforms: []string{"windows", "macos"},
					Shell:     "ui",
					Text:      "从 Hermes 官方网站下载桌面安装包并安装。",
					Status:    "ui-only",
					Notes:     "",
					SourceIDs: []string{"install"},
				},
				{
					Action:    "configure",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "hermes setup --portal",
					Status:    "documented",
					Notes:     "使用 Nous Portal 官方登录。",
					SourceIDs: []string{"quickstart"},
				},
				{
					Action:    "select-model",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "hermes model",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"quickstart"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "hermes",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "start-desktop",
					Platforms: []string{"windows", "macos"},
					Shell:     "terminal",
					Text:      "hermes desktop",
					Status:    "documented",
					Notes:     "需要已安装桌面组件。",
					SourceIDs: []string{"desktop"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "hermes --resume <SESSION_ID>",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "continue",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "hermes --continue",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "list-sessions",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "hermes sessions list",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos", "linux", "wsl"},
					Path:      "~/.hermes/config.yaml",
					Purpose:   "模型、终端等非敏感设置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%USERPROFILE%\\.hermes\\config.yaml",
					Purpose:   "默认 Windows 配置位置。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      "~/.hermes/.env",
					Purpose:   "密钥与环境变量；~ 为用户主目录。",
					SourceIDs: []string{"config"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      "~/.hermes/auth.json",
					Purpose:   "OAuth 凭据，由官方工具管理，不应公开。",
					SourceIDs: []string{"config"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "Nous 官方服务配置",
					Format:    "terminal",
					Content:   "hermes setup --portal\nhermes model",
					Notes:     "使用官方向导写入当前版本配置，不手写过时的 provider 字段。",
					SourceIDs: []string{"quickstart"},
				},
			},
			Notes: []string{"官方脚本可能下载 Git、Python、Node.js 等依赖。", "使用 profile 或 HERMES_HOME 时，目录可能与默认值不同。"},
			Sources: []Source{
				{
					ID:    "quickstart",
					Title: "Hermes 官方快速开始",
					URL:   "https://hermes-agent.nousresearch.com/docs/getting-started/quickstart",
				},
				{
					ID:    "install",
					Title: "Hermes 官方安装",
					URL:   "https://hermes-agent.nousresearch.com/docs/getting-started/installation",
				},
				{
					ID:    "cli",
					Title: "Hermes CLI 会话管理",
					URL:   "https://hermes-agent.nousresearch.com/docs/user-guide/cli",
				},
				{
					ID:    "config",
					Title: "Hermes 配置参考",
					URL:   "https://hermes-agent.nousresearch.com/docs/user-guide/configuration",
				},
				{
					ID:    "desktop",
					Title: "Hermes Desktop",
					URL:   "https://hermes-agent.nousresearch.com/docs/user-guide/desktop",
				},
			},
		},
		{
			ID:        "deepseek-harness",
			Name:      "DeepSeek Harness",
			Surface:   "local-web",
			Platforms: []string{"windows", "macos", "linux"},
			Commands: []Command{
				{
					Action:    "install-and-start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "npx @deepseek-ai/dsh web",
					Status:    "documented",
					Notes:     "npm 按需获取并运行官方 Web UI；不是全局安装。",
					SourceIDs: []string{"readme"},
				},
				{
					Action:    "start",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "npx @deepseek-ai/dsh web",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"readme"},
				},
				{
					Action:    "help",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "terminal",
					Text:      "npx @deepseek-ai/dsh --help",
					Status:    "documented",
					Notes:     "",
					SourceIDs: []string{"cli"},
				},
				{
					Action:    "resume",
					Platforms: []string{"windows", "macos", "linux"},
					Shell:     "ui",
					Text:      "在 Web UI 的已有会话列表中选择要继续的会话。",
					Status:    "ui-only",
					Notes:     "默认 Web profile 未在引用资料中公布通用 CLI resume 命令；不能套用示例 TUI 插件参数。",
					SourceIDs: []string{"quickstart"},
				},
			},
			ConfigFiles: []ConfigFile{
				{
					Platforms: []string{"macos", "linux", "wsl"},
					Path:      "~/.dsh/profiles/web/cordis.patch.yml",
					Purpose:   "默认 Web profile 用户覆盖层。",
					SourceIDs: []string{"reference", "cli"},
				},
				{
					Platforms: []string{"windows"},
					Path:      "%USERPROFILE%\\.dsh\\profiles\\web\\cordis.patch.yml",
					Purpose:   "Windows 默认 Web profile 覆盖层。",
					SourceIDs: []string{"reference", "cli"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      "$DSH_HOME/cordis.patch.yml",
					Purpose:   "主目录覆盖层；DSH_HOME 默认 ~/.dsh，此处是目录标记而非通用 shell 语法。",
					SourceIDs: []string{"reference"},
				},
				{
					Platforms: []string{"windows", "macos", "linux"},
					Path:      "$DSH_HOME/.credentials.yaml",
					Purpose:   "官方 UI 管理的密钥，勿公开。",
					SourceIDs: []string{"reference", "providers"},
				},
			},
			ConfigExamples: []ConfigExample{
				{
					Title:     "官方 DeepSeek 配置",
					Format:    "ui",
					Content:   "启动 Web UI → Settings → Models → DeepSeek → 填写自己的官方 API key → 保存并选择可用模型。",
					Notes:     "使用官方界面维护凭据，不复制历史版本 settings.json。",
					SourceIDs: []string{"quickstart", "providers"},
				},
			},
			Notes: []string{"需满足当前包 engines 要求的 Node.js 和 npm。", "官方目前为快速迭代的开发者预览，入口与配置可能和历史版本不同。"},
			Sources: []Source{
				{
					ID:    "readme",
					Title: "DeepSeek AI 官方 Harness 仓库",
					URL:   "https://github.com/deepseek-ai/deepseek-harness",
				},
				{
					ID:    "cli",
					Title: "Harness CLI 入口",
					URL:   "https://github.com/deepseek-ai/deepseek-harness/blob/master/apps/cli/README.md",
				},
				{
					ID:    "reference",
					Title: "Harness CLI 配置参考",
					URL:   "https://github.com/deepseek-ai/deepseek-harness/blob/master/apps/cli/reference/README.md",
				},
				{
					ID:    "quickstart",
					Title: "Harness Web UI 指南",
					URL:   "https://deepseek-harness.github.io/deepseek-harness/guide/quickstart",
				},
				{
					ID:    "providers",
					Title: "Harness 模型与凭据",
					URL:   "https://deepseek-harness.github.io/deepseek-harness/guide/providers",
				},
			},
		},
	}
}
