package irequest

import (
	"bytes"
	"io"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

// MakeParams GET请求,生成url参数
func MakeParams(pmap map[string]string) (pstr string) {
	p := url.Values{}
	for k, v := range pmap {
		if v != "" {
			p.Add(k, v)
		}
	}
	pstr = strings.Replace(p.Encode(), "+", "%20", -1)
	return
}

func MakeBody(param interface{}) (body io.Reader, err error) {
	var data []byte
	if data, err = jsoniter.Marshal(param); err != nil {
		return
	} else {
		body = bytes.NewBuffer(data)
	}
	return
}
