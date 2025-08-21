@echo off
echo CIT Git模拟器演示脚本
echo ========================
echo.

echo 1. 清理之前的测试数据...
if exist .cit rmdir /s /q .cit
if exist test.txt del test.txt
if exist hello.txt del hello.txt
if exist demo.txt del demo.txt
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
echo Hello, CIT! > hello.txt
echo 这是一个演示文件 > demo.txt
echo 测试文件创建完成！
echo.

echo 5. 查看仓库状态...
cit.exe status
echo.

echo 6. 添加文件到暂存区...
cit.exe add hello.txt
cit.exe add demo.txt
echo 文件添加完成！
echo.

echo 7. 查看暂存区状态...
cit.exe status
echo.

echo 8. 提交更改...
cit.exe commit -m "初始提交：添加演示文件"
if errorlevel 1 (
    echo 提交失败！
    pause
    exit /b 1
)
echo 提交成功！
echo.

echo 9. 查看提交历史...
cit.exe log
echo.

echo 10. 创建新分支...
cit.exe branch feature-demo
echo 分支创建完成！
echo.

echo 11. 查看所有分支...
cit.exe branch
echo.

echo 12. 切换到新分支...
cit.exe checkout feature-demo
echo 分支切换完成！
echo.

echo 13. 在新分支上工作...
echo 新功能代码 > feature.txt
cit.exe add feature.txt
cit.exe commit -m "在新分支上添加功能"
echo 新功能开发完成！
echo.

echo 14. 切换回主分支...
cit.exe checkout main
echo 回到主分支！
echo.

echo 15. 查看最终状态...
cit.exe status
echo.

echo 16. 查看所有分支的提交历史...
cit.exe log
echo.

echo 演示完成！CIT Git模拟器工作正常。
echo.
echo 生成的文件和目录：
dir .cit
echo.
echo 对象存储：
dir .cit\objects
echo.
pause
