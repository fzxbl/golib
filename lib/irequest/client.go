package irequest

import (
	"errors"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
	"time"

	"github.com/fzxbl/golib/lib/ilimit"
	"golang.org/x/net/publicsuffix"
)

type Client struct {
	client       *http.Client
	limit        bool
	limiter      ilimit.Limiter //所有请求都会限流
	isBlockLimit bool
}

type InitCookie struct {
	Host         string //以http或https开头的域名
	Domain       string //cookie要绑定的域名,支持.example.com类型的通配域名
	CookieHeader string //从浏览器复制出来的cookie
}

type ClientOptions struct {
	UseCookie    bool           //是否开启Cookie支持
	InitCookie   *InitCookie    //是否使用已有Cookie初始化，默认为空
	Limiter      ilimit.Limiter //限流，默认为空
	IsBlockLimit bool           // 拿不到令牌时是否阻塞，直到拿到令牌
	Timeout      time.Duration
	Transport    http.RoundTripper
}

type ClientOption func(opts *ClientOptions)

func (i InitCookie) initCookie() (jar *cookiejar.Jar) {
	jar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
	if err != nil {
		panic(err)
	}
	urlObj, err := url.Parse(i.Host)
	if err != nil {
		panic(err)
	}
	cookieParts := strings.Split(i.CookieHeader, ";")
	cookies := make([]*http.Cookie, 0)
	for _, part := range cookieParts {
		if part == "" {
			continue
		}
		kv := strings.SplitN(part, "=", 2)
		if len(kv) != 2 {
			log.Printf("Invalid cookie part: %s", part)
			continue
		}
		name := strings.TrimSpace(kv[0])
		value := strings.TrimSpace(kv[1])
		value = strings.ReplaceAll(value, `"`, ``)
		cookie := &http.Cookie{}
		cookie.Value = value
		cookie.Name = name
		cookie.Domain = i.Domain
		cookies = append(cookies, cookie)
	}

	jar.SetCookies(urlObj, cookies)
	return
}
func NewClient(options ...ClientOption) (c Client) {
	c.client = &http.Client{}
	var opts ClientOptions
	for _, option := range options {
		option(&opts)
	}
	if opts.Limiter != nil {
		c.limit = true
		c.limiter = opts.Limiter
		c.isBlockLimit = opts.IsBlockLimit
	}
	if opts.Transport != nil {
		c.client.Transport = opts.Transport
	}

	if opts.UseCookie {
		if opts.InitCookie != nil {
			cookieJar := opts.InitCookie.initCookie()
			c.client.Jar = cookieJar
		} else {
			cookieJar, err := cookiejar.New(&cookiejar.Options{PublicSuffixList: publicsuffix.List})
			if err != nil {
				panic(err)
			}
			c.client.Jar = cookieJar
		}
	}

	if opts.Timeout > 0 {
		c.client.Timeout = opts.Timeout
	}

	return
}

// WithCookie 开启Cookie支持,并使用initCookie初始化,initCookie为nil时，默认使用空的CookieJar
func WithCookie(initCookie *InitCookie) ClientOption {
	return func(opts *ClientOptions) {
		opts.UseCookie = true
		opts.InitCookie = initCookie
	}
}

// WithClientLimiter 限流，使用该Client的所有请求都会限流
func WithClientLimiter(limiter ilimit.Limiter, isBlockRequest bool) ClientOption {
	return func(opts *ClientOptions) {
		if limiter != nil {
			opts.Limiter = limiter
		} else {
			panic(errors.New(`limiter can not be nil`))
		}
		opts.IsBlockLimit = isBlockRequest
	}
}

func WithTimeout(timeout time.Duration) ClientOption {
	return func(opts *ClientOptions) {
		opts.Timeout = timeout
	}
}

func WithCustomTransport(transport http.RoundTripper) ClientOption {
	return func(opts *ClientOptions) {
		if transport == nil {
			panic(errors.New(`transport can not be nil`))
		}
		opts.Transport = transport
	}
}
