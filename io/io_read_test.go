package main_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

// ReadAll 它在内部建立一个持续扩容的字节切片池，疯狂地向传入的 `Reader` 索取数据，直至彻底撞上 `EOF` 之墙或遭遇无法逾越的致命错误。

func TestReadAll_Usage(t *testing.T) {
	// 模拟一个待全量加载的配置文件数据源
	configStream := strings.NewReader("server_port=8080\nlog_level=debug")

	// 执行全量内存吸取 (享受 Go 1.26 带来的性能飙升红利)
	configBytes, err := io.ReadAll(configStream)

	// 成功断言：不应将 EOF 泄露给上层业务逻辑
	if err != nil {
		t.Fatalf("io.ReadAll 未能妥善处理流结束标志或发生错误: %v", err)
	}

	expectedStr := "server_port=8080\nlog_level=debug"
	if string(configBytes) != expectedStr {
		t.Errorf("全量加载数据损坏，丢失边界完整性")
	}
}

// - `ReadAtLeast(r Reader, bufbyte, min int) (n int, err error)`：其内部部署了一个 `for` 循环不断抽吸。
// 如果在达标（`min` 字节）之前流就干涸了，系统将抛出 `ErrUnexpectedEOF` 以告警数据的残酷截断。
// 同时，它具备智能防御机制：如果传入的 `buf` 物理尺寸甚至装不下 `min` 个字节，它会立即判定逻辑失效并拒绝执行，抛出 `ErrShortBuffer` 错误 。
// - `ReadFull(r Reader, bufbyte) (n int, err error)`：其底层机制直接映射为对 `ReadAtLeast(r, buf, len(buf))`
// 的无缝转发包装 。只要最终获取的数据量哪怕比 `buf` 短缺了一个字节，它都会报出 `ErrUnexpectedEOF`。

func TestWatermarkReads_Usage(t *testing.T) {
	// 模拟一个极其脆弱的短流
	shortStream := strings.NewReader("Tiny")

	// ----- 测试 io.ReadAtLeast -----
	buf1 := make([]byte, 10)
	// 试图从 "Tiny" (4字节) 中强行逼取至少 6 个字节
	n1, err1 := io.ReadAtLeast(shortStream, buf1, 6)

	// 断言：应当捕获到非预期的流中断异常
	if err1 != io.ErrUnexpectedEOF { // 缓冲区够大，但底层流没给够数据
		t.Errorf("io.ReadAtLeast 未能正确报告数据不足，得到错误: %v", err1)
	}
	if n1 != 4 {
		t.Errorf("尽管失败，但仍应如实报告已窃取到的 4 个残存字节")
	}

	// 测试极端防护：提供不合格的缓冲区
	bufTooSmall := make([]byte, 2)
	_, errSize := io.ReadAtLeast(shortStream, bufTooSmall, 6)
	if errSize != io.ErrShortBuffer { // 缓冲区本身就不可能满足 min
		t.Errorf("缺乏短缓冲区防护，期待 ErrShortBuffer")
	}

	// ----- 测试 io.ReadFull -----
	perfectStream := strings.NewReader("Exact668By")
	buf2 := make([]byte, 8)

	n2, err2 := io.ReadFull(perfectStream, buf2)
	if err2 != nil || n2 != 8 {
		t.Fatalf("io.ReadFull 无法填满刚好相符的缓冲区: %v", err2)
	}
}

// **核心机制**：`LimitReader(r Reader, n int64) Reader` 将目标包裹在一个名为 `LimitedReader` 的内部结构体中 。
// 该结构体包含一个单调递减的计数器 `N`。每次上层发起读取时，它会对比申请读取的长度与剩余的 `N` 并进行必要的截短；
// 一旦 `N` 归零，无论下层的原初网络连接是否依然活跃，它都会决绝地向上抛出不可逆的 `EOF` 信号 。
func TestLimitReader_Usage(t *testing.T) {
	// 模拟一个包含危险后门或极长数据的输入流
	maliciousStream := strings.NewReader("SafeData_FollowedBy_Infinite_Malicious_Garbage")

	// 设置安保防线：该通道终生最高只能流出 8 个字节
	secureReader := io.LimitReader(maliciousStream, 8)

	// 使用全量吸取器进行压力测试
	data, err := io.ReadAll(secureReader)

	if err != nil {
		t.Fatalf("包裹层不应导致异常崩溃: %v", err)
	}

	// 状态断言：它被无情地截断了
	if string(data) != "SafeData" {
		t.Errorf("隔离墙失效：期待截断为 'SafeData'，却得到了 '%s'", string(data))
	}
}

func TestTeeReader_Usage(t *testing.T) {
	mainSource := strings.NewReader("Secret Corporate Financial Report")

	// 部署一个旁路监听器 (模拟哈希指纹计算器或监控审计模块)
	var auditor bytes.Buffer

	// 在主流与调用方之间安插窃听网关
	interceptedStream := io.TeeReader(mainSource, &auditor)

	// 终端正常执行业务逻辑，毫无察觉地将流抽干
	businessData, _ := io.ReadAll(interceptedStream)

	// 断言验证：不仅业务端拿到完整数据，审计端也必须积攒了相同的数据留存
	if string(businessData) != "Secret Corporate Financial Report" {
		t.Errorf("TeeReader 损坏了主干业务数据")
	}
	if auditor.String() != "Secret Corporate Financial Report" {
		t.Errorf("旁路监听网关失灵，未能拦截并同步存下数据的完整镜像")
	}
}

func TestSectionAndOffset_Usage(t *testing.T) {
	// ----- NewSectionReader 并发安全读取测试 -----
	// 模拟一段具备并发访问潜质的超长只读大内存数据源
	massiveDataSource := strings.NewReader("PREFIX_SEGMENT_1_SEGMENT_2_SUFFIX")

	// 锁定坐标：忽略前缀，精确切割出 SEGMENT_1 进行独立处理
	// 偏移量为 7，长度限制为 9 字节
	windowedReader := io.NewSectionReader(massiveDataSource, 7, 9)

	extractedPart, _ := io.ReadAll(windowedReader)
	if string(extractedPart) != "SEGMENT_1" {
		t.Errorf("安全切片视窗圈定失败，截取为: %s", string(extractedPart))
	}

	// 调用 Go 1.22.0 增补的透视分析方法
	origin, offset, limit := windowedReader.Outer()
	if origin != massiveDataSource || offset != 7 || limit != 9 {
		t.Errorf("Outer 透视接口反馈元数据有误")
	}

	// ----- (概念说明) NewOffsetWriter 并发写入场景 -----
	// 因标准库缺省在内存中实现了 WriterAt 接口的简便结构，
	// 此处主要展示其对于持久化并发追加的核心思想。
	// 在真实场景中，往往挂载的是 os.File 对象以实现日志簇文件的并发安全落盘。
}
