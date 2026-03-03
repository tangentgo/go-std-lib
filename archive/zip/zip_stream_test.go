package main_test

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func TestFile_Methods_DataOffset_Open_OpenRaw(t *testing.T) {
	// 构造测试数据
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// 强制开启 DEFLATE 压缩，以凸显 Open 与 OpenRaw 的差异
	header := &zip.FileHeader{
		Name:   "data.txt",
		Method: zip.Deflate,
	}

	writer, err := w.CreateHeader(header)
	if err != nil {
		t.Fatalf("CreateHeader 失败: %v", err)
	}

	// 写入极易被压缩的冗余数据
	originalData := bytes.Repeat([]byte("GoArchiveZipTest"), 100)
	if _, err := writer.Write(originalData); err != nil {
		t.Fatalf("写入原始数据失败: %v", err)
	}

	if err := w.Close(); err != nil {
		t.Fatalf("关闭 zip writer 失败: %v", err)
	}

	// 读取 zip
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("NewReader 失败: %v", err)
	}

	if len(r.File) != 1 {
		t.Fatalf("期望 zip 内只有 1 个文件，实际: %d", len(r.File))
	}
	f := r.File[0] // 获取唯一的文件指针

	// --- 测试 1: func (*File) DataOffset ---
	offset, err := f.DataOffset()
	if err != nil {
		t.Fatalf("DataOffset 计算失败: %v", err)
	}
	if offset <= 0 {
		t.Errorf("数据偏移量异常,期望大于0,实际得到: %d", offset)
	}
	t.Logf("数据绝对偏移量定位在第 %d 字节", offset)

	// --- 测试 2: func (*File) Open ---
	// 读取解压后的真实数据
	rc, err := f.Open()
	if err != nil {
		t.Fatalf("Open 解压流初始化失败: %v", err)
	}
	extractedData, err := io.ReadAll(rc)
	_ = rc.Close()
	if err != nil {
		t.Fatalf("ReadAll(Open) 失败: %v", err)
	}

	if !bytes.Equal(extractedData, originalData) {
		t.Errorf("Open 解压后的数据与原始数据不一致")
	}

	// --- 测试 3: func (*File) OpenRaw ---
	// 直接截取底层的 DEFLATE 压缩字节流
	rawReader, err := f.OpenRaw()
	if err != nil {
		t.Fatalf("OpenRaw 原始流获取失败: %v", err)
	}
	rawBytes, err := io.ReadAll(rawReader)
	if err != nil {
		t.Fatalf("ReadAll(OpenRaw) 失败: %v", err)
	}

	t.Logf("解压后尺寸: %d bytes, 压缩态尺寸: %d bytes", len(extractedData), len(rawBytes))

	// 验证 OpenRaw 获取的数据确实经过了压缩，尺寸远小于明文
	if len(rawBytes) >= len(extractedData) {
		t.Errorf("OpenRaw 逻辑错误，获取到的似乎不是压缩流")
	}
}
