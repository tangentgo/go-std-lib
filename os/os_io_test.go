package main_test

import (
	"errors"
	"os"
	"testing"
)

func TestGlobalVariables(t *testing.T) {
	// 测试 Args 变量
	if len(os.Args) < 1 {
		t.Errorf("操作系统必须至少传递一个参数（程序自身路径）给 os.Args")
	}

	// 测试标准输出流的文件描述符编号
	if os.Stdin.Fd() != 0 {
		t.Errorf("期望的 Stdin 描述符为 0，实际为 %d", os.Stdin.Fd())
	}
	if os.Stdout.Fd() != 1 {
		t.Errorf("期望的 Stdout 描述符为 1，实际为 %d", os.Stdout.Fd())
	}
	if os.Stderr.Fd() != 2 {
		t.Errorf("期望的 Stderr 描述符为 2，实际为 %d", os.Stderr.Fd())
	}

	// 测试预定义错误与 errors.Is 的兼容性
	_, err := os.Open("completely_random_file_that_does_not_exist.txt")
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("期望捕获到 os.ErrNotExist 错误，实际捕获到: %v", err)
	}
}
