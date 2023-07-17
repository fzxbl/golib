package ilimit

import (
	"context"
	"time"

	"golang.org/x/time/rate"
)

// Limiter 限速器
type Limiter interface {
	// Allow 是否能获取1个令牌，非阻塞
	Allow() bool
	// WaitCtx 阻塞直到获取1个令牌，当ctx被取消、超时时将返回 error
	Wait(ctx context.Context) (err error)
}

type limiter struct {
	*rate.Limiter
}

// NewLimiter 创建一个容量为capacity的令牌桶，初始状态下桶是满的，并且以间隔interval向桶中增加1个新令牌
func NewLimiter(interval time.Duration, capacity int) Limiter {
	return &limiter{
		Limiter: rate.NewLimiter(rate.Every(interval), capacity),
	}
}
