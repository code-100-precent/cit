package git

import (
	"bytes"
	"cit/internal/storage"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// GitHubAPI 表示GitHub API客户端
type GitHubAPI struct {
	Token      string
	HTTPClient *http.Client
}

// GitHubRepo 表示GitHub仓库信息
type GitHubRepo struct {
	FullName      string `json:"full_name"`
	CloneURL      string `json:"clone_url"`
	SSHURL        string `json:"ssh_url"`
	DefaultBranch string `json:"default_branch"`
}

// GitHubContent 表示GitHub文件内容
type GitHubContent struct {
	Content  string `json:"content"`
	Encoding string `json:"encoding"`
	SHA      string `json:"sha"`
}

// NewGitHubAPI 创建GitHub API客户端
func NewGitHubAPI(token string) *GitHubAPI {
	return &GitHubAPI{
		Token: token,
		HTTPClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// PushToGitHub 真正推送到GitHub
func (r *Repository) PushToGitHub(remoteName, branchName, token string) error {
	// 解析GitHub仓库信息
	repoInfo, err := r.parseGitHubRemote(remoteName)
	if err != nil {
		return fmt.Errorf("解析GitHub远程仓库失败: %v", err)
	}

	// 创建GitHub API客户端
	api := NewGitHubAPI(token)

	// 测试GitHub连接
	fmt.Printf("正在测试GitHub连接...\n")
	if err := api.TestConnection(); err != nil {
		return fmt.Errorf("GitHub连接测试失败: %v", err)
	}

	// 获取当前分支的提交
	commits, err := r.Storage.GetCommitHistory()
	if err != nil {
		return fmt.Errorf("获取提交历史失败: %v", err)
	}

	if len(commits) == 0 {
		return fmt.Errorf("没有可推送的提交")
	}

	fmt.Printf("正在推送到GitHub仓库: %s\n", repoInfo.FullName)
	fmt.Printf("推送分支: %s\n", branchName)
	fmt.Printf("推送提交数: %d\n", len(commits))

	// 从工作目录获取文件列表
	files, err := r.getWorkingDirectoryFiles()
	if err != nil {
		return fmt.Errorf("获取工作目录文件失败: %v", err)
	}

	if len(files) == 0 {
		return fmt.Errorf("工作目录为空，没有文件可推送")
	}

	fmt.Printf("准备推送 %d 个文件...\n", len(files))

	// 推送每个文件到GitHub
	successCount := 0
	for _, filePath := range files {
		fmt.Printf("正在推送文件: %s\n", filePath)
		if err := r.pushFileToGitHub(api, repoInfo, branchName, filePath, ""); err != nil {
			fmt.Printf("❌ 推送文件 %s 失败: %v\n", filePath, err)
			return fmt.Errorf("推送失败，停止推送: %v", err)
		} else {
			successCount++
		}
	}

	if successCount == len(files) {
		fmt.Printf("✅ 成功推送到GitHub仓库: %s\n", repoInfo.FullName)
		fmt.Printf("✅ 成功推送 %d/%d 个文件\n", successCount, len(files))
		return nil
	} else {
		return fmt.Errorf("部分推送失败: %d/%d 个文件成功", successCount, len(files))
	}
}

// parseGitHubRemote 解析GitHub远程仓库信息
func (r *Repository) parseGitHubRemote(remoteName string) (*GitHubRepo, error) {
	remotes, err := r.Storage.ListRemotes()
	if err != nil {
		return nil, err
	}

	var targetRemote *storage.Remote
	for _, remote := range remotes {
		if remote.Name == remoteName {
			targetRemote = remote
			break
		}
	}

	if targetRemote == nil {
		return nil, fmt.Errorf("远程仓库 '%s' 不存在", remoteName)
	}

	// 从URL解析仓库信息
	// 支持格式: https://github.com/owner/repo.git
	// 或者: git@github.com:owner/repo.git
	url := targetRemote.URL
	var fullName string

	if strings.Contains(url, "github.com") {
		if strings.HasPrefix(url, "https://") {
			// https://github.com/owner/repo.git
			parts := strings.Split(url, "/")
			if len(parts) >= 4 {
				fullName = parts[3] + "/" + strings.TrimSuffix(parts[4], ".git")
			}
		} else if strings.HasPrefix(url, "git@") {
			// git@github.com:owner/repo.git
			parts := strings.Split(url, ":")
			if len(parts) == 2 {
				repoParts := strings.Split(parts[1], "/")
				if len(repoParts) >= 2 {
					fullName = repoParts[0] + "/" + strings.TrimSuffix(repoParts[1], ".git")
				}
			}
		}
	}

	if fullName == "" {
		return nil, fmt.Errorf("无法解析GitHub仓库信息: %s", url)
	}

	return &GitHubRepo{
		FullName: fullName,
		CloneURL: url,
	}, nil
}

// pushFileToGitHub 推送单个文件到GitHub
func (r *Repository) pushFileToGitHub(api *GitHubAPI, repo *GitHubRepo, branch, filePath, hash string) error {
	// 读取文件内容
	fileContent, err := r.readFileContent(filePath)
	if err != nil {
		return fmt.Errorf("读取文件内容失败: %v", err)
	}

	// 标准化文件路径，确保GitHub API能正确处理文件夹结构
	normalizedPath := r.normalizeFilePath(filePath)

	// 准备API请求
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s", repo.FullName, normalizedPath)
	
	// 检查文件是否已存在
	existingContent, err := api.getFileContent(repo.FullName, normalizedPath, branch)
	
	requestBody := map[string]interface{}{
		"message": fmt.Sprintf("Update %s via CIT", normalizedPath),
		"content": fileContent,
		"branch":  branch,
	}

	// 如果文件已存在，需要提供SHA
	if err == nil && existingContent.SHA != "" {
		requestBody["sha"] = existingContent.SHA
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("序列化请求失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("PUT", apiURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", "token "+api.Token)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	// 发送请求
	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("发送请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	body, _ := io.ReadAll(resp.Body)
	
	// 详细的错误处理和调试信息
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("GitHub API错误: %s\nURL: %s\n响应: %s", resp.Status, apiURL, string(body))
	}

	// 解析响应确认推送成功
	var result map[string]interface{}
	if err := json.Unmarshal(body, &result); err != nil {
		return fmt.Errorf("解析响应失败: %v", err)
	}

	// 检查响应中的关键字段
	if commitSHA, ok := result["commit"].(map[string]interface{}); ok {
		if sha, ok := commitSHA["sha"].(string); ok {
			fmt.Printf("✓ 文件 %s 推送成功，提交SHA: %s\n", filePath, sha[:8])
		}
	}

	return nil
}

// getFileContent 获取GitHub文件内容
func (api *GitHubAPI) getFileContent(repoFullName, filePath, branch string) (*GitHubContent, error) {
	apiURL := fmt.Sprintf("https://api.github.com/repos/%s/contents/%s?ref=%s", repoFullName, filePath, branch)
	
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "token "+api.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// 文件不存在
		return nil, fmt.Errorf("文件不存在")
	}

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API错误: %s - %s", resp.Status, string(body))
	}

	var content GitHubContent
	if err := json.NewDecoder(resp.Body).Decode(&content); err != nil {
		return nil, err
	}

	return &content, nil
}

// TestConnection 测试GitHub连接
func (api *GitHubAPI) TestConnection() error {
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return fmt.Errorf("创建请求失败: %v", err)
	}

	req.Header.Set("Authorization", "token "+api.Token)
	req.Header.Set("Accept", "application/vnd.github.v3+json")

	resp, err := api.HTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("连接失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("认证失败: %s - %s", resp.Status, string(body))
	}

	var user map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&user); err != nil {
		return fmt.Errorf("解析用户信息失败: %v", err)
	}

	if username, ok := user["login"].(string); ok {
		fmt.Printf("✅ GitHub连接成功，用户: %s\n", username)
	}

	return nil
}

// readFileContent 读取文件内容并编码为base64
func (r *Repository) readFileContent(filePath string) (string, error) {
	// 构建完整的文件路径
	fullPath := filepath.Join(r.Path, filePath)

	// 读取文件内容
	content, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("读取文件失败: %v", err)
	}

	// 编码为base64
	encoded := base64.StdEncoding.EncodeToString(content)
	return encoded, nil
}

// normalizeFilePath 标准化文件路径，确保GitHub API能正确处理
func (r *Repository) normalizeFilePath(filePath string) string {
	// 将Windows路径分隔符转换为GitHub API使用的正斜杠
	normalized := strings.ReplaceAll(filePath, "\\", "/")
	return normalized
}
