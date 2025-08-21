# CIT Git模拟器 - GitHub真实推送指南

## 🎯 目标

让CIT Git模拟器能够**真正推送到GitHub仓库**，而不仅仅是模拟！

## 🚀 新功能特性

### 1. **真正的GitHub推送**
- 使用GitHub REST API v3
- 支持文件内容更新
- 自动处理文件冲突
- 真实的网络传输

### 2. **智能推送模式**
- **模拟模式**: `cit push origin main` (教学用途)
- **真实模式**: `cit push --github-token TOKEN origin main` (实际推送)

### 3. **GitHub API集成**
- 自动解析GitHub仓库信息
- 支持HTTPS和SSH URL格式
- 智能文件内容管理

## 🔑 准备工作

### 1. **创建GitHub仓库**
```bash
# 在GitHub上创建仓库
# 仓库名: cit
# 所有者: code-100-precent
# 完整URL: https://github.com/code-100-precent/cit.git
```

### 2. **生成个人访问令牌**
1. 访问 [GitHub Settings > Developer settings > Personal access tokens](https://github.com/settings/tokens)
2. 点击 "Generate new token (classic)"
3. 选择权限：
   - ✅ `repo` (完整的仓库访问权限)
   - ✅ `workflow` (可选，用于GitHub Actions)
4. 生成令牌并**安全保存**

### 3. **验证仓库权限**
确保你有对 `code-100-precent/cit` 仓库的推送权限。

## 🛠️ 使用方法

### 1. **基本推送流程**
```bash
# 1. 初始化仓库
cit init

# 2. 添加文件
echo "Hello GitHub!" > test.txt
cit add test.txt

# 3. 提交更改
cit commit -m "feat: 添加测试文件"

# 4. 添加远程仓库
cit remote add origin https://github.com/code-100-precent/cit.git

# 5. 推送到GitHub (真实推送)
cit push --github-token YOUR_TOKEN origin main
```

### 2. **命令参数说明**
```bash
cit push [远程名] [分支名] --github-token [令牌]

# 参数说明：
# 远程名: 远程仓库名称 (默认: origin)
# 分支名: 要推送的分支 (默认: 当前分支)
# --github-token: GitHub个人访问令牌
```

### 3. **URL格式支持**
```bash
# HTTPS格式
cit remote add origin https://github.com/code-100-precent/cit.git

# SSH格式 (需要SSH密钥配置)
cit remote add origin git@github.com:code-100-precent/cit.git
```

## 🔧 技术实现

### 1. **GitHub API集成**
```go
// 创建GitHub API客户端
api := NewGitHubAPI(token)

// 推送文件到GitHub
func (r *Repository) PushToGitHub(remoteName, branchName, token string) error {
    // 1. 解析远程仓库信息
    repoInfo, err := r.parseGitHubRemote(remoteName)
    
    // 2. 创建API客户端
    api := NewGitHubAPI(token)
    
    // 3. 推送每个文件
    for filePath, hash := range staging {
        r.pushFileToGitHub(api, repoInfo, branchName, filePath, hash)
    }
}
```

### 2. **文件内容管理**
```go
// 读取文件并编码为base64
func (r *Repository) readFileContent(filePath string) (string, error) {
    content, err := os.ReadFile(fullPath)
    encoded := base64.StdEncoding.EncodeToString(content)
    return encoded, nil
}

// 推送文件到GitHub
func (r *Repository) pushFileToGitHub(api *GitHubAPI, repo *GitHubRepo, branch, filePath, hash string) error {
    // 1. 检查文件是否已存在
    existingContent, err := api.getFileContent(repo.FullName, filePath, branch)
    
    // 2. 准备请求数据
    requestBody := map[string]interface{}{
        "message": fmt.Sprintf("Update %s via CIT", filePath),
        "content": fileContent,
        "branch":  branch,
    }
    
    // 3. 如果文件存在，提供SHA
    if existingContent.SHA != "" {
        requestBody["sha"] = existingContent.SHA
    }
    
    // 4. 发送PUT请求到GitHub API
    resp, err := api.HTTPClient.Do(req)
}
```

### 3. **智能冲突处理**
- 自动检测文件是否已存在
- 提供文件SHA避免冲突
- 支持增量更新

## 📊 推送流程详解

### 1. **准备阶段**
```
本地仓库 → 解析远程信息 → 验证GitHub连接
   ↓              ↓              ↓
暂存区文件    仓库URL解析    API权限检查
```

### 2. **传输阶段**
```
文件读取 → Base64编码 → GitHub API → 文件更新
   ↓           ↓           ↓          ↓
本地文件    编码内容    HTTP PUT    远程存储
```

### 3. **完成阶段**
```
推送完成 → 状态报告 → 错误处理
   ↓          ↓          ↓
成功确认    详细日志    问题诊断
```

## ⚠️ 注意事项

### 1. **安全考虑**
- **永远不要**在代码中硬编码令牌
- 使用环境变量或配置文件存储令牌
- 定期轮换访问令牌
- 最小权限原则

### 2. **API限制**
- GitHub API有速率限制
- 认证用户: 5000次/小时
- 未认证用户: 60次/小时
- 大文件推送可能需要特殊处理

### 3. **错误处理**
- 网络连接问题
- 权限不足
- 文件冲突
- API限制

## 🔍 故障排除

### 1. **常见错误**
```bash
# 权限错误
GitHub API错误: 401 - Bad credentials
解决方案: 检查令牌是否正确，是否有足够权限

# 仓库不存在
GitHub API错误: 404 - Not Found
解决方案: 检查仓库URL是否正确

# 文件冲突
GitHub API错误: 409 - Conflict
解决方案: 先拉取最新更改，解决冲突后重新推送
```

### 2. **调试技巧**
```bash
# 启用详细输出
cit push --github-token TOKEN origin main -v

# 检查远程配置
cit remote list

# 验证GitHub连接
curl -H "Authorization: token YOUR_TOKEN" https://api.github.com/user
```

### 3. **网络问题**
- 检查网络连接
- 验证防火墙设置
- 使用代理服务器（如需要）

## 🎯 实际应用场景

### 1. **个人项目管理**
```bash
# 日常开发工作流
cit add .
cit commit -m "feat: 新功能"
cit push --github-token TOKEN origin main
```

### 2. **教学演示**
```bash
# 展示真实的Git工作流
cit push origin main          # 模拟推送
cit push --github-token TOKEN origin main  # 真实推送
```

### 3. **自动化脚本**
```bash
# 在CI/CD中使用
export GITHUB_TOKEN="your_token"
cit push --github-token $GITHUB_TOKEN origin main
```

## 🔮 未来扩展

### 1. **更多GitHub功能**
- 创建Pull Request
- 管理Issues
- 处理Webhooks
- 支持GitHub Actions

### 2. **其他Git服务**
- GitLab支持
- Bitbucket支持
- 自托管Git服务

### 3. **高级功能**
- 批量文件推送
- 增量同步
- 冲突自动解决
- 推送历史记录

## 📚 学习资源

### 1. **GitHub API文档**
- [GitHub REST API v3](https://docs.github.com/en/rest)
- [Contents API](https://docs.github.com/en/rest/repos/contents)
- [Authentication](https://docs.github.com/en/rest/overview/authentication)

### 2. **Go网络编程**
- [net/http包](https://golang.org/pkg/net/http/)
- [HTTP客户端](https://golang.org/pkg/net/http/#Client)
- [JSON处理](https://golang.org/pkg/encoding/json/)

### 3. **Git内部原理**
- [Git Internals](https://git-scm.com/book/en/v2/Git-Internals-Plumbing-and-Porcelain)
- [Git Protocol](https://git-scm.com/docs/protocol-v2)

## 🎉 总结

通过这个功能，CIT Git模拟器现在可以：

1. **真正推送到GitHub** - 不再是模拟！
2. **学习真实Git工作流** - 体验完整的开发流程
3. **理解GitHub API** - 了解版本控制的网络层面
4. **实践自动化** - 为CI/CD流程做准备

### 🚀 下一步行动

1. **运行演示脚本**: `github_push_demo.bat`
2. **获取GitHub令牌**: 按照指南生成访问令牌
3. **测试真实推送**: `cit push --github-token YOUR_TOKEN origin main`
4. **检查GitHub仓库**: 验证文件是否成功推送

现在你的CIT Git模拟器已经是一个**功能完整的版本控制系统**了！🎯

---

**注意**: 这是教学项目，生产环境请使用标准的Git工具。
