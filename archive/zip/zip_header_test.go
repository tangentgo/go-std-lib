package main_test

import (
	"archive/zip"
	"io/fs"
	"testing"
	"time"
)

func TestFileHeader_AllMethods(t *testing.T) {
	// 1. 模拟一个操作系统级别的文件信息
	// 实际场景中通常通过 os.Stat("filepath") 获取
	mockTime := time.Date(2026, 3, 3, 12, 0, 0, 0, time.UTC)

	// 为了演示，我们手动构造一个 FileHeader 来模拟 FileInfoHeader 的行为
	header := &zip.FileHeader{
		Name:     "core_engine.so",
		Modified: mockTime,
		Method:   zip.Deflate, // 显式声明开启压缩
	}

	// 2. 测试权限位掩码映射 (SetMode 与 Mode)
	// 赋予属主读写执行权限，属组和其他用户读执行权限 (0755)
	targetMode := fs.FileMode(0755)
	header.SetMode(targetMode)

	extractedMode := header.Mode()
	if extractedMode != targetMode {
		t.Errorf("UNIX 权限位映射失败，期望 %v,实际得到 %v", targetMode, extractedMode)
	}

	// 3. 测试与 fs.FileInfo 抽象层的双向转换 (FileInfo)
	fileInfo := header.FileInfo()
	if fileInfo.Name() != "core_engine.so" {
		t.Errorf("FileInfo 接口 Name() 解析异常")
	}
	if fileInfo.Mode() != targetMode {
		t.Errorf("FileInfo 接口 Mode() 透传异常")
	}
	if !fileInfo.ModTime().Equal(mockTime) {
		t.Errorf("FileInfo 接口时间解析异常")
	}

	// 4. 测试废弃的时间方法 (作为向后兼容性验证)
	header.SetModTime(mockTime.Add(time.Hour))
	if !header.ModTime().Equal(mockTime.Add(time.Hour)) {
		t.Errorf("遗留的 SetModTime/ModTime 方法工作异常")
	}
}
