package iboot

import (
	"io"
)

var closeFns []func() error

// TryRegisterCloser 注册 close 方法
func TryRegisterCloser(close interface{}) {
	if c, ok := close.(io.Closer); ok {
		closeFns = append(closeFns, c.Close)
		return
	}
	if fn, ok := close.(func() error); ok {
		closeFns = append(closeFns, fn)
	}
}

// BeforeExit 退出前执行，资源清理、日志落盘等
func BeforeExit() {
	for _, fn := range closeFns {
		_ = fn()
	}
}
