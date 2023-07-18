package irequest

import (
	"log"
	"net/http"
)

type ExampleRoundTripper struct{}

// RoundTrip 编辑这个函数可查看发送前的request，实现debug
func (lrt *ExampleRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	cookies := req.Header.Get("Cookie")
	// 在这里，你看到的将会是由Jar自动添加的cookie
	log.Println(cookies)
	return http.DefaultTransport.RoundTrip(req)
}
