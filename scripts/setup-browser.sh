#!/bin/bash
# Cursor2API 浏览器安装脚本
# 用于安装 Chromium 浏览器依赖

set -e

echo "=========================================="
echo "  Cursor2API 浏览器安装脚本"
echo "=========================================="
echo ""

# 检测操作系统
detect_os() {
    case "$(uname -s)" in
        Darwin*)    echo "macos";;
        Linux*)     echo "linux";;
        MINGW*|MSYS*|CYGWIN*)   echo "windows";;
        *)          echo "unknown";;
    esac
}

OS=$(detect_os)
echo "[信息] 检测到操作系统: $OS"

# 检查是否已安装浏览器
check_existing_browser() {
    echo ""
    echo "[检查] 正在检查已安装的浏览器..."
    
    case "$OS" in
        macos)
            browsers=(
                "/Applications/Google Chrome.app/Contents/MacOS/Google Chrome"
                "/Applications/Chromium.app/Contents/MacOS/Chromium"
                "/Applications/Microsoft Edge.app/Contents/MacOS/Microsoft Edge"
                "/Applications/Brave Browser.app/Contents/MacOS/Brave Browser"
            )
            ;;
        linux)
            browsers=(
                "/usr/bin/chromium"
                "/usr/bin/chromium-browser"
                "/usr/bin/google-chrome"
                "/usr/bin/google-chrome-stable"
                "/snap/bin/chromium"
            )
            ;;
        *)
            browsers=()
            ;;
    esac
    
    for browser in "${browsers[@]}"; do
        if [ -f "$browser" ]; then
            echo "[成功] 找到浏览器: $browser"
            echo ""
            echo "您可以使用此浏览器，设置环境变量:"
            echo "  export BROWSER_PATH=\"$browser\""
            echo ""
            echo "或在 config.yaml 中设置:"
            echo "  browser:"
            echo "    path: \"$browser\""
            return 0
        fi
    done
    
    echo "[信息] 未找到已安装的兼容浏览器"
    return 1
}

# macOS 安装
install_macos() {
    echo ""
    echo "[安装] macOS 浏览器安装选项:"
    echo ""
    echo "方式 1: 使用 Homebrew 安装 Chromium (推荐)"
    echo "  brew install --cask chromium"
    echo ""
    echo "方式 2: 使用 Homebrew 安装 Google Chrome"
    echo "  brew install --cask google-chrome"
    echo ""
    echo "方式 3: 让 go-rod 自动下载 (无需手动安装)"
    echo "  只需保持 config.yaml 中 browser.path 为空"
    echo "  程序启动时会自动下载 Chromium 到 ~/.cache/rod/browser/"
    echo ""
    
    # 检查是否有 Homebrew
    if command -v brew &> /dev/null; then
        echo "[提示] 检测到 Homebrew，是否自动安装 Chromium? (y/n)"
        read -r answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            echo "[安装] 正在安装 Chromium..."
            brew install --cask chromium
            echo "[成功] Chromium 安装完成!"
            echo ""
            echo "浏览器路径: /Applications/Chromium.app/Contents/MacOS/Chromium"
        fi
    else
        echo "[提示] 未检测到 Homebrew"
        echo "安装 Homebrew: /bin/bash -c \"\$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)\""
    fi
}

# Linux 安装
install_linux() {
    echo ""
    echo "[安装] Linux 浏览器安装选项:"
    echo ""
    
    # 检测包管理器
    if command -v apt-get &> /dev/null; then
        echo "方式 1: 使用 apt 安装 (Debian/Ubuntu)"
        echo "  sudo apt-get update && sudo apt-get install -y chromium-browser"
        echo ""
        echo "是否自动安装? (y/n)"
        read -r answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            sudo apt-get update && sudo apt-get install -y chromium-browser
            echo "[成功] 安装完成!"
        fi
    elif command -v yum &> /dev/null; then
        echo "方式 1: 使用 yum 安装 (CentOS/RHEL)"
        echo "  sudo yum install -y chromium"
        echo ""
        echo "是否自动安装? (y/n)"
        read -r answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            sudo yum install -y chromium
            echo "[成功] 安装完成!"
        fi
    elif command -v pacman &> /dev/null; then
        echo "方式 1: 使用 pacman 安装 (Arch)"
        echo "  sudo pacman -S chromium"
        echo ""
        echo "是否自动安装? (y/n)"
        read -r answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            sudo pacman -S chromium
            echo "[成功] 安装完成!"
        fi
    elif command -v apk &> /dev/null; then
        echo "方式 1: 使用 apk 安装 (Alpine)"
        echo "  apk add --no-cache chromium"
        echo ""
        echo "是否自动安装? (y/n)"
        read -r answer
        if [ "$answer" = "y" ] || [ "$answer" = "Y" ]; then
            apk add --no-cache chromium
            echo "[成功] 安装完成!"
        fi
    else
        echo "未检测到已知的包管理器"
    fi
    
    echo ""
    echo "方式 2: 让 go-rod 自动下载 (无需手动安装)"
    echo "  只需保持 config.yaml 中 browser.path 为空"
}

# Docker 安装提示
docker_tips() {
    echo ""
    echo "=========================================="
    echo "  Docker 用户提示"
    echo "=========================================="
    echo ""
    echo "如果在 Docker 中运行，请在 Dockerfile 中添加:"
    echo ""
    echo "# Debian/Ubuntu 基础镜像"
    echo "RUN apt-get update && apt-get install -y \\"
    echo "    chromium \\"
    echo "    fonts-liberation \\"
    echo "    libasound2 \\"
    echo "    libatk-bridge2.0-0 \\"
    echo "    libatk1.0-0 \\"
    echo "    libatspi2.0-0 \\"
    echo "    libcups2 \\"
    echo "    libdbus-1-3 \\"
    echo "    libdrm2 \\"
    echo "    libgbm1 \\"
    echo "    libgtk-3-0 \\"
    echo "    libnspr4 \\"
    echo "    libnss3 \\"
    echo "    libxcomposite1 \\"
    echo "    libxdamage1 \\"
    echo "    libxfixes3 \\"
    echo "    libxkbcommon0 \\"
    echo "    libxrandr2 \\"
    echo "    xdg-utils \\"
    echo "    && rm -rf /var/lib/apt/lists/*"
    echo ""
    echo "# Alpine 基础镜像"
    echo "RUN apk add --no-cache chromium"
    echo ""
    echo "# 设置环境变量"
    echo "ENV BROWSER_PATH=/usr/bin/chromium"
}

# 主流程
main() {
    # 先检查已安装的浏览器
    if check_existing_browser; then
        echo "[完成] 您已有可用的浏览器!"
        docker_tips
        exit 0
    fi
    
    # 根据操作系统安装
    case "$OS" in
        macos)
            install_macos
            ;;
        linux)
            install_linux
            ;;
        windows)
            echo ""
            echo "[信息] Windows 用户请手动安装 Google Chrome 或 Edge 浏览器"
            echo "下载地址: https://www.google.com/chrome/"
            ;;
        *)
            echo "[错误] 不支持的操作系统"
            exit 1
            ;;
    esac
    
    docker_tips
    
    echo ""
    echo "=========================================="
    echo "  安装完成后运行程序"
    echo "=========================================="
    echo ""
    echo "go run cmd/server/main.go"
    echo ""
    echo "或设置环境变量后运行:"
    echo "BROWSER_PATH=/path/to/browser go run cmd/server/main.go"
}

main
