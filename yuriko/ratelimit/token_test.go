package ratelimit

import (
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestTokenBucket_Allow(t *testing.T) {
	bucket := NewTokenBucket("lolll", 5, 10)
	var wg sync.WaitGroup
	for i := 0; i < 1000; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(time.Duration(rand.Int63n(2000)) * time.Millisecond)
			t.Log(bucket.Allow())
			wg.Done()
		}()
	}
	wg.Wait()
}
