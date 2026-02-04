#!/bin/bash

# 今天吃什么 - 启动脚本
# 启动后端服务

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

# PID 文件
BACKEND_PID_FILE="$PROJECT_DIR/.backend.pid"

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

# 检查端口是否被占用
check_port() {
    local port=$1
    if lsof -i :$port > /dev/null 2>&1; then
        return 0  # 端口被占用
    else
        return 1  # 端口空闲
    fi
}

# 等待后端就绪
wait_for_backend() {
    local max_attempts=30
    local attempt=1
    
    log_info "Waiting for backend to be ready..."
    
    while [ $attempt -le $max_attempts ]; do
        if curl -s http://localhost:8080/health > /dev/null 2>&1; then
            log_success "Backend is ready!"
            return 0
        fi
        echo -n "."
        sleep 1
        attempt=$((attempt + 1))
    done
    
    echo ""
    log_error "Backend failed to start within ${max_attempts} seconds"
    return 1
}

# 停止后端
stop_backend() {
    if [ -f "$BACKEND_PID_FILE" ]; then
        local pid=$(cat "$BACKEND_PID_FILE")
        if kill -0 $pid 2>/dev/null; then
            log_info "Stopping backend (PID: $pid)..."
            kill $pid 2>/dev/null || true
            rm -f "$BACKEND_PID_FILE"
            log_success "Backend stopped"
        fi
    fi
    
    # 确保端口释放
    if check_port 8080; then
        log_warn "Port 8080 is still in use, trying to kill..."
        lsof -ti :8080 | xargs kill -9 2>/dev/null || true
    fi
}

# 启动后端
start_backend() {
    log_info "Starting backend..."
    
    # 检查端口
    if check_port 8080; then
        log_warn "Port 8080 is already in use"
        read -p "Kill existing process? (y/n) " -n 1 -r
        echo
        if [[ $REPLY =~ ^[Yy]$ ]]; then
            stop_backend
        else
            log_error "Cannot start backend, port 8080 is in use"
            exit 1
        fi
    fi
    
    cd "$BACKEND_DIR"
    
    # 编译后端
    log_info "Building backend..."
    go build -o server ./cmd/server/main.go
    
    # 启动后端（后台运行）
    ./server > ../backend.log 2>&1 &
    local pid=$!
    echo $pid > "$BACKEND_PID_FILE"
    
    log_info "Backend started with PID: $pid"
    
    # 等待后端就绪
    if ! wait_for_backend; then
        log_error "Failed to start backend"
        cat ../backend.log
        exit 1
    fi
    
    cd "$PROJECT_DIR"
}

# 停止所有服务
stop_all() {
    log_info "Stopping all services..."
    stop_backend
    log_success "All services stopped"
}

# 显示帮助
show_help() {
    echo "Usage: $0 [command]"
    echo ""
    echo "Commands:"
    echo "  start     Start backend (default)"
    echo "  backend   Start backend (same as start)"
    echo "  stop      Stop backend service"
    echo "  restart   Restart backend service"
    echo "  help      Show this help message"
}

# 主函数
main() {
    local command=${1:-start}
    
    case $command in
        start|backend)
            echo -e "${GREEN}"
            echo "╔═══════════════════════════════════════╗"
            echo "║      🍽️  今天吃什么 - What To Eat      ║"
            echo "╚═══════════════════════════════════════╝"
            echo -e "${NC}"
            start_backend
            log_success "Backend is running at http://localhost:8080"
            log_info "Press Ctrl+C to stop"
            tail -f "$PROJECT_DIR/backend.log"
            ;;
        stop)
            stop_all
            ;;
        restart)
            stop_all
            sleep 2
            start_backend
            log_success "Backend is running at http://localhost:8080"
            log_info "Press Ctrl+C to stop"
            tail -f "$PROJECT_DIR/backend.log"
            ;;
        help|--help|-h)
            show_help
            ;;
        *)
            log_error "Unknown command: $command"
            show_help
            exit 1
            ;;
    esac
}

# 捕获 Ctrl+C
trap 'echo ""; log_info "Shutting down..."; stop_backend; exit 0' INT TERM

# 执行主函数
main "$@"
