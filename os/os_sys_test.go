package main_test

import (
	"errors"
	"os"
	"os/exec"
	"runtime"
	"testing"
)

func TestRootSandboxing(t *testing.T) {
	baseDir := "secure_sandbox"
	os.MkdirAll(baseDir+"/public", 0755)
	defer os.RemoveAll(baseDir)

	// 1. 测试 OpenRoot
	root, err := os.OpenRoot(baseDir)
	if err != nil {
		t.Fatalf("OpenRoot 初始化沙箱失败: %v", err)
	}
	defer root.Close() // 测试 Close 方法

	// 2. 测试隔离下的 WriteFile 和 ReadFile (Go 1.25)
	err = root.WriteFile("public/data.json", []byte("{\"status\":\"ok\"}"), 0644)
	if err != nil {
		t.Fatalf("root.WriteFile 执行失败: %v", err)
	}
	content, _ := root.ReadFile("public/data.json")
	if string(content) != "{\"status\":\"ok\"}" {
		t.Errorf("隔离内的数据不匹配")
	}

	// 3. 测试隔离下的重命名 (Go 1.25)
	err = root.Rename("public/data.json", "public/moved.json")
	if err != nil {
		t.Fatalf("root.Rename 执行失败: %v", err)
	}

	// 4. 核心安全验证：测试对抗路径遍历攻击
	_, err = root.ReadFile("../../../etc/shadow")
	if err == nil {
		t.Fatalf("严重漏洞：os.Root 无法阻挡基于 '../' 的相对路径溢出攻击")
	}

	// 5. 测试 FS 适配器集成
	vfs := root.FS()
	fileInVfs, err := vfs.Open("public/moved.json")
	if err != nil {
		t.Errorf("root.FS 适配器无法打开内部验证文件: %v", err)
	}
	fileInVfs.Close()
}
func TestProcessIdentityAndTracing(t *testing.T) {
	// 1. 测试 PID 和 PPID
	pid := os.Getpid()
	if pid <= 0 {
		t.Errorf("Getpid 获取的进程号无效: %d", pid)
	}
	ppid := os.Getppid()
	if ppid < 0 {
		t.Errorf("Getppid 获取的父进程号异常: %d", ppid)
	}

	// 2. 测试 UID 与 GID 系统
	uid := os.Getuid()
	euid := os.Geteuid()
	if uid < 0 || euid < 0 {
		t.Errorf("用户 ID 必须为非负整数")
	}

	// 3. 测试附加组列表
	groups, err := os.Getgroups()
	if err == nil && len(groups) == 0 {
		t.Log("当前用户没有额外归属组")
	}

	// 4. 测试 Executable
	execPath, err := os.Executable()
	if err != nil || execPath == "" {
		t.Errorf("Executable 无法反解析程序运行实体位置: %v", err)
	}

	// os.Exit 不能被放置在标准的 Test 函数中直接执行，
	// 因为这会强行结束整个 go test 的测试驱动容器环境。
}
func TestProcessExecutionAndHandle(t *testing.T) {
	// 构造跨平台的执行环境
	cmdName := "sleep"
	args := []string{"sleep", "1"}

	if runtime.GOOS == "windows" {
		cmdName = "timeout"
		args = []string{"timeout", "1"}
		args = []string{"timeout", "/T", "1", "/NOBREAK"}
	}

	// 必须获取绝对路径
	execPath, err := exec.LookPath(cmdName)
	if err != nil {
		t.Skip("跳过进程测试，因为宿主机缺少基础的阻塞命令工具")
	}

	// 1. 测试 StartProcess
	procAttr := &os.ProcAttr{
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}

	proc, err := os.StartProcess(execPath, args, procAttr)
	if err != nil {
		t.Fatalf("StartProcess 派生新进程失败: %v", err)
	}

	// 2. 测试 Go 1.26 新特性 WithHandle
	err = proc.WithHandle(func(handle uintptr) {
		if handle == 0 {
			t.Errorf("WithHandle 注入的底层 PIDFD 句柄不能为 0")
		}
		// 在此处，开发者可将句柄传入 CGO 或通过 syscall 包直接触发底层 ioctl
	})

	// 仅在不支持 pidfd 的老旧 Linux 或非 Windows 内核上允许 ErrNoHandle
	if err != nil && !errors.Is(err, os.ErrNoHandle) {
		t.Errorf("WithHandle 执行发生未预期异常: %v", err)
	}

	// 3. 测试 Kill 发出终止信号
	if err := proc.Kill(); err != nil {
		t.Fatalf("Kill 强制终止进程失败: %v", err)
	}

	// 4. 测试 Wait 与 ProcessState 遥测数据析构
	state, err := proc.Wait()
	if err != nil {
		t.Fatalf("Wait 等待进程结束挂起失败: %v", err)
	}

	if state.Success() {
		t.Errorf("由于手动发起了 Kill，进程的退出状态不应该被判定为 Success")
	}

	if state.SystemTime() < 0 || state.UserTime() < 0 {
		t.Errorf("进程 CPU 时间统计异常")
	}

	// 测试 Release (由于已经 Wait，此处 Release 主要为清理内存对象结构)
	proc.Release()
}

func TestSystemPropertiesAndDirectories(t *testing.T) {
	// 1. 测试 Hostname
	host, err := os.Hostname()
	if err != nil || host == "" {
		t.Errorf("Hostname 查询失败或值为空")
	}

	// 2. 测试 Getpagesize
	pageSize := os.Getpagesize()
	if pageSize <= 0 {
		t.Errorf("Getpagesize 获取硬件内存页表大小失败: %d", pageSize)
	}

	// 3. 测试规范用户目录系列
	home, _ := os.UserHomeDir()
	config, _ := os.UserConfigDir()
	cache, _ := os.UserCacheDir()

	if home == "" || config == "" || cache == "" {
		t.Errorf("操作系统的环境规范获取失败，某些目录未正确解析")
	}
}
func TestIPCAndVirtualFilesystem(t *testing.T) {
	// 1. 测试 Pipe 进行并发 IPC
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Pipe 创建匿名通道失败: %v", err)
	}

	go func() {
		// 在新协程中写入管道，利用通道传递字符串
		w.Write([]byte("Kernel Ring Buffer Message"))
		w.Close() // 只有关闭写入端，读取端才会收到 EOF，打破阻塞
	}()

	buffer := make([]byte, 64)
	n, _ := r.Read(buffer)
	r.Close()

	if string(buffer[:n]) != "Kernel Ring Buffer Message" {
		t.Errorf("IPC 管道传输数据损坏或丢失，收到： %s", string(buffer[:n]))
	}

	// 2. 测试 DirFS 与 CopyFS 结构转移
	vfsSource := "source_vfs"
	vfsTarget := "target_vfs"
	os.Mkdir(vfsSource, 0755)
	os.WriteFile(vfsSource+"/node.txt", []byte("vfs payload"), 0644)
	defer func() {
		os.RemoveAll(vfsSource)
		os.RemoveAll(vfsTarget)
	}()

	// 将物理层包裹成 FS 接口
	virtualSystem := os.DirFS(vfsSource)

	// 将 FS 接口映射并持久化拷贝到另一侧磁盘
	os.Mkdir(vfsTarget, 0755)
	if err := os.CopyFS(vfsTarget, virtualSystem); err != nil {
		t.Fatalf("CopyFS 转储虚拟文件系统失败: %v", err)
	}

	verifyData, _ := os.ReadFile(vfsTarget + "/node.txt")
	if string(verifyData) != "vfs payload" {
		t.Errorf("CopyFS 数据复制存在丢失")
	}
}
