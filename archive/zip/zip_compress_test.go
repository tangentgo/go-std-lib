package main_test

import (
	"archive/zip"
	"bytes"
	"io"
	"testing"
)

func TestPackage_RegisterCompressor_Decompressor(t *testing.T) {
	// 假设我们在使用一种专有算法：标识号为 88 的 "SuperFast-Z"
	const CustomMethodID uint16 = 88

	// 1. 注册全局解压器：func RegisterDecompressor
	zip.RegisterDecompressor(CustomMethodID, func(r io.Reader) io.ReadCloser {
		// 为了不引入第三方依赖造成编译阻断，此处利用标准库 io.NopCloser 模拟直通解码
		// 在真实业务中，此处将返回诸如 zstd.NewReader(r) 等对象
		return io.NopCloser(r)
	})

	// 2. 注册全局压缩器：func RegisterCompressor
	zip.RegisterCompressor(CustomMethodID, func(w io.Writer) (io.WriteCloser, error) {
		// 返回一个直通的 WriteCloser，作为模拟 "SuperFast-Z" 的压缩逻辑
		return &noopWriteCloser{w}, nil
	})

	// 3. 全链路闭环验证
	buf := new(bytes.Buffer)
	w := zip.NewWriter(buf)

	// 显式指定我们自定义的算法 ID
	fh := &zip.FileHeader{
		Name:   "custom_algo_test.dat",
		Method: CustomMethodID,
	}

	fw, err := w.CreateHeader(fh)
	if err != nil {
		t.Fatalf("使用自定义 Method 初始化头部失败: %v", err)
	}

	payload := []byte("Testing Custom Compressor & Decompressor Pipeline")
	fw.Write(payload)
	w.Close()

	// 触发解压器读取验证
	r, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		t.Fatalf("尝试读取注入了自定义算法的包失败: %v", err)
	}

	rc, err := r.File[0].Open()
	if err != nil {
		t.Fatalf("尝试利用自定义解压器解包失败: %v", err)
	}
	defer rc.Close()

	extracted, _ := io.ReadAll(rc)
	if !bytes.Equal(extracted, payload) {
		t.Errorf("基于自定义算法的端到端管道数据损坏")
	} else {
		t.Log("自定义算法全链路注册与调用成功")
	}
}
