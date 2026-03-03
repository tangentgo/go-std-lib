package main_test

import (
	"archive/zip"
	"bytes"
	"io"
	"os"
	"testing"
)

func TestReader_And_ReadCloser_AllFunctions(t *testing.T) {
	// 准备环节：在内存中构建一个合法的 ZIP 数据源供测试使用
	zipBuffer := new(bytes.Buffer)
	w := zip.NewWriter(zipBuffer)
	w.SetComment("Global Archive Architecture Comment")
	fWriter, _ := w.Create("config/settings.yaml")
	fWriter.Write([]byte("env: production\nversion: 1.0.0"))
	w.Close()

	// --- 测试 1: func NewReader ---
	// 结合 bytes.Reader 提供 io.ReaderAt 接口支持
	readerAt := bytes.NewReader(zipBuffer.Bytes())
	r, err := zip.NewReader(readerAt, int64(zipBuffer.Len()))
	if err != nil {
		t.Fatalf("NewReader 初始化失败: %v", err)
	}

	// 验证全局注释提取
	if r.Comment != "Global Archive Architecture Comment" {
		t.Errorf("未能正确提取归档注释")
	}

	// --- 测试 2: func (*Reader) Open ---
	// 验证 fs.FS 接口的无缝衔接
	fsFile, err := r.Open("config/settings.yaml")
	if err != nil {
		t.Fatalf("通过 fs.FS 接口查找文件失败: %v", err)
	}
	defer fsFile.Close()

	fsInfo, _ := fsFile.Stat()
	if fsInfo.Name() != "settings.yaml" {
		t.Errorf("fs.File 抽象层 Name 解析错误")
	}

	// --- 测试 3: func (*Reader) RegisterDecompressor ---
	// 为此 Reader 实例注册一个空解压器（仅用于方法签名调用演示）
	dummyDecompressor := func(reader io.Reader) io.ReadCloser {
		return io.NopCloser(reader)
	}
	r.RegisterDecompressor(99, dummyDecompressor) // 99 为假定的算法 ID

	// --- 测试 4: func OpenReader 与 func (*ReadCloser) Close ---
	// 由于涉及到真实文件操作，此处演示核心逻辑流
	// 假设我们在磁盘写入了一个临时文件 tmp.zip
	tmpFile := "test_open_reader.zip"
	os.WriteFile(tmpFile, zipBuffer.Bytes(), 0644)
	defer os.Remove(tmpFile) // 清理现场

	rc, err := zip.OpenReader(tmpFile)
	if err != nil {
		t.Fatalf("OpenReader 操作系统级别读取失败: %v", err)
	}
	// 关键：必须关闭释放 fd
	err = rc.Close()
	if err != nil {
		t.Errorf("ReadCloser.Close 释放资源失败: %v", err)
	}
}
