package main_test

import (
	"os"
	"strings"
	"testing"
	"time"
)

func TestDirectoryLifecycle(t *testing.T) {
	baseDir := "test_lifecycle_dir"

	// 1. 测试 Mkdir
	err := os.Mkdir(baseDir, 0750)
	if err != nil && !os.IsExist(err) {
		t.Fatalf("Mkdir 执行失败: %v", err)
	}

	// 2. 测试 MkdirAll
	nestedPath := baseDir + "/level1/level2"
	if err := os.MkdirAll(nestedPath, 0750); err != nil {
		t.Fatalf("MkdirAll 递归创建失败: %v", err)
	}

	// 3. 测试 MkdirTemp
	tempDir, err := os.MkdirTemp(baseDir, "worker-*")
	if err != nil {
		t.Fatalf("MkdirTemp 创建安全临时目录失败: %v", err)
	}

	// 4. 测试 Remove (仅删除空目录或文件)
	if err := os.Remove(tempDir); err != nil {
		t.Fatalf("Remove 失败: %v", err)
	}

	// 5. 测试 RemoveAll (递归清理)
	if err := os.RemoveAll(baseDir); err != nil {
		t.Fatalf("RemoveAll 清理资源树失败: %v", err)
	}
}
func TestRename(t *testing.T) {
	oldName := "rename_source.txt"
	newName := "rename_target.txt"

	os.WriteFile(oldName, []byte("data"), 0644)
	// defer os.Remove(newName)

	if err := os.Rename(oldName, newName); err != nil {
		t.Fatalf("Rename 重命名失败: %v", err)
	}

	if _, err := os.Stat(oldName); !os.IsNotExist(err) {
		t.Errorf("Rename 后，源文件不应继续存在")
	}
}

func TestFileMetadataManipulation(t *testing.T) {
	fileName := "meta_test.txt"
	os.WriteFile(fileName, []byte("metadata"), 0644)
	defer os.Remove(fileName)

	// 1. 测试 Chmod
	if err := os.Chmod(fileName, 0755); err != nil {
		t.Fatalf("Chmod 修改权限失败: %v", err)
	}

	// 2. 测试 Chtimes
	newTime := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	if err := os.Chtimes(fileName, newTime, newTime); err != nil {
		t.Fatalf("Chtimes 更新时间戳失败: %v", err)
	}

	// 3. 测试 Chown / Lchown
	// 由于 Chown 在普通非特权用户下执行通常会返回 operation not permitted，
	// 此处的测试主要验证函数的调用链路是否正常，并妥善处理预期的权限错误。
	err := os.Chown(fileName, os.Getuid(), os.Getgid())
	if err != nil && !os.IsPermission(err) {
		t.Errorf("Chown 抛出了非预期的错误: %v", err)
	}
}

func TestLinkingMechanisms(t *testing.T) {
	original := "core_data.txt"
	hLink := "hard_link.txt"
	sLink := "soft_link.txt"

	os.WriteFile(original, []byte("critical system data"), 0644)
	defer func() {
		os.Remove(original)
		os.Remove(hLink)
		os.Remove(sLink)
	}()

	// 1. 测试 Link (硬链接)
	if err := os.Link(original, hLink); err != nil {
		t.Fatalf("Link 创建硬链接失败: %v", err)
	}

	// 2. 测试 Symlink (符号链接)
	if err := os.Symlink(original, sLink); err != nil {
		t.Fatalf("Symlink 创建符号链接失败: %v", err)
	}

	// 3. 测试 Readlink
	targetPath, err := os.Readlink(sLink)
	if err != nil || targetPath != original {
		t.Errorf("Readlink 解析符号链接失败，期望 %s，实际 %s", original, targetPath)
	}

	// 4. 测试 SameFile 判断物理同一性
	fiOrig, _ := os.Stat(original)
	fiHard, _ := os.Stat(hLink)
	fiSoft, _ := os.Lstat(sLink) // Lstat 获取符号链接自身的信息

	if !os.SameFile(fiOrig, fiHard) {
		t.Errorf("SameFile 失败：硬链接的源与目标应指向同一 Inode")
	}
	if os.SameFile(fiOrig, fiSoft) {
		t.Errorf("SameFile 失败：符号链接自身是一个独立的文件，不应与源文件判定为同一物理文件")
	}
}
func TestTruncateFunction(t *testing.T) {
	fileName := "truncate_test.dat"
	os.WriteFile(fileName, []byte("1234567890"), 0644)
	defer os.Remove(fileName)

	// 测试 Truncate 缩减文件
	if err := os.Truncate(fileName, 5); err != nil {
		t.Fatalf("Truncate 截断文件失败: %v", err)
	}

	data, _ := os.ReadFile(fileName)
	if string(data) != "12345" {
		t.Errorf("Truncate 失败：期望的数据为 '12345'，实际为 %s", data)
	}
}

