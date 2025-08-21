package cmd

import (
	"fmt"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var checkoutCmd = &cobra.Command{
	Use:   "checkout <分支名>",
	Short: "切换到指定分支",
	Long:  "切换到指定的分支，更新工作目录",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		branchName := args[0]

		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 检查分支是否存在
		branches, err := repo.ListBranches()
		if err != nil {
			return fmt.Errorf("获取分支列表失败: %v", err)
		}

		branchExists := false
		for _, branch := range branches {
			if branch.Name == branchName {
				branchExists = true
				break
			}
		}

		if !branchExists {
			return fmt.Errorf("分支 '%s' 不存在", branchName)
		}

		// 切换到指定分支
		if err := repo.CheckoutBranch(branchName); err != nil {
			return fmt.Errorf("切换分支失败: %v", err)
		}

		fmt.Printf("已切换到分支: %s\n", branchName)
		return nil
	},
}
