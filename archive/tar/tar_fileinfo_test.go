package main_test

import (
	"archive/tar"
	"io/fs"
	"testing"
	"time"
)

// mockFSInfo 是一个用于测试的自定义文件信息结构，
// 它同时满足了 fs.FileInfo 和 Go 1.23 引入的 tar.FileInfoNames 接口。
type mockFSInfo struct {
	name    string
	size    int64
	mode    fs.FileMode
	modTime time.Time
	isDir   bool
}

// 实现 fs.FileInfo 接口标准方法
func (m mockFSInfo) Name() string       { return m.name }
func (m mockFSInfo) Size() int64        { return m.size }
func (m mockFSInfo) Mode() fs.FileMode  { return m.mode }
func (m mockFSInfo) ModTime() time.Time { return m.modTime }
func (m mockFSInfo) IsDir() bool        { return m.isDir }
func (m mockFSInfo) Sys() any           { return nil }

// 实现 tar.FileInfoNames 接口扩展方法
func (m mockFSInfo) Uname() (string, error) { return "admin_user", nil }
func (m mockFSInfo) Gname() (string, error) { return "staff_group", nil }

// Test_FileInfoHeader_And_FileInfo_Lifecycle 验证元数据对象在不同系统域之间的转换保真度。
func Test_FileInfoHeader_And_FileInfo_Lifecycle(t *testing.T) {
	// 初始化模拟的时间戳与文件属性
	baseTime := time.Date(2026, 2, 10, 10, 0, 0, 0, time.UTC)
	fi := mockFSInfo{
		name:    "system_config.yaml",
		size:    4096,
		mode:    0640,
		modTime: baseTime,
		isDir:   false,
	}

	// 测试 1: FileInfoHeader (从文件系统域转换到 Tar 归档域)
	// 此函数负责提取系统属性并组装成 tar 规范的 Header 对象。
	hdr, err := tar.FileInfoHeader(fi, "")
	if err != nil {
		t.Fatalf("FileInfoHeader 调用失败并返回错误: %v", err)
	}

	// 验证基础字段映射
	if hdr.Name != "system_config.yaml" {
		t.Errorf("Header 名称未正确映射，期望 system_config.yaml，实际得到 %s", hdr.Name)
	}
	if hdr.Size != 4096 {
		t.Errorf("Header 大小未正确映射，期望 4096，实际得到 %d", hdr.Size)
	}

	// 验证 Go 1.23+ FileInfoNames 接口的自动提取机制
	if hdr.Uname != "admin_user" || hdr.Gname != "staff_group" {
		t.Errorf("未能通过 FileInfoNames 接口正确提取用户标识，当前 Uname: %s, Gname: %s", hdr.Uname, hdr.Gname)
	}

	// 测试 2: (*Header) FileInfo (从 Tar 归档域逆向回文件系统域)
	// 将组装好的 Header 重新转换为标准库通用的 fs.FileInfo 接口形态。
	reconstructedFI := hdr.FileInfo()

	// 验证逆向转换后的数据保真度
	if reconstructedFI.Name() != "system_config.yaml" {
		t.Errorf("重建后的 FileInfo 名称不匹配，得到 %s", reconstructedFI.Name())
	}
	if reconstructedFI.Size() != 4096 {
		t.Errorf("重建后的 FileInfo 大小不匹配，得到 %d", reconstructedFI.Size())
	}
	// 归档模式转换可能会涉及位运算的舍入或掩码处理，但核心权限位应被保留
	if reconstructedFI.Mode().Perm() != 0640 {
		t.Errorf("重建后的 FileInfo 权限位异常，得到 %v", reconstructedFI.Mode().Perm())
	}
}
