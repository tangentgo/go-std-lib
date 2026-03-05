package main_test

import (
	"os"
	"testing"
)

func TestWorkingDirectoryContext(t *testing.T) {
	// 1. 测试 Getwd 记录原始目录
	originalDir, err := os.Getwd()
	if err != nil {
		t.Fatalf("Getwd 获取初始工作目录失败: %v", err)
	}

	// 2. 测试 TempDir
	sysTemp := os.TempDir()
	if sysTemp == "" {
		t.Errorf("TempDir 返回值为空，系统环境异常")
	}

	// 3. 测试 Chdir
	if err := os.Chdir(sysTemp); err != nil {
		t.Fatalf("Chdir 切换工作目录失败: %v", err)
	}

	// 清理：恢复原始目录，防止污染测试运行器上下文
	defer func() {
		if err := os.Chdir(originalDir); err != nil {
			t.Fatalf("恢复工作目录失败: %v", err)
		}
	}()

	currentDir, _ := os.Getwd()
	if currentDir == originalDir {
		t.Errorf("Chdir 之后，工作目录状态未发生实际改变")
	}
}
