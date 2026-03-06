package main_test

import (
	"os"
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
