package irequest

import (
	"errors"
	"time"

	"github.com/fzxbl/golib/lib/ilimit"
)

type AboutResponce struct {
	// 检查返回状态码，不为200时，返回错误信息
	CheckCode  bool
	TargetCode int
	// 是否将返回解析为对象
	Unmarshal bool
	// unmarshal的目标对象
	UnmarshalContainer interface{}
	// 返回http.Response对象,但是不需要再调用Close方法
	HTTPResp bool
	// 返回一个Reader对象，放入Response.Body
	ReaderResp bool
	// 返回[]byte对象
	BytesResp bool
	// 返回string对象
	StringResp bool
}
type AboutLimit struct {
	Limit   bool
	Limiter ilimit.Limiter
	// 拿不到令牌时是否阻塞
	IsBlockLimit bool
}

type RequestOptions struct {
	AboutResponce
	AboutLimit
	RetryCount int
	Timeout    time.Duration
}
type RequestOptionFunc func(opts *RequestOptions)

// WithUnmarshalResp 将返回结果解析为对象
func WithUnmarshalResp(container interface{}) RequestOptionFunc {
	return func(opts *RequestOptions) {
		if container != nil {
			opts.Unmarshal = true
			opts.UnmarshalContainer = container
		} else {
			panic(errors.New("container must be a valid container pointer"))
		}

	}
}

// WithHTTPResp 返回http.Response，但是不需要自己调用resp.Body.Close()
func WithHTTPResp() RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.HTTPResp = true
	}
}

// WithReaderResp 返回一个Reader对象，放入Response.Body
func WithReaderResp() RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.ReaderResp = true
	}
}

// WithBytesResp 返回一个[]byte对象，放入Response.RawContent
func WithByteResp() RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.BytesResp = true
	}
}

// WithStringResp 返回一个string对象，放入Response.Content
func WithStringResp() RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.StringResp = true
	}
}

// WithCheckCode 返回时立刻检查状态码，不为code时，返回错误
func WithCheckCode(code int) RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.CheckCode = true
		opts.TargetCode = code
	}
}

// WithLimiter 请求是否使用限流器
func WithLimiter(limiter ilimit.Limiter, isBlockRequest bool) RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.Limit = true
		opts.IsBlockLimit = isBlockRequest
		opts.Limiter = limiter
	}
}

// WithRetry 请求重试次数
func WithRetry(retryCount int) RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.RetryCount = retryCount
	}
}

// WithRequestTimeout 请求重试次数
func WithRequestTimeout(timeout time.Duration) RequestOptionFunc {
	return func(opts *RequestOptions) {
		opts.Timeout = timeout
	}
}
