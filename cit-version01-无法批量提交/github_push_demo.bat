@echo off
echo CIT Git模拟器 - GitHub真实推送演示
echo ======================================
echo.

echo 重要提示：
echo 1. 你需要有GitHub个人访问令牌 (Personal Access Token)
echo 2. 令牌需要有repo权限
echo 3. 目标仓库必须存在且你有推送权限
echo.

echo 请确保你已经：
echo - 在GitHub上创建了仓库: code-100-precent/cit
echo - 生成了个人访问令牌
echo - 令牌有足够的权限
echo.

pause

echo 1. 清理之前的测试数据...
if exist .cit rmdir /s /q .cit
if exist test.txt del test.txt
echo 清理完成！
echo.

echo 2. 构建项目...
go build -o cit.exe
if errorlevel 1 (
    echo 构建失败！
    pause
    exit /b 1
)
echo 构建成功！
echo.

echo 3. 初始化Git仓库...
cit.exe init
if errorlevel 1 (
    echo 初始化失败！
    pause
    exit /b 1
)
echo 仓库初始化成功！
echo.

echo 4. 创建测试文件...
echo # CIT Git模拟器 > test.txt
echo. >> test.txt
echo 这是一个用Go语言编写的Git模拟器，主要用于教学目的。 >> test.txt
echo 它实现了Git的核心功能，包括： >> test.txt
echo - 仓库初始化 >> test.txt
echo - 文件暂存和提交 >> test.txt
echo - 分支管理 >> test.txt
echo - 提交历史查看 >> test.txt
echo - 基本的合并功能 >> test.txt
echo. >> test.txt
echo 现在支持真正的GitHub推送！ >> test.txt
echo 测试文件创建完成！
echo.

echo 5. 添加文件到暂存区...
cit.exe add test.txt
echo 文件添加完成！
echo.

echo 6. 提交更改...
cit.exe commit -m "feat: 添加CIT项目介绍文件"
if errorlevel 1 (
    echo 提交失败！
    pause
    exit /b 1
)
echo 提交成功！
echo.

echo 7. 添加GitHub远程仓库...
cit.exe remote add origin https://github.com/code-100-precent/cit.git
if errorlevel 1 (
    echo 添加远程仓库失败！
    pause
    exit /b 1
)
echo 远程仓库添加成功！
echo.

echo 8. 查看远程仓库列表...
cit.exe remote list
echo.

echo 9. 查看提交历史...
cit.exe log
echo.

echo 10. 重要：现在需要GitHub令牌进行真实推送！
echo.
echo 请运行以下命令（替换YOUR_TOKEN为你的实际令牌）：
echo cit.exe push --github-token YOUR_TOKEN origin main
echo.
echo 或者，如果你想先测试模拟推送：
echo cit.exe push origin main
echo.

echo 演示准备完成！
echo.
echo 下一步：
echo 1. 获取GitHub个人访问令牌
echo 2. 运行真实推送命令
echo 3. 检查GitHub仓库是否更新
echo.

pause
