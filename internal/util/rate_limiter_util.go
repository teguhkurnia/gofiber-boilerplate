package util

import (
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/redis/go-redis/v9"
)

type RateLimiterUtil struct {
	Redis      *redis.Client
	MaxRequest int64
	Duration   time.Duration
}

func NewRateLimiterUtil(
	redisClient *redis.Client,
	maxRequest int64,
	duration time.Duration,
) *RateLimiterUtil {
	return &RateLimiterUtil{
		Redis:      redisClient,
		MaxRequest: maxRequest,
		Duration:   duration,
	}
}

func (r *RateLimiterUtil) IsAllowed(ctx *fiber.Ctx, key string, localMaxRequest *int64) error {
	increment, err := r.Redis.Incr(ctx.UserContext(), key).Result()
	if err != nil {
		return fiber.NewError(fiber.StatusInternalServerError, "Failed to check rate limit")
	}

	if increment == 1 {
		r.Redis.Expire(ctx.UserContext(), key, r.Duration)
	}

	isAllowed := localMaxRequest != nil && increment <= *localMaxRequest || localMaxRequest == nil && increment <= r.MaxRequest
	if !isAllowed {
		// get remaining time
		remaining, err := r.Redis.TTL(ctx.UserContext(), key).Result()
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to check rate limit")
		}

		ctx.Set("Retry-After", strconv.Itoa(int(remaining.Seconds())))
		ctx.Set("X-RateLimit-Limit", strconv.Itoa(int(r.MaxRequest)))
		ctx.Set("X-RateLimit-Remaining", strconv.Itoa(max(0, int(r.MaxRequest-increment))))
		return fiber.NewError(fiber.StatusTooManyRequests, "Rate limit exceeded")
	}

	ctx.Set("X-RateLimit-Limit", strconv.Itoa(int(r.MaxRequest)))
	ctx.Set("X-RateLimit-Remaining", strconv.Itoa(int(r.MaxRequest-increment)))

	return nil
}
