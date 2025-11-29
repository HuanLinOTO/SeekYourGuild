# SeekYourGuild

<div align="center">

纯 AI 编写，谨慎使用

![Logo](https://img.shields.io/badge/SeekYourGuild-群数据监控-blue?style=for-the-badge)

**专业的 QQ 群数据监控与分析平台**

[![Go](https://img.shields.io/badge/Go-1.21+-00ADD8?style=flat-square&logo=go)](https://go.dev/)
[![Vue](https://img.shields.io/badge/Vue-3.x-4FC08D?style=flat-square&logo=vue.js)](https://vuejs.org/)
[![License](https://img.shields.io/badge/License-MIT-yellow?style=flat-square)](LICENSE)
[![Build](https://img.shields.io/github/actions/workflow/status/HuanLinOTO/SeekYourGuild/build.yml?style=flat-square)](https://github.com/HuanLinOTO/SeekYourGuild/actions)

[功能特性](#-功能特性) • [快速开始](#-快速开始) • [部署指南](#-部署指南) • [API 文档](#-api-文档) • [开发指南](#-开发指南)

</div>

---

## ✨ 功能特性

### 📊 数据监控

- **成员统计**: 实时群成员数量、今日新增/离开
- **消息分析**: 24 小时消息分布、消息密度统计
- **活跃排行**: 用户发言排行榜、QQ 头像展示
- **机器人状态**: 在线状态监控、心跳检测

### 📈 数据洞察

- **群健康指数**: 综合评分（活跃度 + 留存率）
- **互动率分析**: 活跃用户占比统计
- **时段分析**: 深夜/日间/黄金档消息分布
- **趋势预测**: 即将推出

### 🎨 界面特性

- 暗黑主题设计，蓝紫渐变配色
- 响应式布局，支持移动端
- 流畅的动画效果
- 实时数据刷新（30 秒）

---

## 🚀 快速开始

### 使用 Docker（推荐）

```bash
# 克隆项目
git clone https://github.com/HuanLinOTO/SeekYourGuild.git
cd SeekYourGuild

# 配置环境变量
cp .env.example .env
# 编辑 .env 文件，设置 API_KEY 和 ALLOWED_GROUPS

# 启动服务
docker-compose up -d

# 查看日志
docker-compose logs -f
```

访问 http://localhost:8080 查看面板

### 单文件部署（前端嵌入）

后端二进制文件内嵌前端资源，只需一个可执行文件即可运行完整服务。

#### Windows

```powershell
# 构建
.\scripts\build.ps1

# 运行
$env:API_KEY="your-secret-key"
$env:ALLOWED_GROUPS="123456789,987654321"
.\seekyourguild.exe
```

#### Linux / macOS

```bash
# 构建
./scripts/build.sh

# 运行
API_KEY="your-secret-key" ALLOWED_GROUPS="123456789,987654321" ./seekyourguild
```

### 分离开发模式

#### 后端

```bash
cd backend

# 安装依赖
go mod download

# 设置环境变量
export API_KEY="your-secret-key"
export ALLOWED_GROUPS="123456789,987654321"

# 运行（无前端资源时只提供 API）
go run ./cmd/server/
```

#### 前端

```bash
cd frontend

# 安装依赖
yarn install

# 开发模式（连接后端 API）
yarn dev

# 构建生产版本
yarn build
```

---

## 📦 部署指南

详细部署文档请查看 [部署指南](docs/deployment.md)

### 环境变量

| 变量名           | 描述                         | 默认值                     |
| ---------------- | ---------------------------- | -------------------------- |
| `SERVER_HOST`    | 服务监听地址                 | `0.0.0.0`                  |
| `SERVER_PORT`    | 服务端口                     | `8080`                     |
| `DB_DRIVER`      | 数据库类型 (sqlite/postgres) | `sqlite`                   |
| `DB_DSN`         | 数据库连接字符串             | `file:seekyourguild.db...` |
| `API_KEY`        | 机器人上报接口密钥           | `sk-123456`                |
| `ALLOWED_GROUPS` | 允许的群号列表(逗号分隔)     | (空=所有)                  |

### Nginx 反向代理

```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://127.0.0.1:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }
}
```

---

## 📡 API 文档

### 机器人上报接口

所有上报接口需要在 Header 中携带认证信息：

```
Authorization: Bearer <API_KEY>
```

#### 成员变动上报

```http
POST /api/report/member
Content-Type: application/json

{
  "group_id": 123456789,
  "event_type": "join",  // join/leave/kick/sync
  "count": 1,
  "current_total": 100,
  "timestamp": 1732876800
}
```

#### 消息上报

```http
POST /api/report/message
Content-Type: application/json

{
  "group_id": 123456789,
  "user_id": 987654321,
  "nickname": "用户昵称",
  "message_id": "msg_123",
  "message_type": "text",
  "content_length": 10,
  "timestamp": 1732876800
}
```

#### 心跳上报

```http
POST /api/report/heartbeat
Content-Type: application/json

{
  "group_id": 123456789,
  "bot_id": 111222333,
  "status": "online",
  "member_count": 100,
  "timestamp": 1732876800
}
```

### 统计查询接口

| 接口                                         | 描述           |
| -------------------------------------------- | -------------- |
| `GET /api/groups`                            | 获取可用群列表 |
| `GET /api/stats/overview?group_id=xxx`       | 群概览数据     |
| `GET /api/stats/member-changes?group_id=xxx` | 成员变动记录   |
| `GET /api/stats/messages?group_id=xxx`       | 消息统计       |
| `GET /api/stats/active-users?group_id=xxx`   | 活跃用户排行   |

---

## 🛠️ 开发指南

### 项目结构

```
seekyourguild/
├── backend/                 # Go 后端
│   ├── cmd/server/          # 入口
│   ├── internal/
│   │   ├── config/          # 配置
│   │   ├── database/        # 数据库
│   │   ├── handlers/        # API 处理器
│   │   ├── middleware/      # 中间件
│   │   ├── models/          # 数据模型
│   │   └── static/          # 嵌入的前端资源
│   └── go.mod
├── frontend/                # Vue 前端
│   ├── src/
│   │   ├── api/             # API 调用
│   │   ├── components/      # 组件
│   │   ├── stores/          # Pinia Store
│   │   └── views/           # 页面
│   └── package.json
├── scripts/                 # 构建脚本
│   ├── build.sh             # Linux/macOS
│   └── build.ps1            # Windows
├── .github/workflows/       # CI/CD
├── Dockerfile
├── docker-compose.yml
└── README.md
```

### 技术栈

**后端**

- Go 1.21+
- Gin Web Framework
- GORM
- SQLite / PostgreSQL

**前端**

- Vue 3
- Vite 5
- TailwindCSS 3
- Pinia
- Axios

---

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

1. Fork 本项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

---

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

---

<div align="center">

**SeekYourGuild** - 让社群运营更智能

</div>
