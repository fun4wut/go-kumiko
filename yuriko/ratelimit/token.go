package ratelimit

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"math"
	"strconv"
	"sync"
	"time"
)

const (
	numKey      = "num"
	lastTimeKey = "last_time"
)

// TokenBucket 令牌桶算法，通过 redis 来实现 分布式
type TokenBucket struct {
	bucketKey string
	rate      float64 // 令牌生成速率
	capacity  int64   // 桶的容量
	rdb       *redis.Client
	mu        sync.Mutex
	once      sync.Once
}

func NewTokenBucket(bucketKey string, rate float64, capacity int64) *TokenBucket {
	client := redis.NewClient(&redis.Options{
		Addr: "0.0.0.0:6379",
	})
	return &TokenBucket{rdb: client, rate: rate, capacity: capacity, bucketKey: bucketKey}
}

func (tb *TokenBucket) Allow() bool {
	tb.mu.Lock()
	defer tb.mu.Unlock()
	ctx := context.Background()
	currentTime := time.Now().Unix()
	tb.once.Do(func() { // lazy init
		tb.rdb.HSet(ctx, tb.bucketKey, numKey, tb.capacity) // 初始状态，令牌桶是满的
		tb.rdb.HSet(ctx, tb.bucketKey, lastTimeKey, currentTime)
	})
	res, _ := tb.getFields(ctx, numKey, lastTimeKey)
	lastNum := res[0]
	lastTime := res[1]
	delta := tb.rate * float64(currentTime-lastTime)
	currentNum := int64(math.Ceil(min(float64(lastNum)+delta, float64(tb.capacity))))
	if currentNum <= 0 { // 令牌桶爆了，，不能通过
		return false
	}
	// 还有令牌可以获取
	if err := tb.rdb.HMSet(ctx, tb.bucketKey, map[string]any{
		numKey:      currentNum - 1,
		lastTimeKey: currentTime,
	}).Err(); err != nil {
		fmt.Println(err)
	}
	return true
}

func (tb *TokenBucket) getField(ctx context.Context, field string) (int64, error) {
	res, _ := tb.rdb.HGet(ctx, tb.bucketKey, field).Result()
	return strconv.ParseInt(res, 10, 64)
}

func (tb *TokenBucket) getFields(ctx context.Context, fields ...string) (arr []int64, err error) {
	res, err := tb.rdb.HMGet(ctx, tb.bucketKey, fields...).Result()
	if err != nil {
		return
	}
	for _, v := range res {
		vint, _ := strconv.ParseInt(v.(string), 10, 64)
		arr = append(arr, vint)
	}
	return

}

type Number interface {
	int | int64 | float64 | float32
}

func min[T Number](a, b T) T {
	if a < b {
		return a
	}
	return b
}
