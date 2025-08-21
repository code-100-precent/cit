@echo off
echo CIT Git模拟器 - GitHub推送诊断脚本
echo ======================================
echo.

echo 这个脚本将帮助你诊断GitHub推送问题
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

echo 2. 检查远程仓库配置...
cit.exe remote list
echo.

echo 3. 检查暂存区状态...
cit.exe status
echo.

echo 4. 检查提交历史...
cit.exe log
echo.

echo 5. 重要提示：
echo 请确保你已经：
echo - 在GitHub上创建了仓库: code-100-precent/cit
echo - 生成了个人访问令牌 (Personal Access Token)
echo - 令牌有repo权限
echo.

echo 6. 测试GitHub推送（请替换YOUR_TOKEN）：
echo cit.exe push --github-token YOUR_TOKEN origin main
echo.

echo 7. 如果推送失败，请检查：
echo - 令牌是否正确
echo - 仓库是否存在
echo - 网络连接是否正常
echo - 防火墙设置
echo.

echo 8. 调试命令：
echo - 查看详细错误信息
echo - 检查GitHub API响应
echo - 验证文件内容
echo.

pause
