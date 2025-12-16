# Cursor2API

将 Cursor API 转换为 OpenAI/Anthropic 兼容格式的代理服务。

## 原理

本项目利用 [Cursor 文档页面](https://cursor.com/cn/docs) 提供的免费 AI 聊天功能。该页面内置了一个 AI 助手，通过 `https://cursor.com/api/chat` 接口与后端通信。

**关键特点：**
- **无需登录** - 文档页面的 AI 聊天功能对所有访问者开放
- **无需 API Key** - 不需要 Cursor 账号或付费订阅
- **支持多模型** - 可使用 Claude、GPT、Gemini 等模型

本项目通过浏览器自动化技术访问该页面，将请求转发到 Cursor API，并将响应转换为标准的 OpenAI/Anthropic API 格式。

## 功能特性

- **Anthropic Messages API** - 完整支持 `/v1/messages` 接口
- **OpenAI Chat API** - 支持 `/v1/chat/completions` 接口
- **流式响应** - 支持 SSE 流式输出
- **浏览器自动化** - 自动处理人机验证

## 项目结构

```
cursor2api/
├── cmd/server/          # 程序入口
│   └── main.go
├── internal/            # 内部包
│   ├── browser/         # 浏览器自动化服务
│   ├── config/          # 配置管理
│   └── handler/         # HTTP 处理器
├── static/              # 静态文件
├── config.yaml          # 配置文件
└── README.md
```

## 快速开始

### Docker 部署 (推荐)

```bash
# 使用 docker-compose
docker-compose up -d

# 或者手动构建运行
docker build -t cursor2api .
docker run -d -p 3010:3010 --shm-size=2g cursor2api
```

### 本地运行

```bash
# 安装依赖
go mod tidy

# 编译
go build -o cursor2api ./cmd/server

# 运行
./cursor2api
```

服务默认运行在 `http://localhost:3010`

## 浏览器安装

程序需要 Chromium 内核浏览器。有以下几种方式：

### 方式 1: 自动下载 (推荐)

保持 `config.yaml` 中 `browser.path` 为空，程序会：
1. 首先自动检测系统已安装的 Chrome/Chromium/Edge
2. 如果未找到，则自动下载 Chromium 到 `~/.cache/rod/browser/`

### 方式 2: 使用安装脚本

```bash
# 运行安装脚本
./scripts/setup-browser.sh
```

### 方式 3: 手动安装

**macOS:**
```bash
brew install --cask chromium
# 或
brew install --cask google-chrome
```

**Linux (Debian/Ubuntu):**
```bash
sudo apt-get update && sudo apt-get install -y chromium-browser
```

**Linux (Alpine):**
```bash
apk add --no-cache chromium
```

### 方式 4: 使用环境变量

```bash
# 指定浏览器路径
export BROWSER_PATH="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
./cursor2api
```

## 配置

编辑 `config.yaml`：

```yaml
# 服务端口
port: 3010

# 浏览器设置
browser:
  headless: true
  # 留空则自动检测或下载，也可手动指定路径
  path: ""
```

支持的环境变量：
- `PORT` - 覆盖端口配置
- `BROWSER_PATH` - 覆盖浏览器路径

## API 接口

### Anthropic Messages API

```bash
curl http://localhost:3010/v1/messages \
  -H "Content-Type: application/json" \
  -H "x-api-key: any" \
  -d '{
    "model": "claude-sonnet-4-20250514",
    "max_tokens": 1024,
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": true
  }'
```

### OpenAI Chat API

```bash
curl http://localhost:3010/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "gpt-4",
    "messages": [{"role": "user", "content": "Hello"}],
    "stream": true
  }'
```

### 其他接口

- `GET /v1/models` - 获取模型列表
- `GET /health` - 健康检查
- `GET /browser/status` - 浏览器状态

## Claude Code 集成

```bash
# 设置 API 地址
export ANTHROPIC_BASE_URL=http://localhost:3010

# 运行 Claude Code
claude
```

## 支持的模型

| 请求模型 | 映射到 Cursor |
|---------|--------------|
| claude-* | anthropic/claude-sonnet-4.5 |
| gpt-* | openai/gpt-5-nano |
| gemini-* | google/gemini-2.5-flash |

## 依赖

- Go 1.21+
- Chromium 浏览器

## 免责声明

本项目仅供学习和研究目的使用。

- 本项目是一个非官方的第三方工具，与 Cursor 官方无任何关联
- 使用本项目可能违反 Cursor 的服务条款，请自行承担风险
- 本项目不提供任何形式的担保，包括但不限于适销性、特定用途适用性
- 作者不对使用本项目造成的任何直接或间接损失负责
- 请勿将本项目用于商业用途或任何违法活动

使用本项目即表示您已阅读并同意以上声明。

## 许可证

MIT
