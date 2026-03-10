package main_test

import (
	"fmt"
	"testing"
)

func TestSprint(t *testing.T) {
	// 底层原理：触发字符串的堆内存分配，无系统调用
	result := fmt.Sprint("连接重试:", 5, "次")
	fmt.Print("TestSprint生成字符串: ", result, "\n")
}
func TestSprintf(t *testing.T) {
	// 底层原理：解析 %q 为字符串添加双引号，%g 紧凑输出浮点数。
	// 最终生成新字符串返回。
	query := fmt.Sprintf("SELECT * FROM users WHERE name=%q AND score>%g", "admin", 90.5)
	fmt.Print("TestSprintf生成的SQL: ", query, "\n")
}
func TestSprintln(t *testing.T) {
	// 底层原理：在缓冲末尾自动追加换行符，常用于构建包含多行的数据块
	line := fmt.Sprintln("Record:", "XYZ", 100)
	fmt.Print("TestSprintln单行内容: ", line)
}
