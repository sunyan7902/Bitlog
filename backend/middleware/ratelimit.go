package middleware

import (
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// IPRecord 登录记录结构
type IPRecord struct {
	Count     int
	Timestamp time.Time
}

var (
	rateLimitMap = make(map[string]*IPRecord)
	mu           sync.Mutex
)

const (
	rateLimitWindow = 5 * time.Minute
	maxAttempts     = 5
)

// GetClientIP 获取可靠的真实客户端 IP
func GetClientIP(c *gin.Context) string {
	forwarded := c.GetHeader("X-Forwarded-For")
	if forwarded != "" {
		parts := strings.Split(forwarded, ",")
		return strings.TrimSpace(parts[0])
	}
	realIP := c.GetHeader("X-Real-IP")
	if realIP != "" {
		return realIP
	}
	return c.ClientIP()
}

// ClearRateLimit 登录成功时清除对应 IP 计数
func ClearRateLimit(ip string) {
	mu.Lock()
	defer mu.Unlock()
	delete(rateLimitMap, ip)
}

// IPRateLimiter 防爆破限流拦截器
func IPRateLimiter() gin.HandlerFunc {
	return func(c *gin.Context) {
		ip := GetClientIP(c)
		now := time.Now()

		mu.Lock()
		record, exists := rateLimitMap[ip]

		// 若记录已超过滑动窗口，则清除重置
		if exists && now.Sub(record.Timestamp) > rateLimitWindow {
			delete(rateLimitMap, ip)
			exists = false
		}

		if exists && record.Count >= maxAttempts {
			mu.Unlock()
			c.JSON(http.StatusTooManyRequests, gin.H{
				"error": "Too many login attempts, please try again later",
			})
			c.Abort()
			return
		}

		// 记录新的登录尝试
		if !exists {
			rateLimitMap[ip] = &IPRecord{
				Count:     1,
				Timestamp: now,
			}
		} else {
			record.Count++
		}
		mu.Unlock()

		c.Next()
	}
}
