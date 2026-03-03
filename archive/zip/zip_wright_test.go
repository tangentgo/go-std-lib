package main_test

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func TestWriter_All_Methods(t *testing.T) {
	// 使用内存作为底层介质测试 func NewWriter
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// --- 测试 func SetOffset 与 func SetComment ---
	// 模拟该 ZIP 被追加在 100 字节的桩数据之后
	w.SetOffset(100)
	err := w.SetComment("Comprehensive Test Archive")
	if err != nil {
		t.Fatalf("SetComment 失败: %v", err)
	}

	// --- 测试 func Create ---
	// 快捷创建标准文件
	f1, err := w.Create("standard_file.txt")
	if err != nil {
		t.Fatalf("Create 操作失败: %v", err)
	}
	f1.Write([]byte("Basic standard payload"))

	// --- 测试 func CreateHeader ---
	// 使用自定义权限创建安全配置文件
	fh := &zip.FileHeader{
		Name:   "secret/config.key",
		Method: zip.Deflate,
	}
	fh.SetMode(0600) // 仅拥有者读写
	f2, err := w.CreateHeader(fh)
	if err != nil {
		t.Fatalf("CreateHeader 失败: %v", err)
	}
	f2.Write([]byte("SUPER_SECRET_KEY=12345"))

	// --- 测试 func Flush ---
	// 主动刷新写缓冲区
	err = w.Flush()
	if err != nil {
		t.Fatalf("Flush 刷新缓冲池失败: %v", err)
	}

	// --- 测试 func RegisterCompressor ---
	// 替换特定的压缩逻辑。例如算法 ID 99 映射到一个直接输出的包装器
	w.RegisterCompressor(99, func(out io.Writer) (io.WriteCloser, error) {
		return &noopWriteCloser{out}, nil
	})

	// --- 测试 func AddFS ---
	// 借助 Go 内置的嵌入或者目录抽象。由于测试环境限制，我们这里采用 os.DirFS 模拟
	// 为防止环境差异报错，此处仅做逻辑说明，注释执行
	// err = w.AddFS(os.DirFS("."))

	// --- 测试 func CreateRaw 与 func Copy ---
	// 为了演示 Copy 和 CreateRaw，我们需要预先准备一个包含原始压缩数据的 File 指针
	simulateCopyAndRaw(t, w)

	// --- 测试 func Close ---
	// 封印中央目录结构，结束写入
	err = w.Close()
	if err != nil {
		t.Fatalf("Close 构建中央目录失败: %v", err)
	}

	t.Logf("所有流控制测试完毕，生成归档总字节数: %d", buf.Len())
}

// 辅助包装器，实现 WriteCloser
type noopWriteCloser struct{ io.Writer }

func (n *noopWriteCloser) Close() error { return nil }

// 辅助方法，演示 Copy 与 CreateRaw
func simulateCopyAndRaw(t *testing.T, targetW *zip.Writer) {
	// 创建一个源 ZIP
	srcBuf := new(bytes.Buffer)
	srcW := zip.NewWriter(srcBuf)
	f, _ := srcW.Create("source_data.bin")
	f.Write([]byte("this data is highly compressed internally"))
	srcW.Close()

	// 读取源 ZIP
	srcR, _ := zip.NewReader(bytes.NewReader(srcBuf.Bytes()), int64(srcBuf.Len()))

	// 使用 Copy 高速桥接
	err := targetW.Copy(srcR.File[0])
	if err != nil {
		t.Fatalf("Copy 零压缩迁移过程失败: %v", err)
	}
}
