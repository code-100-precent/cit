package cmd

import (
	"fmt"
	"strings"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "显示提交历史",
	Long:  "显示Git仓库的提交历史记录",
	RunE: func(cmd *cobra.Command, args []string) error {
		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 获取提交历史
		commits, err := repo.GetCommitHistory()
		if err != nil {
			return fmt.Errorf("获取提交历史失败: %v", err)
		}

		if len(commits) == 0 {
			fmt.Println("暂无提交记录")
			return nil
		}

		// 显示提交历史
		fmt.Println("提交历史:")
		fmt.Println(strings.Repeat("=", 80))
		
		for i, commit := range commits {
			fmt.Printf("提交 #%d\n", len(commits)-i)
			fmt.Printf("ID: %s\n", commit.ID)
			fmt.Printf("信息: %s\n", commit.Message)
			fmt.Printf("作者: %s\n", commit.Author)
			fmt.Printf("时间: %s\n", commit.Timestamp)
			if commit.ParentID != "" {
				fmt.Printf("父提交: %s\n", commit.ParentID)
			}
			if i < len(commits)-1 {
				fmt.Println(strings.Repeat("-", 40))
			}
		}

		return nil
	},
}
