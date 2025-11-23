package middleware

import (
	"context"
	"flight-aggregator/internal/config"
	"flight-aggregator/internal/models"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/echo/v4"
)

type RedisSlidingWindow struct {
	client *redis.Client
	limit  int
	window time.Duration
}

func NewRedisSlidingWindowRateLimit() echo.MiddlewareFunc {
	cfg := config.MustLoad()
	
	// Parse Redis URL if it contains credentials
	var rdb *redis.Client
	if len(cfg.RedisAddr) > 8 && cfg.RedisAddr[:8] == "redis://" {
		opt, err := redis.ParseURL(cfg.RedisAddr)
		if err != nil {
			log.Printf("Failed to parse Redis URL: %v, falling back to simple connection", err)
			rdb = redis.NewClient(&redis.Options{
				Addr: "localhost:6379",
				DB:   0,
			})
		} else {
			rdb = redis.NewClient(opt)
		}
	} else {
		rdb = redis.NewClient(&redis.Options{
			Addr: cfg.RedisAddr,
			DB:   0,
		})
	}

	rsw := &RedisSlidingWindow{
		client: rdb,
		limit:  cfg.RateLimitCount,
		window: cfg.RateLimitWindow,
	}

	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			// Check Rate Limit Using IP
			clientIP := c.RealIP()
			if clientIP == "" {
				clientIP = c.Request().RemoteAddr
			}

			allowed, err := rsw.Allow(c.Request().Context(), clientIP)
			if err != nil {
				// If Redis fails, allow the request (fail open)
				return next(c)
			}

			if !allowed {
				errorResp := models.ErrorResponse{
					Status:  "error",
					Code:    "RATE_LIMIT_EXCEEDED",
					Message: "Too many requests, please try again later",
				}
				return c.JSON(http.StatusTooManyRequests, errorResp)
			}

			return next(c)
		}
	}
}

func (rsw *RedisSlidingWindow) Allow(ctx context.Context, key string) (bool, error) {
	now := time.Now().UnixMilli()
	windowStart := now - rsw.window.Milliseconds()
	redisKey := fmt.Sprintf("rate_limit:%s", key)

	// Remove old entries and count current requests
	pipe := rsw.client.Pipeline()
	pipe.ZRemRangeByScore(ctx, redisKey, "0", strconv.FormatInt(windowStart, 10))
	pipe.ZCard(ctx, redisKey)
	pipe.ZAdd(ctx, redisKey, &redis.Z{Score: float64(now), Member: now})
	pipe.Expire(ctx, redisKey, rsw.window)

	results, err := pipe.Exec(ctx)
	if err != nil {
		return false, err
	}

	count := results[1].(*redis.IntCmd).Val()
	return count < int64(rsw.limit), nil
}



