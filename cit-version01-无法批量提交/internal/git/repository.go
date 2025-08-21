package git

import (
	"crypto/sha1"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"cit/internal/storage"
	"cit/internal/utils"
)

// Repository 表示一个Git仓库
type Repository struct {
	ID            string    `json:"id"`
	Path          string    `json:"path"`
	CreatedAt     time.Time `json:"created_at"`
	CurrentBranch string    `json:"current_branch"`
	Storage       *storage.Storage
}

// InitRepository 初始化一个新的Git仓库
func InitRepository(path string) (*Repository, error) {
	// 生成仓库ID
	repoID := generateRepositoryID(path)

	// 创建仓库目录结构
	gitDir := filepath.Join(path, ".cit-version01-无法批量提交")
	if err := os.MkdirAll(gitDir, 0755); err != nil {
		return nil, fmt.Errorf("创建仓库目录失败: %v", err)
	}

	// 创建子目录
	subdirs := []string{"objects", "refs", "refs/heads", "refs/tags"}
	for _, subdir := range subdirs {
		subdirPath := filepath.Join(gitDir, subdir)
		if err := os.MkdirAll(subdirPath, 0755); err != nil {
			return nil, fmt.Errorf("创建子目录失败: %v", err)
		}
	}

	// 初始化存储
	storage, err := storage.NewStorage(gitDir)
	if err != nil {
		return nil, fmt.Errorf("初始化存储失败: %v", err)
	}

	// 创建仓库对象
	repo := &Repository{
		ID:            repoID,
		Path:          path,
		CreatedAt:     time.Now(),
		CurrentBranch: "main",
		Storage:       storage,
	}

	// 创建主分支
	if err := repo.CreateBranch("main"); err != nil {
		return nil, fmt.Errorf("创建主分支失败: %v", err)
	}

	// 保存仓库信息
	if err := repo.save(); err != nil {
		return nil, fmt.Errorf("保存仓库信息失败: %v", err)
	}

	return repo, nil
}

// FindRepository 查找Git仓库
func FindRepository(path string) (*Repository, error) {
	currentPath, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	// 向上查找.cit目录
	for {
		gitDir := filepath.Join(currentPath, ".cit-version01-无法批量提交")
		if _, err := os.Stat(gitDir); err == nil {
			// 找到仓库，加载信息
			return loadRepository(gitDir)
		}

		parent := filepath.Dir(currentPath)
		if parent == currentPath {
			break
		}
		currentPath = parent
	}

	return nil, fmt.Errorf("未找到Git仓库")
}

// AddToStaging 将文件添加到暂存区
func (r *Repository) AddToStaging(filePath string) error {
	// 获取文件的相对路径
	relPath, err := filepath.Rel(r.Path, filePath)
	if err != nil {
		return fmt.Errorf("获取相对路径失败: %v", err)
	}

	// 计算文件哈希
	hash, err := utils.CalculateFileHash(filePath)
	if err != nil {
		return fmt.Errorf("计算文件哈希失败: %v", err)
	}

	// 存储文件对象
	if err := r.Storage.StoreObject(hash, filePath); err != nil {
		return fmt.Errorf("存储文件对象失败: %v", err)
	}

	// 添加到暂存区
	return r.Storage.AddToStaging(relPath, hash)
}

// Commit 提交暂存区的更改
func (r *Repository) Commit(message string) (*storage.Commit, error) {
	// 获取暂存区内容
	staging, err := r.Storage.GetStaging()
	if err != nil {
		return nil, fmt.Errorf("获取暂存区失败: %v", err)
	}

	if len(staging) == 0 {
		return nil, fmt.Errorf("暂存区为空")
	}

	// 获取当前分支的最新提交
	var parentID string
	if currentCommit, err := r.Storage.GetBranchHead(r.CurrentBranch); err == nil {
		parentID = currentCommit
	}

	// 创建提交对象
	commit := &storage.Commit{
		ID:        utils.GenerateID(),
		Message:   message,
		Author:    getCurrentUser(),
		Timestamp: time.Now(),
		ParentID:  parentID,
		TreeHash:  "", // 这里应该计算树对象的哈希
	}

	// 保存提交对象
	if err := r.Storage.StoreCommit(commit); err != nil {
		return nil, fmt.Errorf("保存提交失败: %v", err)
	}

	// 更新分支头
	if err := r.Storage.UpdateBranchHead(r.CurrentBranch, commit.ID); err != nil {
		return nil, fmt.Errorf("更新分支头失败: %v", err)
	}

	// 清空暂存区
	if err := r.Storage.ClearStaging(); err != nil {
		return nil, fmt.Errorf("清空暂存区失败: %v", err)
	}

	return commit, nil
}

// GetStatus 获取仓库状态
func (r *Repository) GetStatus() (*Status, error) {
	status := &Status{
		CurrentBranch: r.CurrentBranch,
	}

	// 获取最新提交
	if commit, err := r.Storage.GetBranchHead(r.CurrentBranch); err == nil {
		status.LastCommit = commit
	}

	// 获取暂存区文件
	if staging, err := r.Storage.GetStaging(); err == nil {
		status.StagedFiles = make([]string, 0, len(staging))
		for file := range staging {
			status.StagedFiles = append(status.StagedFiles, file)
		}
	}

	// 获取工作目录状态
	workdirStatus, err := r.getWorkdirStatus()
	if err == nil {
		status.ModifiedFiles = workdirStatus.ModifiedFiles
		status.UntrackedFiles = workdirStatus.UntrackedFiles
	}

	return status, nil
}

// GetCommitHistory 获取提交历史
func (r *Repository) GetCommitHistory() ([]*storage.Commit, error) {
	return r.Storage.GetCommitHistory()
}

