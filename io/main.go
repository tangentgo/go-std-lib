package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	pr, pw := io.Pipe()

	// 生产者：不断往 PipeWriter 写数据（模拟实时产生的日志/数据）
	go func() {
		defer pw.Close() // 关闭写端，读端最终会收到 EOF
		for i := 1; i <= 5; i++ {
			line := fmt.Sprintf("line %d: hello pipe\n", i)
			_, err := io.WriteString(pw, line)
			if err != nil {
				// 如果读端提前关闭，写端这里会拿到错误
				log.Println("writer error:", err)
				return
			}
			time.Sleep(200 * time.Millisecond)
		}
	}()

	// 消费者：从 PipeReader 读数据，做流式处理（这里：转大写），再输出到 stdout
	// 你也可以把输出换成 gzip writer、网络连接、文件等。
	reader := bufio.NewReader(pr)
	for {
		s, err := reader.ReadString('\n')
		if len(s) > 0 {
			processed := strings.ToUpper(s)
			_, _ = io.WriteString(os.Stdout, processed)
		}
		if err != nil {
			if err == io.EOF {
				fmt.Println("DONE (EOF)")
				break
			}
			// 如果写端 CloseWithError，这里会拿到那个错误
			log.Println("reader error:", err)
			break
		}
	}
}
