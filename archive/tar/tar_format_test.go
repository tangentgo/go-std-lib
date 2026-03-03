package main_test

import (
	"archive/tar"
	"testing"
)

// Test_Format_String 验证不同归档格式的字符串表示形式。
// 该测试详尽检查了 Go 1.26.0 中支持的所有 Format 常量，
// 确保其 String() 方法返回预期的标准化文本。
func Test_Format_String(t *testing.T) {
	tests := []struct {
		name     string
		format   tar.Format
		expected string
	}{
		{"Unknown Format", tar.FormatUnknown, "<unknown>"},
		{"USTAR Format", tar.FormatUSTAR, "USTAR"},
		{"PAX Format", tar.FormatPAX, "PAX"},
		{"GNU Format", tar.FormatGNU, "GNU"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			result := tc.format.String()
			if result != tc.expected {
				t.Errorf("Format.String() 返回值异常: 期望获得 %q, 实际获得 %q", tc.expected, result)
			}
		})
	}
}
