package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "提交暂存区的更改",
	Long:  "将暂存区的更改提交到仓库，创建一个新的提交记录",
	RunE: func(cmd *cobra.Command, args []string) error {
		message, _ := cmd.Flags().GetString("message")
		if message == "" {
			return fmt.Errorf("必须提供提交信息，使用 -m 标志")
		}

		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 检查是否使用了 -a 标志
		addAll, _ := cmd.Flags().GetBool("all")
		if addAll {
			// 自动添加所有已跟踪的修改文件
			if err := autoAddModifiedFiles(repo); err != nil {
				return fmt.Errorf("自动添加文件失败: %v", err)
			}
		}

		// 检查暂存区是否有内容
		if repo.IsStagingEmpty() {
			return fmt.Errorf("暂存区为空，没有可提交的更改")
		}

		// 创建提交
		commit, err := repo.Commit(message)
		if err != nil {
			return fmt.Errorf("提交失败: %v", err)
		}

		fmt.Printf("提交成功！\n")
		fmt.Printf("提交ID: %s\n", commit.ID)
		fmt.Printf("提交信息: %s\n", commit.Message)
		fmt.Printf("作者: %s\n", commit.Author)
		fmt.Printf("时间: %s\n", commit.Timestamp)

		return nil
	},
}

func init() {
	commitCmd.Flags().StringP("message", "m", "", "提交信息")
	commitCmd.Flags().BoolP("all", "a", false, "自动添加所有已跟踪的修改文件")
	commitCmd.MarkFlagRequired("message")
}

// autoAddModifiedFiles 自动添加所有已跟踪的修改文件
func autoAddModifiedFiles(repo *git.Repository) error {
	// 简化实现：添加当前目录下所有非隐藏文件
	// 在实际的Git中，这会只添加已跟踪的修改文件
	return filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录和隐藏文件/目录
		if info.IsDir() {
			if strings.HasPrefix(info.Name(), ".") {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过隐藏文件
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// 添加文件到暂存区
		absPath, err := filepath.Abs(path)
		if err != nil {
			return err
		}

		return repo.AddToStaging(absPath)
	})
}
