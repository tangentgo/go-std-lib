package main_test

import (
	"fmt"
	"testing"
)

// 构造一个包装类型进行演示
type Wrapper struct {
	Value int
}

func (w Wrapper) Format(f fmt.State, verb rune) {
	// 底层原理：从 State 环境中提取所有标志位、宽度和精度，
	// 重组为 "%05d" 这样的指令字符串，以备后续二次调度。
	formatDirective := fmt.FormatString(f, verb)
	fmt.Fprintf(f, "Wrapper: %d,formatDirective: %s", w.Value, formatDirective)
}

func TestFormatString(t *testing.T) {
	w := Wrapper{Value: 42}
	// 将触发 Wrapper.Format，并在其内部用 FormatString 重建 "%05d"
	fmt.Printf("TestFormatString 输出: %05d\n", w)
}
