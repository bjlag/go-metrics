package limiter

// RateLimiter ограничивает число одновременных запросов. Основан на буферизованном канале.
type RateLimiter struct {
	ch chan struct{}
}

// NewRateLimiter создает rate limiter.
func NewRateLimiter(reqMax int) *RateLimiter {
	return &RateLimiter{
		ch: make(chan struct{}, reqMax),
	}
}

// Acquire занимает ресурс.
func (l *RateLimiter) Acquire() {
	l.ch <- struct{}{}
}

// Release освобождает ресурс.
func (l *RateLimiter) Release() {
	<-l.ch
}