// ListBranches 列出所有分支
func (r *Repository) ListBranches() ([]*storage.Branch, error) {
	return r.Storage.ListBranches()
}

// CreateBranch 创建新分支
func (r *Repository) CreateBranch(name string) error {
	// 检查分支是否已存在
	branches, err := r.Storage.ListBranches()
	if err != nil {
		return err
	}

	for _, branch := range branches {
		if branch.Name == name {
			return fmt.Errorf("分支 '%s' 已存在", name)
		}
	}

	// 获取当前分支头
	var head string
	if currentCommit, err := r.Storage.GetBranchHead(r.CurrentBranch); err == nil {
		head = currentCommit
	}

	// 创建新分支
	branch := &storage.Branch{
		Name: name,
		Head: head,
	}

	return r.Storage.CreateBranch(branch)
}

// CheckoutBranch 切换到指定分支
func (r *Repository) CheckoutBranch(name string) error {
	// 检查分支是否存在
	branches, err := r.Storage.ListBranches()
	if err != nil {
		return err
	}

	var targetBranch *storage.Branch
	for _, branch := range branches {
		if branch.Name == name {
			targetBranch = branch
			break
		}
	}

	if targetBranch == nil {
		return fmt.Errorf("分支 '%s' 不存在", name)
	}

	// 更新当前分支
	r.CurrentBranch = name
	return r.save()
}

// GetCurrentBranch 获取当前分支名
func (r *Repository) GetCurrentBranch() string {
	return r.CurrentBranch
}

// IsStagingEmpty 检查暂存区是否为空
func (r *Repository) IsStagingEmpty() bool {
	staging, err := r.Storage.GetStaging()
	if err != nil {
		return true
	}
	return len(staging) == 0
}

// AddRemote 添加远程仓库
func (r *Repository) AddRemote(name, url string) error {
	return r.Storage.AddRemote(&storage.Remote{
		Name: name,
		URL:  url,
	})
}

// ListRemotes 列出所有远程仓库
func (r *Repository) ListRemotes() ([]*storage.Remote, error) {
	return r.Storage.ListRemotes()
}

// RemoveRemote 删除远程仓库
func (r *Repository) RemoveRemote(name string) error {
	return r.Storage.RemoveRemote(name)
}

// Push 推送到远程仓库
func (r *Repository) Push(remoteName, branchName string) error {
	// 获取远程仓库信息
	remotes, err := r.Storage.ListRemotes()
	if err != nil {
		return fmt.Errorf("获取远程仓库失败: %v", err)
	}

	var targetRemote *storage.Remote
	for _, remote := range remotes {
		if remote.Name == remoteName {
			targetRemote = remote
			break
		}
	}

	if targetRemote == nil {
		return fmt.Errorf("远程仓库 '%s' 不存在", remoteName)
	}

	// 获取当前分支的提交
	commits, err := r.Storage.GetCommitHistory()
	if err != nil {
		return fmt.Errorf("获取提交历史失败: %v", err)
	}

	if len(commits) == 0 {
		return fmt.Errorf("没有可推送的提交")
	}

	// 模拟推送过程
	fmt.Printf("正在推送到 %s (%s)...\n", remoteName, targetRemote.URL)
	fmt.Printf("推送分支: %s\n", branchName)
	fmt.Printf("推送提交数: %d\n", len(commits))

	// 在实际实现中，这里会：
	// 1. 连接到远程仓库
	// 2. 计算需要推送的对象
	// 3. 通过网络传输数据
	// 4. 更新远程引用

	return nil
}

// 私有方法

func (r *Repository) save() error {
	repoFile := filepath.Join(r.Path, ".cit-version01-无法批量提交", "repository.json")
	data, err := json.MarshalIndent(r, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(repoFile, data, 0644)
}

func loadRepository(gitDir string) (*Repository, error) {
	repoFile := filepath.Join(gitDir, "repository.json")
	data, err := os.ReadFile(repoFile)
	if err != nil {
		return nil, err
	}

	var repo Repository
	if err := json.Unmarshal(data, &repo); err != nil {
		return nil, err
	}

	// 重新初始化存储
	storage, err := storage.NewStorage(gitDir)
	if err != nil {
		return nil, err
	}
	repo.Storage = storage

	return &repo, nil
}

func generateRepositoryID(path string) string {
	data := fmt.Sprintf("%s-%d", path, time.Now().UnixNano())
	hash := sha1.Sum([]byte(data))
	return fmt.Sprintf("%x", hash[:8])
}

func getCurrentUser() string {
	// 这里应该从环境变量或配置文件获取用户信息
	// 简化实现，返回默认值
	return "user@example.com"
}

func (r *Repository) getWorkdirStatus() (*WorkdirStatus, error) {
	// 这里应该实现工作目录状态检查
	// 简化实现，返回空状态
	return &WorkdirStatus{
		ModifiedFiles:  []string{},
		UntrackedFiles: []string{},
	}, nil
}

// getWorkingDirectoryFiles 获取工作目录中的所有文件
func (r *Repository) getWorkingDirectoryFiles() ([]string, error) {
	var files []string

	// 遍历工作目录，收集所有文件
	err := filepath.Walk(r.Path, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录
		if info.IsDir() {
			return nil
		}

		// 跳过.cit目录
		if strings.Contains(path, ".cit-version01-无法批量提交") {
			return nil
		}

		// 跳过隐藏文件
		if strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}

		// 获取相对路径
		relPath, err := filepath.Rel(r.Path, path)
		if err != nil {
			return err
		}

		// 跳过可执行文件
		if strings.HasSuffix(path, ".exe") {
			return nil
		}

		files = append(files, relPath)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("遍历工作目录失败: %v", err)
	}

	return files, nil
}
