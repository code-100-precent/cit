@echo off
echo CIT Git模拟器 - 远程功能演示
echo ================================
echo.

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
echo Hello, Remote World! > test.txt
echo 测试文件创建完成！
echo.

echo 5. 添加文件到暂存区...
cit.exe add test.txt
echo 文件添加完成！
echo.

echo 6. 提交更改...
cit.exe commit -m "初始提交：添加测试文件"
if errorlevel 1 (
    echo 提交失败！
    pause
    exit /b 1
)
echo 提交成功！
echo.

echo 7. 查看远程仓库列表...
cit.exe remote list
echo.

echo 8. 添加远程仓库...
cit.exe remote add origin https://github.com/user/repo.git
if errorlevel 1 (
    echo 添加远程仓库失败！
    pause
    exit /b 1
)
echo 远程仓库添加成功！
echo.

echo 9. 再次查看远程仓库列表...
cit.exe remote list
echo.

echo 10. 推送到远程仓库...
cit.exe push origin main
if errorlevel 1 (
    echo 推送失败！
    pause
    exit /b 1
)
echo 推送成功！
echo.

echo 11. 查看提交历史...
cit.exe log
echo.

echo 12. 查看仓库状态...
cit.exe status
echo.

echo 演示完成！CIT Git模拟器的远程功能工作正常。
echo.
echo 生成的文件和目录：
dir .cit
echo.
echo 远程仓库配置：
type .cit\remotes.json
echo.
pause
