package main_test

import (
	"bytes"
	"fmt"
	"os"
	"testing"
)

func TestFprint(t *testing.T) {
	// 使用 bytes.Buffer 作为 io.Writer 的实现者
	var buf bytes.Buffer
	// 底层原理：渲染完毕后，调用 buf.Write() 将字节写入内存数组
	n, err := fmt.Fprint(&buf, "诊断信息: ", 404, " ", "Not Found")
	if err == nil {
		fmt.Print("TestFprint成功写入字节数: ", n, "\n")
	}
}
func TestFprintf(t *testing.T) {
	var buf bytes.Buffer
	// 底层原理：针对 %04d 进行整数宽度控制及前导零填充，
	// 结果直接通过 io.Writer 接口注入 buf 中。
	n, err := fmt.Fprintf(&buf, "用户UID: %04d, 权限: %s", os.Getuid(), os.Getenv("USER"))
	if err == nil {
		fmt.Print("TestFprintf结果: ", buf.String(), " (", n, " bytes)\n")
	}
}
func TestFprintln(t *testing.T) {
	var buf bytes.Buffer
	// 底层原理：确保流记录呈现行缓冲特性，方便下游以 \n 进行按行反序列化。
	n, err := fmt.Fprintln(&buf, "WARN", "Disk Space Low", 85.5)
	if err == nil {
		fmt.Print("n:", n, "\n", "TestFprintln流内容:\n", buf.String())
	}
}
