package main

import "fmt"

func main() {
	var firstName string
	var age int
	// 底层原理：对 os.Stdin 发起 read 阻塞，读取后依靠状态机以空白符为界截断 Token。
	// 注意：实际执行时需在终端输入，如 "John 30"
	fmt.Print("TestScan 请输入姓名和年龄 (空格分隔): ")
	n, err := fmt.Scan(&firstName, &age)
	if err != nil {
		fmt.Println(n, err.Error())
	}
	fmt.Printf("扫描得到: %s, %d岁\n", firstName, age)
	fmt.Printf("firstname: %s , age: %d \n", firstName, age)
}
