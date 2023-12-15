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

func (c Client) doWithRetry(req *http.Request, options []RequestOptionFunc) (resp IResponse, err error) {
	// 处置不同的选项
	var opts RequestOptions
	for _, opt := range options {
		opt(&opts)
	}
	for i := 0; i <= opts.RetryCount; i++ {
		resp, err = c.do(req, opts)
		if err != nil {
			return
		}
	}

	return
}

func (c Client) do(req *http.Request, opts RequestOptions) (resp response, err error) {
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

	var r response

	// 状态通用字段拷贝
	r.statusCode = originResp.StatusCode
	r.status = originResp.Status

	// 备份原始body
	var rawContent []byte
	if rawContent, err = io.ReadAll(originResp.Body); err != nil {
		return r, err
	}

	// 检查返回状态码
	if opts.ResponceOpt.CheckCode {
		if r.statusCode != opts.TargetCode {
			err = fmt.Errorf("status code not match. expected: %d actual: %d", opts.TargetCode, r.statusCode)
			r.content = string(rawContent)
			return r, err
		}
	}
	// 开始处理不同的返回类型
	if opts.ResponceOpt.HTTPResp {
		r.httpResp = *originResp
		r.httpResp.Body = io.NopCloser(bytes.NewBuffer(rawContent))
	}

	if opts.ResponceOpt.BytesResp {
		r.rawContent = rawContent
	}

	if opts.ResponceOpt.ReaderResp {
		r.body = bytes.NewReader(rawContent)
	}

	if opts.ResponceOpt.Unmarshal {
		if err = jsoniter.Unmarshal(rawContent, opts.ResponceOpt.UnmarshalContainer); err != nil {
			return r, err
		}
	}

	if opts.ResponceOpt.StringResp {
		r.content = string(rawContent)
	}
	return r, nil
}

func (c Client) Do(req *http.Request, options ...RequestOptionFunc) (resp IResponse, err error) {
	resp, err = c.doWithRetry(req, options)
	return
}
func (c Client) GetURL(URL string, options ...RequestOptionFunc) (resp IResponse, err error) {
	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return
	}
	resp, err = c.doWithRetry(req, options)
	return
}

func (c Client) PostURL(URL string, body io.Reader, options ...RequestOptionFunc) (resp IResponse, err error) {
	req, err := http.NewRequest(http.MethodPost, URL, body)
	if err != nil {
		return
	}
	resp, err = c.doWithRetry(req, options)
	return
}
