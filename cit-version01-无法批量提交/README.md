# CIT - Go语言实现的Git模拟器

这是一个用Go语言编写的Git模拟器，主要用于教学目的，帮助理解Git的核心概念和工作原理。

## 🚀 功能特性

- **仓库管理**: 初始化Git仓库，管理仓库配置
- **文件操作**: 添加文件到暂存区，跟踪文件变更
- **提交系统**: 创建提交记录，管理版本历史
- **分支管理**: 创建、切换、删除分支
- **状态查看**: 查看工作目录、暂存区和提交状态
- **历史记录**: 查看提交历史和分支信息

## 📦 安装和运行

### 前置要求
- Go 1.21 或更高版本

### 安装依赖
```bash
go mod tidy
```

### 构建项目
```bash
# 构建命令行版本
go build -o cit-version01-无法批量提交.exe

# 构建测试版本
go build -o test.exe test.go
```

### 运行测试
```bash
# 运行功能测试
./test.exe
```

## 🛠️ 基本命令

### 初始化仓库
```bash
# 在当前目录初始化仓库
cit-version01-无法批量提交 init

# 在指定目录初始化仓库
cit-version01-无法批量提交 init /path/to/project
```

### 文件管理
```bash
# 添加文件到暂存区
cit-version01-无法批量提交 add <filename>

# 添加多个文件
cit-version01-无法批量提交 add file1.txt file2.txt

# 添加目录
cit-version01-无法批量提交 add src/
```

### 提交更改
```bash
# 提交暂存区的更改
cit-version01-无法批量提交 commit -m "提交信息"
```

### 查看状态
```bash
# 查看仓库状态
cit-version01-无法批量提交 status
```

### 查看历史
```bash
# 查看提交历史
cit-version01-无法批量提交 log
```

### 分支操作
```bash
# 列出所有分支
cit-version01-无法批量提交 branch

# 创建新分支
cit-version01-无法批量提交 branch <branch-name>

# 切换到指定分支
cit-version01-无法批量提交 checkout <branch-name>
```

## 🏗️ 项目结构

```
cit/
├── cmd/                    # 命令行界面
│   ├── root.go           # 根命令
│   ├── init.go           # 初始化命令
│   ├── add.go            # 添加命令
│   ├── commit.go         # 提交命令
│   ├── status.go         # 状态命令
│   ├── log.go            # 日志命令
│   ├── branch.go         # 分支命令
│   └── checkout.go       # 切换命令
├── internal/              # 内部包
│   ├── git/              # Git核心逻辑
│   │   ├── repository.go # 仓库管理
│   │   └── models.go     # 数据模型
│   ├── storage/          # 数据存储
│   │   └── storage.go    # 存储实现
│   └── utils/            # 工具函数
│       └── utils.go      # 通用工具
├── pkg/                  # 公共包
├── main.go               # 主程序入口
├── test.go               # 功能测试
└── go.mod                # Go模块文件
```

## 📁 仓库结构

当初始化仓库后，会在项目目录下创建 `.cit` 目录：

```
.cit/
├── objects/              # 对象存储
│   ├── [hash1]/         # 按哈希分组的对象
│   └── [hash2]/
├── refs/                 # 引用管理
│   ├── heads/           # 分支引用
│   └── tags/            # 标签引用
├── repository.json       # 仓库配置
├── branches.json         # 分支信息
├── commits.json          # 提交历史
└── staging.json          # 暂存区状态
```

## 🔍 核心概念

### 1. 对象存储
- **Blob对象**: 存储文件内容
- **Tree对象**: 存储目录结构
- **Commit对象**: 存储提交信息

### 2. 暂存区
- 工作目录和提交之间的中间状态
- 通过 `add` 命令添加文件
- 通过 `commit` 命令提交更改

### 3. 分支系统
- 每个分支指向一个提交
- 支持创建、切换、删除操作
- 默认分支为 `main`

### 4. 提交历史
- 每个提交都有唯一的哈希ID
- 支持父子关系，形成提交链
- 包含作者、时间、消息等信息

## 🧪 使用示例

### 基本工作流程
```bash
# 1. 初始化仓库
cit-version01-无法批量提交 init

# 2. 创建文件
echo "Hello, CIT!" > hello.txt

# 3. 添加到暂存区
cit-version01-无法批量提交 add hello.txt

# 4. 提交更改
cit-version01-无法批量提交 commit -m "添加问候文件"

# 5. 查看状态
cit-version01-无法批量提交 status

# 6. 查看历史
cit-version01-无法批量提交 log
```

### 分支操作示例
```bash
# 创建功能分支
cit-version01-无法批量提交 branch feature-login

# 切换到功能分支
cit-version01-无法批量提交 checkout feature-login

# 在新分支上工作
echo "登录功能" > login.txt
cit-version01-无法批量提交 add login.txt
cit-version01-无法批量提交 commit -m "实现登录功能"

# 切换回主分支
cit-version01-无法批量提交 checkout main
```

## 🎯 学习目标

通过这个项目，你可以学习：

1. **Git工作原理**: 理解版本控制系统的核心概念
2. **Go语言编程**: 练习文件操作、数据结构、错误处理
3. **命令行开发**: 学习使用Cobra框架开发CLI工具
4. **系统设计**: 理解数据存储、对象模型、状态管理

## ⚠️ 注意事项

- 这是一个教学项目，功能相对简化
- 不建议用于生产环境
- 主要用于学习和理解Git的工作原理
- 支持基本的文件操作和版本管理

## 🔧 开发计划

- [ ] 实现树对象和Blob对象
- [ ] 添加合并功能
- [ ] 支持标签管理
- [ ] 实现远程仓库操作
- [ ] 添加配置文件支持
- [ ] 优化性能和存储效率

## 📝 许可证

本项目仅用于教学目的，请勿用于商业用途。

## 🤝 贡献

欢迎提交Issue和Pull Request来改进这个项目！

---

**享受学习Git的乐趣！** 🎉
