package irequest

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	jsoniter "github.com/json-iterator/go"
)

type Response struct {
	Status     string
	StatusCode int
	RawContent []byte
	Content    string
	HTTPResp   http.Response
	Body       io.ReadSeeker
}

func (c Client) do(req *http.Request, options []RequestOptionFunc) (resp Response, err error) {
	// 处置不同的选项
	var opts RequestOptions
	for _, opt := range options {
		opt(&opts)
	}
	for i := 0; i <= opts.RetryCount; i++ {
		resp, err = c.do1(req, opts)
		if err != nil {
			return
		}
	}

	return
}

func (c Client) do1(req *http.Request, opts RequestOptions) (resp Response, err error) {
	// client限流器
	var originResp *http.Response
	if c.limit {
		if c.isBlockLimit {
			if err = c.limiter.Wait(context.Background()); err != nil {
				return
			}
		} else {
			if ok := c.limiter.Allow(); !ok {
				err = errors.New(`limit by client limiter`)
				return
			}
		}
	}

	// 请求限流器
	if opts.AboutLimit.Limit {
		if opts.AboutLimit.IsBlockLimit {
			if err = opts.AboutLimit.Limiter.Wait(context.Background()); err != nil {
				return
			}
		} else {
			if ok := opts.AboutLimit.Limiter.Allow(); !ok {
				err = errors.New(`limit by request limiter`)
				return
			}
		}
	}

	var ctx context.Context
	if opts.Timeout != 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(context.Background(), opts.Timeout)
		defer cancel() //及时释放ctx资源、断开连接
		req = req.WithContext(ctx)
	}

	// 发起请求
	if originResp, err = c.client.Do(req); err != nil {
		return
	}
	defer originResp.Body.Close()

	// 状态通用字段拷贝
	resp.StatusCode = originResp.StatusCode
	resp.Status = originResp.Status

	// 备份原始body
	var rawContent []byte
	if rawContent, err = io.ReadAll(originResp.Body); err != nil {
		return
	}

	// 检查返回状态码
	if opts.AboutResponce.CheckCode {
		if resp.StatusCode != opts.TargetCode {
			err = fmt.Errorf("status code not match. expected: %d actual: %d", opts.TargetCode, resp.StatusCode)
			resp.Content = string(rawContent)
			return
		}
	}
	// 开始处理不同的返回类型
	if opts.AboutResponce.HTTPResp {
		resp.HTTPResp = *originResp
		resp.HTTPResp.Body = io.NopCloser(bytes.NewBuffer(rawContent))
	}

	// 开始处理不同的返回类型
	if opts.AboutResponce.HTTPResp {
		resp.HTTPResp = *originResp
		resp.HTTPResp.Body = io.NopCloser(bytes.NewBuffer(rawContent))
	}

	if opts.AboutResponce.BytesResp {
		resp.RawContent = rawContent
	}

	if opts.AboutResponce.ReaderResp {
		resp.Body = bytes.NewReader(rawContent)
	}

	if opts.AboutResponce.Unmarshal {
		if err = jsoniter.Unmarshal(rawContent, opts.AboutResponce.UnmarshalContainer); err != nil {
			return
		}
	}

	if opts.AboutResponce.StringResp {
		resp.Content = string(rawContent)
	}
	return
}

func (c Client) Do(req *http.Request, options ...RequestOptionFunc) (resp Response, err error) {
	resp, err = c.do(req, options)
	return
}
func (c Client) GetURL(URL string, options ...RequestOptionFunc) (resp Response, err error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return
	}
	resp, err = c.do(req, options)
	return
}

func (c Client) PostURL(URL string, body io.Reader, options ...RequestOptionFunc) (resp Response, err error) {
	req, err := http.NewRequest(http.MethodPost, URL, body)
	if err != nil {
		return
	}
	resp, err = c.do(req, options)
	return
}
