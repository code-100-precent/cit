package cmd

import (
	"fmt"
	"strings"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var mergeCmd = &cobra.Command{
	Use:   "merge [分支名]",
	Short: "合并指定分支到当前分支",
	Long:  "将指定分支的更改合并到当前分支，自动处理冲突",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		sourceBranch := args[0]

		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 检查是否尝试合并到自己
		if sourceBranch == repo.GetCurrentBranch() {
			return fmt.Errorf("不能合并分支到自身")
		}

		fmt.Printf("正在合并分支 '%s' 到 '%s'...\n", sourceBranch, repo.GetCurrentBranch())

		// 执行合并
		result, err := repo.MergeBranch(sourceBranch)
		if err != nil {
			return fmt.Errorf("合并失败: %v", err)
		}

		// 显示合并结果
		if result.Success {
			fmt.Printf("✅ %s\n", result.Message)
		} else {
			fmt.Printf("❌ %s\n", result.Message)
			fmt.Printf("发现 %d 个冲突文件:\n", len(result.Conflicts))

			for _, conflict := range result.Conflicts {
				fmt.Printf("  - %s\n", conflict)
			}

			fmt.Println("\n请使用以下命令解决冲突:")
			fmt.Printf("  cit-version01-无法批量提交 resolve <文件路径> <策略>\n")
			fmt.Printf("  策略选项: ours(保留我们的), theirs(保留他们的), both(保留两者)\n")
			fmt.Printf("  例如: cit-version01-无法批量提交 resolve %s ours\n", result.Conflicts[0])
		}

		return nil
	},
}

var resolveCmd = &cobra.Command{
	Use:   "resolve [文件路径] [策略]",
	Short: "解决文件冲突",
	Long:  "使用指定策略解决文件冲突",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		filePath := args[0]
		resolution := args[1]

		// 验证解决策略
		validStrategies := []string{"ours", "theirs", "both"}
		isValid := false
		for _, strategy := range validStrategies {
			if resolution == strategy {
				isValid = true
				break
			}
		}
		if !isValid {
			return fmt.Errorf("无效的解决策略 '%s'，支持: %s", resolution, strings.Join(validStrategies, ", "))
		}

		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 解决冲突
		if err := repo.ResolveConflict(filePath, resolution); err != nil {
			return fmt.Errorf("解决冲突失败: %v", err)
		}

		fmt.Printf("✅ 成功解决文件 '%s' 的冲突 (策略: %s)\n", filePath, resolution)
		return nil
	},
}

var conflictsCmd = &cobra.Command{
	Use:   "conflicts",
	Short: "显示当前冲突状态",
	Long:  "列出所有有冲突的文件",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 获取冲突状态
		conflicts, err := repo.GetConflictStatus()
		if err != nil {
			return fmt.Errorf("获取冲突状态失败: %v", err)
		}

		if len(conflicts) == 0 {
			fmt.Println("✅ 当前没有冲突")
			return nil
		}

		fmt.Printf("发现 %d 个冲突文件:\n", len(conflicts))
		for _, conflict := range conflicts {
			fmt.Printf("  - %s\n", conflict)
		}

		fmt.Println("\n使用以下命令解决冲突:")
		fmt.Printf("  cit-version01-无法批量提交 resolve <文件路径> <策略>\n")
		fmt.Printf("  策略选项: ours(保留我们的), theirs(保留他们的), both(保留两者)\n")

		return nil
	},
}

func init() {
	// 这里不需要注册到根命令，因为已经在root.go中注册了
}
