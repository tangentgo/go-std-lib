package main_test

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"
)

// Test_Reader_Lifecycle 验证流式读取组件的初始化、迭代以及数据提取功能。
func Test_Reader_Lifecycle(t *testing.T) {
	// 阶段一：在内存中合成一个微型的合法 tar 归档作为测试载体。
	// 这里预先植入两个实体：一个包含数据的文本文件，一个无数据的目录结构。
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// 写入文件实体
	err := tw.WriteHeader(&tar.Header{
		Name:     "data.txt",
		Mode:     0600,
		Size:     14, // 预期写入 14 字节
		Typeflag: tar.TypeReg,
	})
	if err != nil {
		t.Fatalf("准备测试数据阶段失败 (WriteHeader): %v", err)
	}
	tw.Write([]byte("TarballContent")) // 写入精确的 14 字节

	// 写入目录实体
	err = tw.WriteHeader(&tar.Header{
		Name:     "assets/",
		Mode:     0755,
		Size:     0, // 目录体必须具备零尺寸的特征
		Typeflag: tar.TypeDir,
	})
	if err != nil {
		t.Fatalf("准备测试目录阶段失败: %v", err)
	}
	tw.Close() // 封印归档，写入尾部终止符

	// 保存到磁盘
	// if err := os.WriteFile("out.tar", buf.Bytes(), 0644); err != nil {
	// 	t.Fatalf("保存 out.tar 失败: %v", err)
	// }
	// 阶段二：正式测试 Reader 的系列方法
	// 1. 测试 tar.NewReader
	tr := tar.NewReader(&buf)

	// 2. 测试 (*Reader) Next - 定位首个文件
	h, err := tr.Next()
	if err != nil {
		t.Fatalf("tr.Next() 在定位首个文件时遭遇未预期错误: %v", err)
	}
	if h.Name != "data.txt" || h.Size != 14 {
		t.Errorf("解析出的 Header 元数据失真，获得名称: %s, 大小: %d", h.Name, h.Size)
	}

	// 3. 测试 (*Reader) Read - 提取实体载荷
	dataBuffer := make([]byte, 100)
	n, err := tr.Read(dataBuffer)
	// 在此处，Read 若未触及数据边界可能返回 (n, nil)，也可能在读完时返回 (n, io.EOF)
	if err != nil && err != io.EOF {
		t.Fatalf("tr.Read() 提取数据块失败: %v", err)
	}
	extractedContent := string(dataBuffer[:n])
	if extractedContent != "TarballContent" {
		t.Errorf("数据提取产生损坏，预期 'TarballContent'，实际获得 '%s'", extractedContent)
	}

	// 再次调用 Next，跨越目录实体
	dirHeader, err := tr.Next()
	if err != nil {
		t.Fatalf("tr.Next() 导航至目录实体时发生故障: %v", err)
	}
	if dirHeader.Typeflag != tar.TypeDir || dirHeader.Name != "assets/" {
		t.Errorf("目录实体的元数据识别错误")
	}

	// 再次调用 Next，预期触达归档终止符
	_, err = tr.Next()
	if err != io.EOF {
		t.Errorf("迭代未能安全终止，期望捕获 io.EOF，实际截获 %v", err)
	}
}
