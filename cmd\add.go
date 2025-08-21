package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add [文件或目录]",
	Short: "将文件添加到暂存区",
	Long:  "将指定的文件或目录添加到Git暂存区，准备下一次提交",
	Args:  cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		// 查找Git仓库
		repo, err := git.FindRepository(".")
		if err != nil {
			return fmt.Errorf("未找到Git仓库: %v", err)
		}

		// 添加文件到暂存区
		for _, path := range args {
			if err := addPath(repo, path); err != nil {
				fmt.Printf("警告: 添加 %s 失败: %v\n", path, err)
			}
		}

		return nil
	},
}

// addPath 添加路径（文件或目录）到暂存区
func addPath(repo *git.Repository, path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		return fmt.Errorf("路径错误: %v", err)
	}

	// 检查路径是否存在
	info, err := os.Stat(absPath)
	if os.IsNotExist(err) {
		return fmt.Errorf("路径不存在")
	}
	if err != nil {
		return fmt.Errorf("访问路径失败: %v", err)
	}

	if info.IsDir() {
		// 处理目录 - 递归添加所有文件
		return addDirectory(repo, absPath, path)
	} else {
		// 处理单个文件
		return addFile(repo, absPath, path)
	}
}

// addFile 添加单个文件到暂存区
func addFile(repo *git.Repository, absPath, originalPath string) error {
	if err := repo.AddToStaging(absPath); err != nil {
		return err
	}
	fmt.Printf("已添加 %s 到暂存区\n", originalPath)
	return nil
}

// addDirectory 递归添加目录中的所有文件
func addDirectory(repo *git.Repository, absPath, originalPath string) error {
	return filepath.Walk(absPath, func(filePath string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// 跳过目录本身和.cit目录
		if info.IsDir() {
			if strings.HasSuffix(filePath, ".cit") || strings.Contains(filePath, ".cit"+string(filepath.Separator)) {
				return filepath.SkipDir
			}
			return nil
		}

		// 跳过隐藏文件和系统文件
		if strings.HasPrefix(info.Name(), ".") {
			return nil
		}

		// 添加文件
		relPath, err := filepath.Rel(absPath, filePath)
		if err != nil {
			relPath = info.Name()
		}

		if originalPath == "." {
			fmt.Printf("已添加 %s 到暂存区\n", relPath)
		} else {
			fmt.Printf("已添加 %s 到暂存区\n", filepath.Join(originalPath, relPath))
		}

		return repo.AddToStaging(filePath)
	})
}
