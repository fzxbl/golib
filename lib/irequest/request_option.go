package irequest

import (
	"errors"

	"github.com/fzxbl/golib/lib/ilimit"
)

type AboutResponce struct {
	// 检查返回状态码，不为200时，返回错误信息
	CheckCode bool
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
}
type RequestOption func(opts *RequestOptions)

// WithUnmarshalResp 将返回结果解析为对象
func WithUnmarshalResp(container interface{}) RequestOption {
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
func WithHTTPResp() RequestOption {
	return func(opts *RequestOptions) {
		opts.HTTPResp = true
	}
}

// WithReaderResp 返回一个Reader对象，放入Response.Body
func WithReaderResp() RequestOption {
	return func(opts *RequestOptions) {
		opts.ReaderResp = true
	}
}

// WithBytesResp 返回一个[]byte对象，放入Response.RawContent
func WithByteResp() RequestOption {
	return func(opts *RequestOptions) {
		opts.BytesResp = true
	}
}

// WithStringResp 返回一个string对象，放入Response.Content
func WithStringResp() RequestOption {
	return func(opts *RequestOptions) {
		opts.StringResp = true
	}
}

// WithLimiter 请求是否使用限流器
func WithLimiter(limiter ilimit.Limiter, isBlockRequest bool) RequestOption {
	return func(opts *RequestOptions) {
		opts.Limit = true
		opts.IsBlockLimit = isBlockRequest
		opts.Limiter = limiter
	}
}
