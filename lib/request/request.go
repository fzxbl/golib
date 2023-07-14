package irequest

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
)

// GetURLWithUnmarshal GET请求，将返回内容存入提供的container
func GetURLWithUnmarshal(url string, timeout time.Duration, resultContainer interface{}) (err error) {
	client := http.Client{Timeout: timeout}
	var resp *http.Response
	if resp, err = client.Get(url); err != nil {
		return
	}
	defer resp.Body.Close()
	var ret []byte
	if ret, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status is %d,content is %s", resp.StatusCode, string(ret))
		return
	}
	err = jsoniter.Unmarshal(ret, resultContainer)
	return
}

// GetURL GET请求，将原结果返回
func GetURLRaw(url string, timeout time.Duration) (ret []byte, err error) {
	client := &http.Client{Timeout: timeout}
	var resp *http.Response
	resp, err = client.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if ret, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status is %d,content is %s", resp.StatusCode, string(ret))
	}
	return
}

// PostURLWithUnmarshal POST请求，将返回内容存入提供的container
func PostURLWithUnmarshal(url, contentType string, body io.Reader, timeout time.Duration, resultContainer interface{}) (err error) {
	client := http.Client{Timeout: timeout}
	var resp *http.Response
	if resp, err = client.Post(url, contentType, body); err != nil {
		return
	}
	defer resp.Body.Close()
	var ret []byte
	if ret, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status is %d,content is %s", resp.StatusCode, string(ret))
	}
	err = jsoniter.Unmarshal(ret, resultContainer)
	return
}

// PostURL POST请求，将原结果返回
func PostURLRaw(url, contentType string, body io.Reader, timeout time.Duration) (ret []byte, err error) {
	client := &http.Client{Timeout: timeout}
	var resp *http.Response
	resp, err = client.Post(url, contentType, body)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if ret, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status is %d,content is %s", resp.StatusCode, string(ret))
	}
	return
}

// HTTPDo 完成请求，将原结果返回
func HTTPDoRaw(req *http.Request, timeout time.Duration) (ret []byte, err error) {
	client := &http.Client{Timeout: timeout}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	if ret, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status is %d,content is %s", resp.StatusCode, string(ret))
	}
	return
}

// HTTPDoWithUnmarshal 完成请求，将返回内容存入提供的container
func HTTPDoWithUnmarshal(url, contentType string, body io.Reader, timeout time.Duration, resultContainer interface{}) (err error) {
	client := http.Client{Timeout: timeout}
	var resp *http.Response
	if resp, err = client.Post(url, contentType, body); err != nil {
		return
	}
	defer resp.Body.Close()
	var ret []byte
	if ret, err = ioutil.ReadAll(resp.Body); err != nil {
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("http status is %d,content is %s", resp.StatusCode, string(ret))
	}
	err = jsoniter.Unmarshal(ret, resultContainer)
	return
}

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