func TestFileInstantiation(t *testing.T) {
	// 1. 测试 Create
	fCreate, err := os.Create("instantiation_test.txt")
	if err != nil {
		t.Fatalf("Create 失败: %v", err)
	}
	fCreate.Close()
	defer os.Remove("instantiation_test.txt")

	// 2. 测试 OpenFile (追加写入模式)
	fOpenFile, err := os.OpenFile("instantiation_test.txt", os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		t.Fatalf("OpenFile 追加模式打开失败: %v", err)
	}
	fOpenFile.Close()

	// 3. 测试 Open (只读模式)
	fOpen, err := os.Open("instantiation_test.txt")
	if err != nil {
		t.Fatalf("Open 失败: %v", err)
	}
	fOpen.Close()

	// 4. 测试 CreateTemp
	fTemp, err := os.CreateTemp("", "secure-temp-*.log")
	if err != nil {
		t.Fatalf("CreateTemp 失败: %v", err)
	}
	defer os.Remove(fTemp.Name())
	fTemp.Close()
}
func TestHighLevelIO(t *testing.T) {
	fileName := "high_level.txt"
	payload := []byte("Automated descriptor management.")
	defer os.Remove(fileName)

	// 1. 测试 WriteFile
	if err := os.WriteFile(fileName, payload, 0644); err != nil {
		t.Fatalf("WriteFile 写入失败: %v", err)
	}

	// 2. 测试 ReadFile
	readBack, err := os.ReadFile(fileName)
	if err != nil {
		t.Fatalf("ReadFile 读取失败: %v", err)
	}

	if string(readBack) != string(payload) {
		t.Errorf("数据完整性受损，期望 %s，实际 %s", payload, readBack)
	}
}
func TestFileStreamIO(t *testing.T) {
	f, _ := os.Create("stream.bin")
	defer os.Remove("stream.bin")
	defer f.Close()

	// 1. 测试 Write / WriteString
	f.WriteString("HEAD")
	f.Write([]byte("-BODY"))

	// 2. 测试 WriteAt (在指定位置进行原子覆写，不改变游标)
	f.WriteAt([]byte("TAIL"), 9) // 此时文件内容：HEAD-BODYTAIL

	// 3. 测试 ReadAt
	buffer := make([]byte, 4)
	f.ReadAt(buffer, 5) // 从偏移量 5 读取 4 字节
	if string(buffer) != "BODY" {
		t.Errorf("ReadAt 未命中目标数据，实际读取: %s", buffer)
	}

	// 4. 测试 ReadFrom
	source := strings.NewReader("-FOOTER")
	f.Seek(0, 2) // 将游标移动到末尾 (SEEK_END)
	f.ReadFrom(source)

	// 验证全量状态
	finalData, _ := os.ReadFile("stream.bin")
	if string(finalData) != "HEAD-BODYTAIL-FOOTER" {
		t.Errorf("数据流合成错误，结果为: %s", finalData)
	}
}
