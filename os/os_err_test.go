package main_test

import (
	"errors"
	"os"
	"testing"
)

func TestErrorEvaluations(t *testing.T) {
	// 1. 测试 IsNotExist 与 modern errors.Is 的等效性
	_, err := os.Stat("an_impossible_file_path_12345.dat")
	if err == nil {
		t.Fatalf("对于不存在的文件，系统没有产生预期的失败")
	}
	if !os.IsNotExist(err) {
		t.Errorf("IsNotExist 断言失败：无法有效识别缺失资源")
	}
	if !errors.Is(err, os.ErrNotExist) {
		t.Errorf("errors.Is 无法向后兼容判断 ErrNotExist")
	}

	// 2. 测试 IsExist
	conflictFile := "conflict.bin"
	os.WriteFile(conflictFile, nil, 0644)
	defer os.Remove(conflictFile)

	_, err = os.OpenFile(conflictFile, os.O_CREATE|os.O_EXCL, 0644)
	if !os.IsExist(err) {
		t.Errorf("IsExist 断言失败：无法有效识别并发资源排他冲突")
	}

	// 3. 测试 NewSyscallError 和 IsPermission
	// 模拟一个伪造的文件权限错误
	basePermissionErr := os.ErrPermission
	wrappedErr := os.NewSyscallError("mock_write_call", basePermissionErr)

	if !os.IsPermission(wrappedErr) {
		t.Errorf("NewSyscallError 包装破坏了原有的上下文关联链路，导致 IsPermission 鉴定失效")
	}
}
