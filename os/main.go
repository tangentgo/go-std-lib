package main

import (
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// 1. 创建或打开一个日志文件 (如果文件不存在则创建，如果存在则追加写入)
	logFile, err := os.OpenFile("process_output.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		fmt.Printf("打开日志文件失败: %v\n", err)
		return
	}
	// 记得在函数退出时关闭文件句柄
	defer logFile.Close()

	// 2. 准备执行的命令
	// 这里我们用 sh -c 执行两句话：一句正常打印，一句故意报错(打印到标准错误)
	cmdName := "sh"
	args := []string{"sh", "-c", "echo '这是一条正常输出'; ls /non_existent_folder"}

	execPath, err := exec.LookPath(cmdName)
	if err != nil {
		fmt.Printf("找不到命令: %v\n", err)
		return
	}

	// 3. 核心修改点：配置 ProcAttr 的 Files 数组
	// 索引 0 (Stdin)  -> nil      (切断输入，子进程就像在后台运行，不需要键盘)
	// 索引 1 (Stdout) -> logFile  (把标准输出接到日志文件)
	// 索引 2 (Stderr) -> logFile  (把标准错误也接到同一个日志文件)
	procAttr := &os.ProcAttr{
		Files: []*os.File{nil, logFile, logFile},
	}

	fmt.Println("正在悄悄启动子进程，请观察屏幕，你将不会看到子进程的输出...")

	// 4. 启动进程
	proc, err := os.StartProcess(execPath, args, procAttr)
	if err != nil {
		fmt.Printf("启动进程失败: %v\n", err)
		return
	}

	// 5. 等待进程结束
	state, err := proc.Wait()
	if err != nil {
		fmt.Printf("等待进程结束失败: %v\n", err)
		return
	}

	fmt.Printf("子进程执行完毕，退出状态: %v\n", state.Success())
	fmt.Println("现在你可以去查看当前目录下的 process_output.log 文件了！")
}
