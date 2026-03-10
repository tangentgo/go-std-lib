package main_test

import (
	"fmt"
	"strings"
	"testing"
)

func TestScan(b *testing.T) {
	var firstName string
	var age int
	// 底层原理：对 os.Stdin 发起 read 阻塞，读取后依靠状态机以空白符为界截断 Token。
	// 注意：实际执行时需在终端输入，如 "John 30"
	fmt.Print("TestScan 请输入姓名和年龄 (空格分隔): ")
	_, _ = fmt.Scan(&firstName, &age)
	fmt.Printf("扫描得到: %s, %d岁\n", firstName, age)
	fmt.Println("TestScan 已就绪 (注释以跳过终端阻塞)")
}
func TestScanf(t *testing.T) {
	var ip string
	var port int
	// 底层原理：FSM要求输入严格符合 "IP:XXX PORT:XXX" 格式，否则返回 err
	fmt.Print("TestScanf 请按格式输入 IP:127.0.0.1 PORT:8080 : ")
	_, _ = fmt.Scanf("IP:%s PORT:%d", &ip, &port)
	fmt.Printf("扫描结果: IP=%s, Port=%d\n", ip, port)
	fmt.Println("TestScanf 已就绪 (注释以跳过终端阻塞)")
}
func TestScanln(t *testing.T) {
	var item string
	var price float64
	// 底层原理：内部的 nlIsEnd 标志位被设置为 true，一旦 ReadRune 探测到 '\n'，立即中止状态机。
	fmt.Print("TestScanln 请在一行内输入商品和价格: ")
	_, _ = fmt.Scanln(&item, &price)
	fmt.Printf("读取到: %s 价格: %g\n", item, price)
	fmt.Println("TestScanln 已就绪 (注释以跳过终端阻塞)")
}
func TestFscan(t *testing.T) {
	// 使用 strings.NewReader 模拟一个实现了 io.Reader 的流
	reader := strings.NewReader("SystemA 99.9")
	var node string
	var uptime float64
	// 底层原理：从 reader 中循环调用 Read 方法提取字节
	n, err := fmt.Fscan(reader, &node, &uptime)
	if err == nil {
		fmt.Printf("TestFscan 成功提取 %d 个变量: 节点=%s, 在线率=%g\n", n, node, uptime)
	}
}
func TestFscanf(t *testing.T) {
	reader := strings.NewReader("MAC: fa:bc:12:34:56:78, Status: Active")
	var mac string
	var status string
	// 底层原理：FSM 验证流内容是否匹配 "MAC: %s, Status: %s"
	n, err := fmt.Fscanf(reader, "MAC: %s, Status: %s", &mac, &status)
	if err == nil {
		fmt.Printf("TestFscanf 成功提取 %d 项: MAC=%s, 状态=%s\n", n, mac, status)
	}
}
func TestFscanln(t *testing.T) {
	reader := strings.NewReader("Header Value\nIgnored Data")
	var k, v string
	// 底层原理：只扫描遇到 '\n' 前的元素，丢弃 "Ignored Data"
	n, err := fmt.Fscanln(reader, &k, &v)
	if err == nil {
		fmt.Printf("TestFscanln 首行读取了 %d 个值: %s, %s\n", n, k, v)
	}
}
func TestSscan(t *testing.T) {
	data := "True 1024"
	var flag bool
	var size int
	// 底层原理：零系统调用，使用内部 stringReader 避免内存拷贝，利用反射装载值
	n, err := fmt.Sscan(data, &flag, &size)
	if err == nil {
		fmt.Printf("TestSscan 解包了 %d 个数据: bool=%t, int=%d\n", n, flag, size)
	}
}
