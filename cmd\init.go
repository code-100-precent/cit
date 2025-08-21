package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"cit/internal/git"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "初始化一个新的Git仓库",
	Long:  "在当前目录或指定目录初始化一个新的Git仓库",
	Args:  cobra.MaximumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		dir := "."
		if len(args) > 0 {
			dir = args[0]
		}

		// 获取绝对路径
		absPath, err := filepath.Abs(dir)
		if err != nil {
			return fmt.Errorf("获取路径失败: %v", err)
		}

		// 检查目录是否存在
		if _, err := os.Stat(absPath); os.IsNotExist(err) {
			return fmt.Errorf("目录不存在: %s", absPath)
		}

		// 初始化仓库
		repo, err := git.InitRepository(absPath)
		if err != nil {
			return fmt.Errorf("初始化仓库失败: %v", err)
		}

		fmt.Printf("已在 %s 初始化空的Git仓库\n", absPath)
		fmt.Printf("仓库ID: %s\n", repo.ID)
		return nil
	},
}
