package ilimit

import (
	"context"
	"testing"
	"time"
)

func Test_Limiter(t *testing.T) {
	ctx := context.Background()
	l := NewLimiter(time.Second, 1)
	for i := 0; i < 5; i++ {
		if err := l.Wait(ctx); err != nil {
			t.Fatal(err)
		} else {
			t.Logf("%s", time.Now().Format(time.DateTime))
		}
	}
	time.Sleep(time.Second)
	for i := 0; i < 5; i++ {
		if l.Allow() {
			t.Logf("GET 1 %s", time.Now().Format(time.DateTime))
		} else {
			t.Logf("FAIL 1 %s", time.Now().Format(time.DateTime))
		}
	}
}
