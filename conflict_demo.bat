@echo off
echo CIT Git模拟器 - 冲突解决演示
echo ================================
echo.

echo 这个脚本将演示CIT的冲突解决功能
echo.

echo 1. 检查当前目录...
if not exist .cit (
    echo ❌ 错误: 当前目录不是CIT仓库
    echo 请先运行: cit init
    pause
    exit /b 1
)
echo ✅ 当前目录是CIT仓库
echo.

echo 2. 创建演示分支...
cit.exe branch feature-branch
echo.

echo 3. 切换到功能分支...
cit.exe checkout feature-branch
echo.

echo 4. 在功能分支上修改文件...
echo # 功能分支的修改 > feature.txt
echo 这是功能分支的修改内容 >> feature.txt
echo 添加新功能 >> feature.txt
echo.

echo 5. 提交功能分支的更改...
cit.exe add feature.txt
cit.exe commit -m "feat: 在功能分支上添加新功能"
echo.

echo 6. 切换回主分支...
cit.exe checkout main
echo.

echo 7. 在主分支上修改同一个文件...
echo # 主分支的修改 > feature.txt
echo 这是主分支的修改内容 >> feature.txt
echo 修复bug >> feature.txt
echo.

echo 8. 提交主分支的更改...
cit.exe add feature.txt
cit.exe commit -m "fix: 在主分支上修复bug"
echo.

echo 9. 尝试合并功能分支...
echo 这将产生冲突...
cit.exe merge feature-branch
echo.

echo 10. 查看冲突状态...
cit.exe conflicts
echo.

echo 11. 解决冲突 (保留我们的内容)...
cit.exe resolve feature.txt ours
echo.

echo 12. 再次查看冲突状态...
cit.exe conflicts
echo.

echo 13. 提交解决后的更改...
cit.exe add feature.txt
cit.exe commit -m "merge: 解决冲突，保留主分支内容"
echo.

echo 14. 查看合并后的状态...
cit.exe status
echo.

echo 15. 查看提交历史...
cit.exe log
echo.

echo ✅ 冲突解决演示完成！
echo.
echo 演示内容:
echo - 创建分支并修改相同文件
echo - 检测合并冲突
echo - 生成冲突标记
echo - 使用策略解决冲突
echo - 完成合并提交
echo.

pause
