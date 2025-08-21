package cmd

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "cit",
	Short: "CIT - Go语言实现的Git模拟器",
	Long: `CIT是一个用Go语言编写的Git模拟器，主要用于教学目的。
它实现了Git的核心功能，包括：
- 仓库初始化
- 文件暂存和提交
- 分支管理
- 提交历史查看
- 基本的合并功能`,
}

func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// 添加子命令
	rootCmd.AddCommand(initCmd)
	rootCmd.AddCommand(addCmd)
	rootCmd.AddCommand(commitCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(logCmd)
	rootCmd.AddCommand(branchCmd)
	rootCmd.AddCommand(checkoutCmd)
	rootCmd.AddCommand(pushCmd)
	rootCmd.AddCommand(remoteCmd)
}
