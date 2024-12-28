package middleware

import (
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/time/rate"
)

// IPRateLimiter IP限流器
type IPRateLimiter struct {
	ips    map[string]*rate.Limiter
	mu     *sync.RWMutex
	rate   rate.Limit
	burst  int
	ttl    time.Duration
	lastOp map[string]time.Time
}

// NewIPRateLimiter 创建一个新的IP限流器
func NewIPRateLimiter(r rate.Limit, burst int, ttl time.Duration) *IPRateLimiter {
	i := &IPRateLimiter{
		ips:    make(map[string]*rate.Limiter),
		mu:     &sync.RWMutex{},
		rate:   r,
		burst:  burst,
		ttl:    ttl,
		lastOp: make(map[string]time.Time),
	}

	// 启动清理过期限流器的goroutine
	go i.cleanupLoop()
	return i
}

// AddIP 创建一个新的限流器并添加到map中
func (i *IPRateLimiter) AddIP(ip string) *rate.Limiter {
	i.mu.Lock()
	defer i.mu.Unlock()

	limiter := rate.NewLimiter(i.rate, i.burst)
	i.ips[ip] = limiter
	i.lastOp[ip] = time.Now()

	return limiter
}

// GetLimiter 获取IP对应的限流器，如果不存在则创建一个新的
func (i *IPRateLimiter) GetLimiter(ip string) *rate.Limiter {
	i.mu.Lock()
	limiter, exists := i.ips[ip]

	if !exists {
		i.mu.Unlock()
		return i.AddIP(ip)
	}

	i.lastOp[ip] = time.Now()
	i.mu.Unlock()

	return limiter
}

// cleanupLoop 定期清理过期的限流器
func (i *IPRateLimiter) cleanupLoop() {
	ticker := time.NewTicker(i.ttl)
	for range ticker.C {
		i.mu.Lock()
		for ip, lastOp := range i.lastOp {
			if time.Since(lastOp) > i.ttl {
				delete(i.ips, ip)
				delete(i.lastOp, ip)
			}
		}
		i.mu.Unlock()
	}
}

// RateLimitMiddleware IP限流中间件
func RateLimitMiddleware(limiter *IPRateLimiter) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取IP地址
		ip := c.ClientIP()
		
		// 获取限流器
		l := limiter.GetLimiter(ip)
		if !l.Allow() {
			c.JSON(http.StatusTooManyRequests, gin.H{
				"code": 429,
				"msg":  "请求过于频繁，请稍后再试",
			})
			c.Abort()
			return
		}
		
		c.Next()
	}
}
