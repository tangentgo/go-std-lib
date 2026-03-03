package main_test

import (
	"archive/zip"
	"bytes"
	"errors"
	"testing"
)

// TestConstantsAndVariables 演示常量的逻辑判断与错误变量的拦截
func TestConstantsAndVariables(t *testing.T) {
	// 验证常量定义
	if zip.Store != 0 || zip.Deflate != 8 {
		t.Fatalf("核心压缩常量被意外篡改")
	}

	// 场景 1: 模拟 ErrFormat
	invalidData := []byte("这是一段伪造的、不符合 PKWARE 规范的随机数据流")
	_, err := zip.NewReader(bytes.NewReader(invalidData), int64(len(invalidData)))
	if err == nil {
		t.Fatal("预期遭遇解析失败，但却成功实例化了 Reader")
	}
	// 工程规范：使用 errors.Is 进行解包判定
	if !errors.Is(err, zip.ErrFormat) {
		t.Errorf("预期错误类型为 ErrFormat,实际得到: %v", err)
	} else {
		t.Log("成功拦截不合法 ZIP 格式异常")
	}
}
