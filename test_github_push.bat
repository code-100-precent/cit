@echo off
echo CIT Git模拟器 - GitHub推送测试
echo ================================
echo.

echo 重要提示：
echo 请确保你已经：
echo 1. 在GitHub上创建了仓库: code-100-precent/cit
echo 2. 生成了个人访问令牌 (Personal Access Token)
echo 3. 令牌有repo权限
echo.

echo 请输入你的GitHub个人访问令牌:
set /p GITHUB_TOKEN="令牌: "

if "%GITHUB_TOKEN%"=="" (
    echo ❌ 错误: 令牌不能为空
    pause
    exit /b 1
)

echo.
echo 开始测试GitHub推送...
echo.

echo 1. 测试GitHub连接...
cit.exe push --github-token %GITHUB_TOKEN% origin main

if errorlevel 1 (
    echo.
    echo ❌ GitHub推送失败！
    echo.
    echo 可能的原因：
    echo - 令牌无效或过期
    echo - 仓库不存在或无权限
    echo - 网络连接问题
    echo - 防火墙阻止
    echo.
    echo 调试步骤：
    echo 1. 检查令牌是否正确
    echo 2. 确认仓库 https://github.com/code-100-precent/cit 存在
    echo 3. 验证令牌有repo权限
    echo 4. 检查网络连接
    echo.
) else (
    echo.
    echo ✅ GitHub推送成功！
    echo.
    echo 请检查你的GitHub仓库：
    echo https://github.com/code-100-precent/cit
    echo.
    echo 如果仓库中没有看到文件，可能是：
    echo - 分支名称不匹配
    echo - 文件路径问题
    echo - GitHub API延迟
    echo.
)

pause
