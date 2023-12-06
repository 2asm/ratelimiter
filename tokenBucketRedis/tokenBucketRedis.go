package main

import (
	"context"
	"fmt"
	"log"
	"strconv"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

type Limiter struct {
	id     int64 // client identification useid, ip
	mu     sync.Mutex
	cli    *redis.Client
	limit  float64
	burst  int
	// tokens float64 // stored in redis as id_tokenbucket_redis_tokens
	// last is the last time the limiter's tokens field was updated (Unix Nano)
	// last int64 // stored in redis as id_tokenbucket_redis_last
}

// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
func NewLimiter(user_id int64, r float64, b int, client *redis.Client) *Limiter {
    l := &Limiter{
		id:    user_id,
		cli:   client,
		limit: r,
		burst: b,
	}
    l.set_redis_last(0)
    l.set_redis_tokens(float64(b))
    return l
}

// Limit returns the maximum overall event rate.
func (lim *Limiter) Limit() float64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}

// Burst returns the maximum burst size. Burst is the maximum number of tokens
// that can be consumed in a single call to Allow, Reserve, or Wait, so higher
// Burst values allow more events to happen at once.
// A zero Burst allows no events, unless limit == Inf.
func (lim *Limiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}
// Allow reports whether an event may happen now.
func (lim *Limiter) Allow() bool {
	return lim.Allow1(time.Now().UnixNano())
}

// helper function for Allow
func (lim *Limiter) Allow1(t int64) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	tokens := lim.advance(t)

	// consume one token
	tokens -= 1

	if tokens < 0 {
        lim.set_redis_tokens(0)
		return false
	} else {
        lim.set_redis_last(t)
        lim.set_redis_tokens(tokens)
		return true
	}
}

func (lim *Limiter) get_redis_last() int64 {
    var ctx = context.Background()
    last_key := fmt.Sprintf("%v_tokenbucket_redis_last", lim.id)
    last_string, err := lim.cli.Get(ctx, last_key).Result()
    if err != nil {
        log.Fatalf("redis query failed\n")
    }
    last, _ := strconv.ParseInt(last_string, 10, 64)
    return last
}

func (lim * Limiter) get_redis_tokens() float64 {
    var ctx = context.Background()
    tokens_key := fmt.Sprintf("%v_tokenbucket_redis_tokens", lim.id)
    tokens_string, err := lim.cli.Get(ctx, tokens_key).Result()
    if err != nil {
        log.Fatalf("redis query failed\n")
    }
    lim_tokens, _ := strconv.ParseFloat(tokens_string, 64)
    return lim_tokens
}

func (lim *Limiter) set_redis_last(x int64) {
	var ctx = context.Background()
	last_key := fmt.Sprintf("%v_tokenbucket_redis_last", lim.id)
    _, err := lim.cli.Set(ctx, last_key, x, 0).Result()
    if err != nil {
        log.Fatalf("redis query failed\n")
    }
}

func (lim * Limiter) set_redis_tokens(x float64) {
	var ctx = context.Background()
	tokens_key := fmt.Sprintf("%v_tokenbucket_redis_tokens", lim.id)
    _, err := lim.cli.Set(ctx, tokens_key, x, 0).Result()
    if err != nil {
        log.Fatalf("redis query failed\n")
    }
}

// advance calculates and returns an updated state for lim resulting from the passage of time.
// lim is not changed.
// advance requires that lim.mu is held.
func (lim *Limiter) advance(t int64) float64 {
    last := lim.get_redis_last()
    lim_tokens := lim.get_redis_tokens()
    if t < last {
        last = t
    }

    // Calculate the new number of tokens, due to time that passed.
    elapsed := t - last
    delta := float64(elapsed) * float64(lim.limit) / 1e9
    tokens := lim_tokens + delta

    burst := float64(lim.burst)
    if tokens > burst {
        tokens = burst
    }
    return tokens
}

func main() {
	redis_client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})
	limiter := NewLimiter(3213, 2, 6, redis_client)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 1999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
}
