package middlewares

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

type RateLimiter struct {
	sync.RWMutex
	limits map[string][]time.Time
	window time.Duration
	limit  int
	redis  *redis.Client
}

func NewRateLimiter(window time.Duration, limit int, redisClient *redis.Client) *RateLimiter {
	return &RateLimiter{
		limits: make(map[string][]time.Time),
		window: window,
		limit:  limit,
		redis:  redisClient,
	}
}

func RateLimitMiddleware(window time.Duration, limit int, redisClient *redis.Client) gin.HandlerFunc {
	limiter := NewRateLimiter(window, limit, redisClient)

	return func(c *gin.Context) {
		ip := c.ClientIP()

		if limiter.redis != nil {
			// Usando Redis para rate limiting
			allowed, remaining, err := limiter.checkRedis(c, ip)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "rate limit check failed"})
				c.Abort()
				return
			}

			if !allowed {
				c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
				c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
				c.Abort()
				return
			}
		} else {
			// Usando memória para rate limiting
			allowed, remaining := limiter.checkMemory(ip)
			if !allowed {
				c.Header("X-RateLimit-Remaining", fmt.Sprintf("%d", remaining))
				c.Header("X-RateLimit-Limit", fmt.Sprintf("%d", limit))
				c.JSON(http.StatusTooManyRequests, gin.H{"error": "rate limit exceeded"})
				c.Abort()
				return
			}
		}

		c.Next()
	}
}

func (rl *RateLimiter) checkMemory(ip string) (bool, int) {
	rl.Lock()
	defer rl.Unlock()

	now := time.Now()
	window := now.Add(-rl.window)

	// Remove timestamps antigos
	if times, exists := rl.limits[ip]; exists {
		var valid []time.Time
		for _, t := range times {
			if t.After(window) {
				valid = append(valid, t)
			}
		}
		rl.limits[ip] = valid
	}

	// Verifica o limite
	times := rl.limits[ip]
	if len(times) >= rl.limit {
		return false, 0
	}

	// Adiciona novo timestamp
	rl.limits[ip] = append(times, now)

	return true, rl.limit - len(rl.limits[ip])
}

func (rl *RateLimiter) checkRedis(c *gin.Context, ip string) (bool, int, error) {
	ctx := c.Request.Context()
	key := fmt.Sprintf("ratelimit:%s", ip)

	// Adiciona timestamp atual e remove timestamps antigos
	now := time.Now().UnixNano()
	member := redis.Z{
		Score:  float64(now),
		Member: now,
	}

	pipe := rl.redis.Pipeline()
	pipe.ZAdd(ctx, key, member)
	pipe.ZRemRangeByScore(ctx, key, "0", fmt.Sprintf("%d", now-rl.window.Nanoseconds()))

	// Conta requisições no período
	countCmd := pipe.ZCard(ctx, key)
	pipe.Expire(ctx, key, rl.window)

	_, err := pipe.Exec(ctx)
	if err != nil {
		return false, 0, err
	}

	count := countCmd.Val()
	remaining := rl.limit - int(count)

	return remaining >= 0, remaining, nil
}
