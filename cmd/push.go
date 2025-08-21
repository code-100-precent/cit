package cmd

import (
	"fmt"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push [远程名] [分支名]",
	Short: "推送提交到远程仓库",
	Long:  "将本地分支的提交推送到远程仓库",
	Args:  cobra.MaximumNArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 获取远程名和分支名
		remoteName := "origin"
		branchName := repo.GetCurrentBranch()

		if len(args) > 0 {
			remoteName = args[0]
		}
		if len(args) > 1 {
			branchName = args[1]
		}

		// 检查是否是GitHub推送
		githubToken, _ := cmd.Flags().GetString("github-token")
		if githubToken != "" {
			// 使用GitHub API推送
			if err := repo.PushToGitHub(remoteName, branchName, githubToken); err != nil {
				return fmt.Errorf("GitHub推送失败: %v", err)
			}
		} else {
			// 使用模拟推送
			if err := repo.Push(remoteName, branchName); err != nil {
				return fmt.Errorf("推送失败: %v", err)
			}
			fmt.Printf("成功推送到 %s/%s\n", remoteName, branchName)
		}

		return nil
	},
}

func init() {
	pushCmd.Flags().String("github-token", "", "GitHub个人访问令牌 (用于真实推送)")
}

func init() {
	// 添加远程仓库管理命令到remoteCmd
	remoteCmd.AddCommand(addRemoteCmd)
	remoteCmd.AddCommand(listRemoteCmd)
	remoteCmd.AddCommand(removeRemoteCmd)
}

// 远程仓库管理命令
var remoteCmd = &cobra.Command{
	Use:   "remote",
	Short: "管理远程仓库",
	Long:  "添加、删除、查看远程仓库",
}

var addRemoteCmd = &cobra.Command{
	Use:   "add [远程名] [URL]",
	Short: "添加远程仓库",
	Long:  "添加一个新的远程仓库",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteName := args[0]
		remoteURL := args[1]

		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		if err := repo.AddRemote(remoteName, remoteURL); err != nil {
			return fmt.Errorf("添加远程仓库失败: %v", err)
		}

		fmt.Printf("已添加远程仓库: %s -> %s\n", remoteName, remoteURL)
		return nil
	},
}

var listRemoteCmd = &cobra.Command{
	Use:   "list",
	Short: "列出所有远程仓库",
	Long:  "显示所有配置的远程仓库",
	RunE: func(cmd *cobra.Command, args []string) error {
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		remotes, err := repo.ListRemotes()
		if err != nil {
			return fmt.Errorf("获取远程仓库列表失败: %v", err)
		}

		if len(remotes) == 0 {
			fmt.Println("暂无远程仓库")
			return nil
		}

		fmt.Println("远程仓库:")
		for _, remote := range remotes {
			fmt.Printf("  %s\t%s\n", remote.Name, remote.URL)
		}

		return nil
	},
}

var removeRemoteCmd = &cobra.Command{
	Use:   "remove [远程名]",
	Short: "删除远程仓库",
	Long:  "删除指定的远程仓库",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		remoteName := args[0]

		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		if err := repo.RemoveRemote(remoteName); err != nil {
			return fmt.Errorf("删除远程仓库失败: %v", err)
		}

		fmt.Printf("已删除远程仓库: %s\n", remoteName)
		return nil
	},
}

func init() {
	remoteCmd.AddCommand(addRemoteCmd)
	remoteCmd.AddCommand(listRemoteCmd)
	remoteCmd.AddCommand(removeRemoteCmd)
}
