package main_test

import (
	"archive/tar"
	"bytes"
	"errors"
	"testing"
)

// Test_Writer_Workflow 验证归档写入组件的初始化、属性注入、载荷管控与安全封印流程。
func Test_Writer_Workflow(t *testing.T) {
	var buf bytes.Buffer
	// 1. 初始化 Writer
	tw := tar.NewWriter(&buf)

	// 构造测试文件头。使用显式指针实例化。
	hdr := &tar.Header{
		Name:     "config/settings.ini",
		Mode:     0644,
		Size:     10, // 严密声明后续仅容许写入 10 个字节
		Typeflag: tar.TypeReg,
	}

	// 2. 测试 (*Writer) WriteHeader - 注入元数据并建立状态栅栏
	if err := tw.WriteHeader(hdr); err != nil {
		t.Fatalf("写入 Header 元数据区时失败: %v", err)
	}

	// 3. 测试 (*Writer) Write - 写入合法数据载荷
	validData := []byte("0123456789") // 正好 10 字节
	n, err := tw.Write(validData)
	if err != nil {
		t.Fatalf("执行写入流操作时遭遇故障: %v", err)
	}
	if n != 10 {
		t.Errorf("写入字节数统计错乱，预期 10，实际报告 %d", n)
	}

	// 验证越界防御机制：尝试写入超出 Header.Size 范围的第 11 个字节
	_, err = tw.Write([]byte("A"))
	if !errors.Is(err, tar.ErrWriteTooLong) {
		t.Errorf("越界拦截失效！期望触发 tar.ErrWriteTooLong，实际获得: %v", err)
	}

	// 4. 测试 (*Writer) Flush - 手动干预区块对齐
	// 此处主动调用以填平剩余的 502 个字节（512 - 10）为全零填充符
	if err := tw.Flush(); err != nil {
		t.Fatalf("手动对齐块边界时发生错误: %v", err)
	}

	// 5. 测试 (*Writer) Close - 执行终止协议并封印结构
	if err := tw.Close(); err != nil {
		t.Fatalf("归档封印操作崩溃: %v", err)
	}

	// 验证封印后拦截机制：封闭状态机不再受理新文件的注入
	err = tw.WriteHeader(&tar.Header{Name: "illegal_late_entry.txt", Size: 0})
	if !errors.Is(err, tar.ErrWriteAfterClose) {
		t.Errorf("封印状态违规检测失效！期望触发 tar.ErrWriteAfterClose，实际获得: %v", err)
	}

	// 验证最终产物的物理结构：
	// Header 区块(512) + 数据及填充区块(512) + EOF终止区块*2(1024) = 2048 字节
	if buf.Len() != 2048 {
		t.Errorf("最终归档尺寸计算违背协议，期望获取 2048 字节，实际得到 %d 字节", buf.Len())
	}
}
