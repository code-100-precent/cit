package git

// Status 表示仓库状态
type Status struct {
	CurrentBranch  string   `json:"current_branch"`
	LastCommit     string   `json:"last_commit"`
	StagedFiles   []string `json:"staged_files"`
	ModifiedFiles []string `json:"modified_files"`
	UntrackedFiles []string `json:"untracked_files"`
}

// WorkdirStatus 表示工作目录状态
type WorkdirStatus struct {
	ModifiedFiles  []string `json:"modified_files"`
	UntrackedFiles []string `json:"untracked_files"`
}

// Remote 表示远程仓库
type Remote struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

// PushResult 表示推送结果
type PushResult struct {
	RemoteName string   `json:"remote_name"`
	BranchName string   `json:"branch_name"`
	PushedCommits []string `json:"pushed_commits"`
	TotalObjects int    `json:"total_objects"`
}
