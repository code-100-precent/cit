package cmd

import (
	"fmt"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "显示仓库状态",
	Long:  "显示工作目录和暂存区的状态",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 获取状态
		status, err := repo.GetStatus()
		if err != nil {
			return fmt.Errorf("获取状态失败: %v", err)
		}

		// 显示状态
		fmt.Println("仓库状态:")
		fmt.Printf("当前分支: %s\n", status.CurrentBranch)
		fmt.Printf("最新提交: %s\n", status.LastCommit)
		
		if len(status.StagedFiles) > 0 {
			fmt.Println("\n暂存区文件:")
			for _, file := range status.StagedFiles {
				fmt.Printf("  %s\n", file)
			}
		}

		if len(status.ModifiedFiles) > 0 {
			fmt.Println("\n已修改但未暂存的文件:")
			for _, file := range status.ModifiedFiles {
				fmt.Printf("  %s\n", file)
			}
		}

		if len(status.UntrackedFiles) > 0 {
			fmt.Println("\n未跟踪的文件:")
			for _, file := range status.UntrackedFiles {
				fmt.Printf("  %s\n", file)
			}
		}

		return nil
	},
}
