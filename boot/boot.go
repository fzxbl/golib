package iboot

import (
	"io"
)

var closeFns []func() error

// TryRegisterCloser 注册 close 方法
func TryRegisterCloser(comp interface{}) {
	if c, ok := comp.(io.Closer); ok {
		closeFns = append(closeFns, c.Close)
		return
	}
	if fn, ok := comp.(func() error); ok {
		closeFns = append(closeFns, fn)
	}
}

// BeforeExit 退出前执行，资源清理、日志落盘等
func BeforeExit() {
	for _, fn := range closeFns {
		_ = fn()
	}
}
