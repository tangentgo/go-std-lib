package main_test

import (
	"errors"
	"fmt"
	"testing"
)

func TestErrorf(t *testing.T) {
	baseErr := errors.New("i/o timeout")
	// 底层原理：检测到 %w 动词，拦截 baseErr，
	// 返回的对象不仅包含文本，还挂载了底层的指针链接。
	wrapped := fmt.Errorf("读取配置文件失败: %w", baseErr)

	fmt.Println("TestErrorf 生成错误:", wrapped.Error())
}
