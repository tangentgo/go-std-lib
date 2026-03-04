package main_test

import (
	"io"
	"testing"
)

func TestPipe_Usage(t *testing.T) {
	// 创建一个无缓冲的同步管道
	pr, pw := io.Pipe()

	// 定义用于线程间通信的确认通道
	done := make(chan struct{})

	// 启动生产者协程 (异步生成数据)
	go func() {
		defer pw.Close() // 关键：写完必须主动关闭，否则消费者将永久死锁

		payload := []byte("Streaming large JSON payload via Pipe without RAM cost.")
		// 此处的 Write 会陷入同步阻塞，直到外部主协程开始拉取
		n, err := pw.Write(payload)

		if err != nil {
			t.Errorf("管道写入异常: %v", err)
		}
		if n != len(payload) {
			t.Errorf("管道未能同步全部数据")
		}
		close(done) // 释放信号
	}()

	// 主协程扮演消费者 (异步拉取数据)
	// 借用 io.ReadAll 不断抽取，直到写端触发 Close
	receivedData, err := io.ReadAll(pr)

	if err != nil {
		t.Fatalf("从管道读取数据遭遇失败: %v", err)
	}
	if string(receivedData) != "Streaming large JSON payload via Pipe without RAM cost." {
		t.Errorf("管道传输导致数据失真")
	}

	// 确保异步流程完整结束
	<-done
}

// 必须强调，在同一个协程内对 Pipe 进行读写是毫无意义且绝对会引发致命死锁（Deadlock）的。
