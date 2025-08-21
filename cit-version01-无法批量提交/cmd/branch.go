package cmd

import (
	"fmt"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var branchCmd = &cobra.Command{
	Use:   "branch [分支名]",
	Short: "管理分支",
	Long:  "创建、列出或删除分支",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		if len(args) == 0 {
			// 列出所有分支
			branches, err := repo.ListBranches()
			if err != nil {
				return fmt.Errorf("获取分支列表失败: %v", err)
			}

			fmt.Println("分支列表:")
			for _, branch := range branches {
				if branch.Name == repo.GetCurrentBranch() {
					fmt.Printf("* %s\n", branch.Name)
				} else {
					fmt.Printf("  %s\n", branch.Name)
				}
			}
			return nil
		}

		// 创建新分支
		branchName := args[0]
		if err := repo.CreateBranch(branchName); err != nil {
			return fmt.Errorf("创建分支失败: %v", err)
		}

		fmt.Printf("已创建分支: %s\n", branchName)
		return nil
	},
}
