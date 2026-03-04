package main_test

import (
	"bytes"
	"io"
	"testing"
)

func TestWriteString_Usage(t *testing.T) {
	// bytes.Buffer 内部实现了高效的 WriteString 方法
	var optimalDest bytes.Buffer

	// 利用高速通道倾泻字符串，在此路径下内存开销趋近于零
	n, err := io.WriteString(&optimalDest, "High Performance Log Output")

	if err != nil {
		t.Fatalf("io.WriteString 写入失败: %v", err)
	}

	expectedStr := "High Performance Log Output"
	if n != len(expectedStr) || optimalDest.String() != expectedStr {
		t.Errorf("利用高速通道转储的字符串内容发生偏差")
	}
}

// **核心机制**：函数签名 `WriteString(w Writer, s string) (n int, err error)`。
// 当被调用时，它会在运行时进行动态反射级别的接口探针检查：
// 如果底层目标 `w` 原生支持 Go 1.12 引入的 `StringWriter` 接口
// （例如 `bytes.Buffer` 和 `strings.Builder` 就原生支持），
// 它将直接移交并无损地灌入字符串；仅当目标非常原始且不具备该能力时，
// 它才会作为兼容性兜底方案，无可奈何地进行那次昂贵的 `byte` 类型分配与转换 。
