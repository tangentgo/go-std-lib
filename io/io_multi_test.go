package main_test

import (
	"bytes"
	"io"
	"strings"
	"testing"
)

func TestMultiMechanisms_Usage(t *testing.T) {
	// ----- MultiReader 测试场景 -----
	r1 := strings.NewReader("Segment_A | ")
	r2 := strings.NewReader("Segment_B | ")
	r3 := strings.NewReader("Segment_C")

	// 构建逻辑级联视图
	mergedStream := io.MultiReader(r1, r2, r3)
	completeData, _ := io.ReadAll(mergedStream)

	expectedMerged := "Segment_A | Segment_B | Segment_C"
	if string(completeData) != expectedMerged {
		t.Errorf("MultiReader 拼接错乱: %s", string(completeData))
	}

	// ----- MultiWriter 测试场景 -----
	var displayConsole bytes.Buffer // 模拟屏幕输出
	var diskFile bytes.Buffer       // 模拟磁盘落地

	// 构建一转二的分发集线器
	broadcaster := io.MultiWriter(&displayConsole, &diskFile)

	msg := []byte("System Crash Alert!")
	_, err := broadcaster.Write(msg)

	if err != nil {
		t.Fatalf("广播分发出现故障: %v", err)
	}

	// 断言：双端必须毫无差别地接收到相同数据
	if displayConsole.String() != string(msg) || diskFile.String() != string(msg) {
		t.Errorf("MultiWriter 分发不均匀或存在丢失")
	}
}
