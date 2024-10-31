package limiter

type RateLimiter struct {
	ch chan struct{}
}

func NewRateLimiter(reqMax int) *RateLimiter {
	return &RateLimiter{
		ch: make(chan struct{}, reqMax),
	}
}

func (l *RateLimiter) Acquire() {
	l.ch <- struct{}{}
}

func (l *RateLimiter) Release() {
	<-l.ch
}
