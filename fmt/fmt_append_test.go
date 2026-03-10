package main_test

import (
	"fmt"
	"testing"
)

func TestAppend(t *testing.T) {
	// 底层原理：预分配容量为 64 的底层数组
	buf := make([]byte, 0, 64)
	// 将结果追加至底部的可用容量内，避免了新建字符串引发的逃逸分配
	buf = fmt.Append(buf, "SessionID:", 999)
	fmt.Print("TestAppend切片内容: ", string(buf), "\n")
}
func TestAppendf(t *testing.T) {
	buf := make([]byte, 0, 128)
	// 底层原理：状态机执行期间，直接利用传入切片的内存作为暂存区
	buf = fmt.Appendf(buf, "Request %s processed in %d ms", "/api/v1", 42)
	fmt.Print("TestAppendf切片内容: ", string(buf), "\n")
}
func TestAppendln(t *testing.T) {
	buf := make([]byte, 0, 64)
	// 底层原理：在切片尾部插入参数的字面量，并在最后推入 '\n'
	buf = fmt.Appendln(buf, "Transaction", "Committed", "OK")
	fmt.Print("TestAppendln切片内容: ", string(buf))
}
