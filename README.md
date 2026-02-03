# 今天吃什么 - What To Eat

一个具有"仪式感"的午餐决策 App，帮助选择困难症患者做出今天吃什么的决定。

## 功能特性

- 🎰 老虎机风格的随机决策动画
- 📋 餐厅列表管理（添加/删除）
- 📊 最近5天用餐历史记录
- 🎯 加权随机算法（避免连续吃同样的食物）
- 🔐 用户登录/注册

## 技术栈

### 前端

- Vue 3 + TypeScript
- Vite
- Pinia (状态管理)
- Vue Router
- Axios

### 后端

- Go 1.21+
- Gin Web Framework
- GORM
- JWT 认证
- MySQL 8.0

## 快速开始

### 方式一：一键启动（推荐）

确保 MySQL 服务已启动，然后执行：

```bash
# 赋予执行权限（首次运行）
chmod +x start.sh

# 启动项目（自动创建数据库和表）
./start.sh
```

启动脚本会自动：

1. 编译并启动后端
2. 初始化数据库和数据表
3. 启动前端开发服务器

### 方式二：手动启动

#### 1. 配置数据库

编辑 `backend/config/config.yaml`：

```yaml
database:
  host: "localhost"
  port: "3306"
  user: "root"
  password: "your-password"
  dbname: "what_to_eat"
```

#### 2. 启动后端

```bash
cd backend
go mod tidy
go run cmd/server/main.go
```

后端启动时会自动创建数据库和表，并插入默认餐厅数据。

#### 3. 启动前端

```bash
cd frontend
npm install
npm run dev
```

前端将在 <http://localhost:5173> 启动。

### 启动脚本命令

```bash
./start.sh start     # 启动后端和前端（默认）
./start.sh stop      # 停止所有服务
./start.sh restart   # 重启所有服务
./start.sh backend   # 只启动后端
./start.sh frontend  # 只启动前端
./start.sh help      # 显示帮助
```

### 4. 开始使用

1. 访问 <http://localhost:5173>
2. 注册一个新账号或登录
3. 在右侧添加喜欢的餐厅
4. 点击"开始决策"按钮，观看老虎机动画
5. 查看历史记录，防止连续吃同样的食物

## API 接口

### 认证

- `POST /api/auth/register` - 用户注册
- `POST /api/auth/login` - 用户登录

### 餐厅管理

- `GET /api/restaurants` - 获取餐厅列表
- `POST /api/restaurants` - 添加餐厅
- `DELETE /api/restaurants/:id` - 删除餐厅

### 决策

- `POST /api/decide` - 执行随机决策
- `GET /api/history` - 获取最近5天的历史记录

## 加权随机算法

为了避免用户连续多天吃同样的食物，系统实现了加权随机算法：

1. 获取用户最近3次的决策记录
2. 如果某餐厅在最近3次中出现过，其被选中的概率降低50%
3. 使用加权随机选择最终结果

## 项目结构

```
.
├── start.sh                 # 一键启动脚本
├── frontend/                # Vue 3 前端
│   ├── src/
│   │   ├── api/            # API 请求封装
│   │   ├── components/     # 组件
│   │   ├── stores/         # Pinia 状态管理
│   │   ├── views/          # 页面
│   │   └── router/         # 路由配置
│   └── package.json
├── backend/                 # Go 后端
│   ├── cmd/server/         # 入口文件
│   ├── config/             # Viper 配置
│   │   ├── config.go       # 配置加载逻辑
│   │   └── config.yaml     # 配置文件
│   ├── internal/
│   │   ├── handler/        # HTTP 处理器
│   │   ├── service/        # 业务逻辑
│   │   ├── repository/     # 数据访问（含数据库初始化）
│   │   └── model/          # 数据模型
│   ├── pkg/middleware/     # 中间件
│   └── go.mod
├── sql/
│   └── init.sql            # 数据库初始化脚本（备用）
└── README.md
```

## 配置说明

### 后端配置

配置文件位于 `backend/config/config.yaml`，使用 Viper 管理：

```yaml
server:
  port: "8080"           # 服务端口
  mode: "debug"          # 运行模式：debug/release/test

database:
  host: "localhost"      # 数据库主机
  port: "3306"           # 数据库端口
  user: "root"           # 数据库用户名
  password: "123456"     # 数据库密码
  dbname: "what_to_eat"  # 数据库名
  max_idle_conns: 10     # 最大空闲连接数
  max_open_conns: 100    # 最大打开连接数
  conn_max_lifetime: 3600 # 连接最大生命周期（秒）

jwt:
  secret: "your-secret-key"  # JWT 签名密钥
  expire_time: 24            # Token 过期时间（小时）
```

也支持环境变量覆盖，前缀为 `APP_`：

- `APP_DATABASE_HOST`
- `APP_DATABASE_PASSWORD`
- `APP_JWT_SECRET`

### 前端配置

配置文件位于 `frontend/.env`：

| 变量 | 默认值 | 说明 |
|------|--------|------|
| VITE_API_BASE_URL | <http://localhost:8080> | 后端 API 地址 |
