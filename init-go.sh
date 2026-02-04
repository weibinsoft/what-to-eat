#!/bin/bash

# Go 语言环境初始化脚本
# 用于初始化项目的 Go 开发环境

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 项目根目录
PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"

# 日志函数
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查命令是否存在
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# 检查 Go 是否已安装
check_go_installed() {
    if command_exists go; then
        GO_VERSION=$(go version | awk '{print $3}' | sed 's/go//')
        log_success "Go 已安装，版本: $GO_VERSION"
        return 0
    else
        log_warn "Go 未安装"
        return 1
    fi
}

# 检查 Go 版本是否符合要求
check_go_version() {
    local required_version="1.23"
    local current_version=$(go version | awk '{print $3}' | sed 's/go//' | cut -d. -f1,2)
    
    if [ "$(printf '%s\n' "$required_version" "$current_version" | sort -V | head -n1)" != "$required_version" ]; then
        log_error "Go 版本过低，需要 >= $required_version，当前: $current_version"
        return 1
    fi
    log_success "Go 版本符合要求: $current_version"
    return 0
}

# 使用 sdkman 安装 Go
install_go_with_sdkman() {
    if command_exists sdk; then
        log_info "使用 sdkman 安装 Go..."
        sdk install go 1.23.0 || {
            log_error "sdkman 安装 Go 失败"
            return 1
        }
        # 设置环境变量
        export SDKMAN_DIR="$HOME/.sdkman"
        [[ -s "$SDKMAN_DIR/bin/sdkman-init.sh" ]] && source "$SDKMAN_DIR/bin/sdkman-init.sh"
        log_success "Go 安装成功"
        return 0
    else
        log_warn "sdkman 未安装，跳过自动安装"
        return 1
    fi
}

# 配置 Go 环境变量
setup_go_env() {
    log_info "配置 Go 环境变量..."
    
    # 如果使用 sdkman
    if [ -s "$HOME/.sdkman/bin/sdkman-init.sh" ]; then
        export SDKMAN_DIR="$HOME/.sdkman"
        source "$SDKMAN_DIR/bin/sdkman-init.sh"
    fi
    
    # 设置 Go 代理（可选，针对国内用户）
    if [ -z "$GOPROXY" ]; then
        export GOPROXY=https://goproxy.cn,direct
        log_info "设置 GOPROXY: $GOPROXY"
    fi
    
    # 设置 GOSUMDB（可选）
    if [ -z "$GOSUMDB" ]; then
        export GOSUMDB=sum.golang.google.cn
        log_info "设置 GOSUMDB: $GOSUMDB"
    fi
    
    # 显示当前 Go 环境
    if command_exists go; then
        log_info "GOROOT: $(go env GOROOT)"
        log_info "GOPATH: $(go env GOPATH)"
        log_info "GOPROXY: $(go env GOPROXY)"
    fi
}

# 下载 Go 依赖
download_dependencies() {
    log_info "进入 backend 目录..."
    cd "$BACKEND_DIR" || {
        log_error "无法进入 backend 目录"
        return 1
    }
    
    log_info "下载 Go 模块依赖..."
    go mod download || {
        log_error "下载依赖失败"
        return 1
    }
    
    log_info "整理依赖..."
    go mod tidy || {
        log_error "整理依赖失败"
        return 1
    }
    
    log_success "依赖下载完成"
}

# 验证 Go 环境
verify_go_env() {
    log_info "验证 Go 环境..."
    
    cd "$BACKEND_DIR" || return 1
    
    # 检查 go.mod 文件
    if [ ! -f "go.mod" ]; then
        log_error "go.mod 文件不存在"
        return 1
    fi
    
    # 尝试编译（不生成二进制文件）
    log_info "验证代码编译..."
    go build -o /dev/null ./cmd/server || {
        log_error "代码编译失败"
        return 1
    }
    
    log_success "Go 环境验证通过"
}

# 主函数
main() {
    log_info "开始初始化 Go 语言环境..."
    echo ""
    
    # 检查 Go 是否已安装
    if ! check_go_installed; then
        log_info "尝试自动安装 Go..."
        if ! install_go_with_sdkman; then
            log_error "无法自动安装 Go，请手动安装："
            echo "  1. 访问 https://go.dev/dl/ 下载并安装"
            echo "  2. 或使用 Homebrew: brew install go"
            echo "  3. 或使用 sdkman: sdk install go"
            exit 1
        fi
    fi
    
    # 检查 Go 版本
    if ! check_go_version; then
        log_error "请升级 Go 版本"
        exit 1
    fi
    
    # 配置环境变量
    setup_go_env
    
    # 下载依赖
    download_dependencies
    
    # 验证环境
    verify_go_env
    
    echo ""
    log_success "Go 语言环境初始化完成！"
    log_info "可以运行以下命令启动后端："
    echo "  cd backend && go run cmd/server/main.go"
    echo "  或使用项目启动脚本："
    echo "  ./start.sh backend"
}

# 执行主函数
main
