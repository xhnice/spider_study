package utils

import (
	"bytes"
	"runtime"
	"sync_study/logger"
)

// HandleError -
func HandleError(err error) {
	if err != nil {
		logger.Error(err)
	}
}

// Recover -
func Recover(position string) {
	if err := recover(); err != nil {
		logger.Errorf("[%s]Panic: %s", position, err)
		logger.Errorf(string(PanicTrace(10)))
	}
}

// PanicTrace trace panic stack info
func PanicTrace(kb int) []byte {
	s := []byte("/src/runtime/panic.go")
	e := []byte("\ngoroutine")
	line := []byte("\n")
	stack := make([]byte, kb<<10)
	// Stack 调用goroutine的堆栈跟踪格式化为buf,并返回写入buf的字节数
	length := runtime.Stack(stack, true)
	start := bytes.Index(stack, s)
	if start < 0 {
		start = 0
	}

	stack = stack[start:length]
	start = bytes.Index(stack, line) + 1
	stack = stack[start:]
	end := bytes.LastIndex(stack, line)
	if end != -1 {
		stack = stack[:end]
	}
	end = bytes.Index(stack, e)
	if end != -1 {
		stack = stack[:end]
	}
	stack = bytes.TrimRight(stack, string(line))
	return stack
}
