package main_test

import (
	"os"
	"testing"
)

func TestEnvironmentVariablesManagement(t *testing.T) {
	testKey := "GO_126_TEST_ENV"
	testVal := "active_state"

	// 1. 测试 Setenv
	if err := os.Setenv(testKey, testVal); err != nil {
		t.Fatalf("Setenv 执行失败: %v", err)
	}

	// 2. 测试 Getenv
	if got := os.Getenv(testKey); got != testVal {
		t.Errorf("Getenv 返回值异常，期望 %s，实际 %s", testVal, got)
	}

	// 3. 测试 LookupEnv
	emptyKey := "GO_126_EMPTY_ENV"
	os.Setenv(emptyKey, "")
	val, exists := os.LookupEnv(emptyKey)
	if !exists || val != "" {
		t.Errorf("LookupEnv 无法正确识别值为空但实际存在的环境变量")
	}

	// 4. 测试 Environ
	envList := os.Environ()
	found := false
	for _, e := range envList {
		if e == testKey+"="+testVal {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Environ 返回的列表中未包含新设置的环境变量")
	}

	// 5. 测试 Unsetenv
	os.Unsetenv(testKey)
	if _, exists := os.LookupEnv(testKey); exists {
		t.Errorf("Unsetenv 执行后，环境变量依然存在")
	}

	// 注意：Clearenv 具有破坏性，为防止影响系统的其他测试组件，不在此处调用
}

func TestEnvironmentExpansion(t *testing.T) {
	os.Setenv("LOG_LEVEL", "DEBUG")
	os.Setenv("APP_NAME", "GoServer")
	defer os.Unsetenv("LOG_LEVEL")
	defer os.Unsetenv("APP_NAME")

	// 测试 ExpandEnv
	input := "Service ${APP_NAME} is running at $LOG_LEVEL level."
	expected := "Service GoServer is running at DEBUG level."
	if got := os.ExpandEnv(input); got != expected {
		t.Errorf("ExpandEnv 替换失败，得到: %s", got)
	}

	// 测试 Expand
	customInput := "Token is ${TOKEN}"
	gotCustom := os.Expand(customInput, func(key string) string {
		if key == "TOKEN" {
			return "xyz123"
		}
		return "unknown"
	})
	if gotCustom != "Token is xyz123" {
		t.Errorf("Expand 自定义映射失败，得到: %s", gotCustom)
	}
}
