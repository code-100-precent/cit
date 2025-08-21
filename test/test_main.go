package test

import (
	"fmt"
	"os"
	"strings"

	"cit/internal/git"
)

func RunTest() {
	fmt.Println("CIT - Go语言实现的Git模拟器")
	fmt.Println(strings.Repeat("=", 40))

	// 测试初始化仓库
	fmt.Println("\n1. 测试初始化仓库...")
	repo, err := git.InitRepository(".")
	if err != nil {
		fmt.Printf("初始化仓库失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("仓库初始化成功！ID: %s\n", repo.ID)

	// 测试创建测试文件
	fmt.Println("\n2. 创建测试文件...")
	testContent := "这是一个测试文件\n用于演示CIT的功能"
	if err := os.WriteFile("test.txt", []byte(testContent), 0644); err != nil {
		fmt.Printf("创建测试文件失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("测试文件创建成功")

	// 测试添加文件到暂存区
	fmt.Println("\n3. 添加文件到暂存区...")
	if err := repo.AddToStaging("test.txt"); err != nil {
		fmt.Printf("添加文件到暂存区失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("文件已添加到暂存区")

	// 测试提交
	fmt.Println("\n4. 提交更改...")
	commit, err := repo.Commit("初始提交：添加测试文件")
	if err != nil {
		fmt.Printf("提交失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("提交成功！ID: %s\n", commit.ID)

	// 测试查看状态
	fmt.Println("\n5. 查看仓库状态...")
	status, err := repo.GetStatus()
	if err != nil {
		fmt.Printf("获取状态失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Printf("当前分支: %s\n", status.CurrentBranch)
	fmt.Printf("最新提交: %s\n", status.LastCommit)

	// 测试查看提交历史
	fmt.Println("\n6. 查看提交历史...")
	commits, err := repo.GetCommitHistory()
	if err != nil {
		fmt.Printf("获取提交历史失败: %v", err)
		os.Exit(1)
	}
	fmt.Printf("共有 %d 个提交:\n", len(commits))
	for i, c := range commits {
		fmt.Printf("  %d. %s - %s\n", i+1, c.ID[:8], c.Message)
	}

	// 测试分支操作
	fmt.Println("\n7. 测试分支操作...")
	if err := repo.CreateBranch("feature"); err != nil {
		fmt.Printf("创建分支失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("分支 'feature' 创建成功")

	branches, err := repo.ListBranches()
	if err != nil {
		fmt.Printf("列出分支失败: %v\n", err)
		os.Exit(1)
	}
	fmt.Println("所有分支:")
	for _, branch := range branches {
		if branch.Name == repo.GetCurrentBranch() {
			fmt.Printf("  * %s\n", branch.Name)
		} else {
			fmt.Printf("    %s\n", branch.Name)
		}
	}

	fmt.Println("\n测试完成！CIT Git模拟器工作正常。")
}
