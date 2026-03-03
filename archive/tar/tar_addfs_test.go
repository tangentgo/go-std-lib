package main_test

import (
	"archive/tar"
	"bytes"
	"io"
	"testing"
	"testing/fstest"
)

// Test_Writer_AddFS 验证高层级的系统目录递归归档功能。
// 该用例利用标准库的 fstest.MapFS 构建虚拟文件树进行模拟打包。
func Test_Writer_AddFS(t *testing.T) {
	// 在内存域构建一个微型虚拟文件树
	virtualFS := fstest.MapFS{
		"index.html":       {Data: []byte("<html></html>"), Mode: 0644},
		"scripts/main.js":  {Data: []byte("console.log('hi');"), Mode: 0644},
		"scripts/empty.js": {Data: []byte(""), Mode: 0644}, // 模拟空文件边缘场景
	}

	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)

	// 测试 (*Writer) AddFS - 一键式打包虚拟文件树
	if err := tw.AddFS(virtualFS); err != nil {
		t.Fatalf("执行 AddFS 递归注入操作失败: %v", err)
	}

	// 收尾封印
	if err := tw.Close(); err != nil {
		t.Fatalf("执行 Close 封印失败: %v", err)
	}
	// if err := os.WriteFile("out.tar", buf.Bytes(), 0644); err != nil {
	// 	t.Fatalf("保存 out.tar 失败: %v", err)
	// }
	// 反向解构验证：利用 Reader 验证 AddFS 是否正确录入了所有层级的实体
	tr := tar.NewReader(&buf)
	var archivedPaths []string

	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("通过 Reader 逆向迭代文件时发生异常: %v", err)
		}
		archivedPaths = append(archivedPaths, hdr.Name)
	}

	// 根据 fs.FS 的层级遍历规则，预期的实体应包含目录和全部子文件
	// 注意：fstest 可能会视具体情况推断上级目录结构，这里至少要存在定义的文件
	foundCount := len(archivedPaths)
	if foundCount < 3 {
		t.Errorf("AddFS 未能打包预期数量的文件实体，仅捕获到 %d 个: %v", foundCount, archivedPaths)
	}
}
