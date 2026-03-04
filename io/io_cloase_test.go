package main_test

import (
	"io"
	"strings"
	"testing"
)

func TestNopCloser_Usage(t *testing.T) {
	// 这是一个只拥有 Read 方法、不含任何必须清理资源的纯净内存池
	pureMemoryStream := strings.NewReader("Fake HTTP Request Body Payload")

	// 发起升格伪装，披上 ReadCloser 的外衣，以通过高级 API 的类型严控审查
	fullyCompliantStream := io.NopCloser(pureMemoryStream)

	// 业务框架习惯性地在 defer 中调用销毁指令
	err := fullyCompliantStream.Close()

	// 断言：空壳销毁指令必须绝对平滑过渡，不引发任何实质错误
	if err != nil {
		t.Fatalf("伪装器缺陷：NopCloser 的销毁调用发生了非预期的错误反馈: %v", err)
	}

	// 证实伪装并未破坏其实质的传输血脉
	survivingData, _ := io.ReadAll(fullyCompliantStream)
	if string(survivingData) != "Fake HTTP Request Body Payload" {
		t.Errorf("伪装行动扭曲了原本的数据流向")
	}
}
