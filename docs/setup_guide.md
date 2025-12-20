# 开发环境搭建指南

## 1. 环境要求

### 1.1 操作系统

- Windows 10 或 Windows 11
- 建议使用 64 位系统

### 1.2 开发工具

- **IDE**：Visual Studio Code 或 GoLand + PyCharm
- **版本控制**：Git
- **数据库工具**：Navicat 或 MySQL Workbench
- **Redis 工具**：Redis Desktop Manager

### 1.3 语言环境

#### 1.3.1 Go 环境

- **Go 版本**：1.20+
- **安装路径**：建议安装到 `C:\Go` 或其他不含空格的路径
- **环境变量**：
  - `GOROOT`：Go 安装路径
  - `GOPATH`：Go 工作目录，建议设置为 `%USERPROFILE%\go`
  - `PATH`：添加 `%GOROOT%\bin` 和 `%GOPATH%\bin`

#### 1.3.2 Python 环境

- **Python 版本**：3.9+
- **安装路径**：建议安装到 `C:\Python39` 或其他不含空格的路径
- **环境变量**：
  - `PATH`：添加 Python 安装路径和 Scripts 目录
- **包管理工具**：
  - pip 或 pip3
  - 建议使用虚拟环境

### 1.4 数据库环境

#### 1.4.1 MySQL

- **版本**：8.0+
- **安装**：建议使用 MySQL Installer 进行安装
- **配置**：
  - 端口：3306
  - 用户名：root
  - 密码：建议设置强密码

#### 1.4.2 Redis

- **版本**：7.0+
- **安装**：建议使用 Redis for Windows 或 WSL2 下的 Redis
- **配置**：
  - 端口：6379
  - 密码：建议设置密码

## 2. 项目初始化

### 2.1 克隆项目

```bash
git clone <项目仓库地址>
cd DdOaListDownload
```

### 2.2 目录结构

```
DdOaListDownload/
├── README.md              # 项目说明文档
├── backend/               # 后端目录
│   ├── go/                # Go 服务（处理新版 API）
│   └── python/             # Python 服务（处理旧版 API）
├── frontend/              # 前端目录
└── docs/                 # 文档目录
    ├── setup_guide.md      # 开发环境搭建指南
    ├── development_guide.md  # 开发规范指南
    └── implementation_plan.md  # 实施方案
```

## 3. 后端环境搭建

### 3.1 Go 服务环境

#### 3.1.1 安装依赖

```bash
cd backend/go
go mod tidy
```

#### 3.1.2 配置文件

- 复制 `config/config.example.go` 为 `config/config.go`
- 修改配置文件中的数据库、Redis 等连接信息

#### 3.1.3 运行服务

```bash
cd backend/go
go run main.go
```

### 3.2 Python 服务环境

#### 3.2.1 创建虚拟环境

```bash
cd backend/python
python -m venv venv
```

#### 3.2.2 激活虚拟环境

```bash
# Windows CMD
venv\Scripts\activate.bat

# Windows PowerShell
venv\Scripts\Activate.ps1
```

#### 3.2.3 安装依赖

```bash
pip install -r requirements.txt
```

#### 3.2.4 配置文件

- 复制 `config/config.example.py` 为 `config/config.py`
- 修改配置文件中的数据库、Redis 等连接信息

#### 3.2.5 运行服务

```bash
python app.py
```

## 4. 前端环境搭建

### 4.1 安装依赖

```bash
cd frontend
npm install
```

### 4.2 配置文件

- 复制 `.env.example` 为 `.env`
- 修改配置文件中的 API 地址等信息

### 4.3 运行服务

```bash
npm run dev
```

## 5. 数据库初始化

### 5.1 创建数据库

```sql
CREATE DATABASE dd_oa_download CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 5.2 执行迁移

#### 5.2.1 Go 服务迁移

```bash
cd backend/go
# 执行数据库迁移
```

#### 5.2.2 Python 服务迁移

```bash
cd backend/python
# 执行数据库迁移
```

## 6. 访问应用

- 前端应用：http://localhost:5173
- Go 后端 API：http://localhost:8080
- Python 后端 API：http://localhost:8081

## 7. 常见问题

### 7.1 Go 依赖安装失败

- 检查网络连接
- 配置 GOPROXY：
  ```bash
go env -w GOPROXY=https://goproxy.io,direct
  ```

### 7.2 Python 依赖安装失败

- 升级 pip：
  ```bash
  pip install --upgrade pip
  ```
- 使用国内镜像：
  ```bash
  pip install -r requirements.txt -i https://pypi.tuna.tsinghua.edu.cn/simple
  ```

### 7.3 数据库连接失败

- 检查 MySQL 服务是否启动
- 检查用户名、密码是否正确
- 检查数据库是否存在

### 7.4 Redis 连接失败

- 检查 Redis 服务是否启动
- 检查密码是否正确

## 8. 开发流程

1. 创建分支
2. 开发功能
3. 编写测试
4. 提交代码
5. 创建 PR
6. 代码 review
7. 合并到主分支

作者: cjx
邮箱: xx4125517@126.com
时间: 2025-12-20 14:30:00