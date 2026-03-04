package main_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestCopy_Usage(t *testing.T) {
	// 模拟一个包含大段有效数据的外部数据源 (Reader)
	sourceData := "Go 1.26.0 IO Copy Mechanism Test Data Stream."
	src := strings.NewReader(sourceData)

	// 模拟一个内存数据接收端 (Writer)
	var dst bytes.Buffer

	// 启动转储引擎，并记录流经的字节总数
	written, err := io.Copy(&dst, src)

	// 工程规范：Copy 成功遇到 EOF 时，返回的 err 应当为 nil，而不是 EOF
	if err != nil {
		t.Fatalf("io.Copy 在流转过程中发生意外错误: %v", err)
	}

	// 状态断言：校验写入的字节数是否与源长度一致
	if written != int64(len(sourceData)) {
		t.Errorf("数据传输截断：期望传输 %d 字节，实际传输 %d 字节", len(sourceData), written)
	}

	// 内容断言：校验数据保真度
	if dst.String() != sourceData {
		t.Errorf("数据损坏：目的端接收到的内容与源端不符")
	}
}
func TestCopyBuffer_Usage(t *testing.T) {
	src := strings.NewReader("Controlled Memory Allocation Stream Transfer")
	var dst bytes.Buffer

	// 严苛内存控制策略：刻意预分配一个极小尺寸的缓冲区（8字节）
	// 这在底层会强制 CopyBuffer 进行多次微小的循环读取
	reusableBuffer := make([]byte, 8)

	// 执行受控拷贝
	written, err := io.CopyBuffer(&dst, src, reusableBuffer)

	if err != nil {
		t.Fatalf("io.CopyBuffer 执行失败: %v", err)
	}

	expectedStr := "Controlled Memory Allocation Stream Transfer"
	if written != int64(len(expectedStr)) || dst.String() != expectedStr {
		t.Errorf("使用外置极小缓冲区导致数据重组异常")
	}
}

func TestCopyN_Usage(t *testing.T) {
	// 模拟一个包含数百字节的冗长报文流
	src := strings.NewReader("Header_Data||Payload_Body_Extremely_Long...")
	var dst bytes.Buffer

	// 业务逻辑：协议规范指出，前 11 个字节是唯一的元数据
	targetLength := int64(11)
	written, err := io.CopyN(&dst, src, targetLength)

	if err != nil {
		t.Fatalf("io.CopyN 不应报错: %v", err)
	}

	// 状态验证：仅应精确提取 11 个字节
	if written != targetLength {
		t.Errorf("越界操作：期望复制 %d 字节，实际得到 %d 字节", targetLength, written)
	}
	if dst.String() != "Header_Data" {
		t.Errorf("提取内容篡改：得到异常数据 [%s]", dst.String())
	}

	// 异常断言测试：请求索取的长度超出了源的极限
	var dstOOB bytes.Buffer
	_, errOOB := io.CopyN(&dstOOB, src, 1000)
	// 此时应当引发预期的 EOF 中断
	if errOOB != io.EOF {
		t.Errorf("针对耗尽流的越界提取应当触发 EOF，但得到: %v", errOOB)
	}
}
