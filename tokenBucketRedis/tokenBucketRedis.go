package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// Redis transactions use optimistic locking.
const maxRetries = 1

type Limiter struct {
	id         int64 // client identification useid, ip
	cli        *redis.Client
	limit      float64
	burst      int
	last_key   string
	tokens_key string
}

func NewLimiter(user_id int64, r float64, b int, client *redis.Client) *Limiter {
	last_key := fmt.Sprintf("%v_tokenbucket_redis_last", user_id)
	tokens_key := fmt.Sprintf("%v_tokenbucket_redis_tokens", user_id)
	l := &Limiter{
		id:    user_id,
		cli:   client,
		limit: r,
		burst: b,
        last_key: last_key,
        tokens_key: tokens_key,
	}
	var ctx = context.Background()
	l.cli.Set(ctx, last_key, 0, 0)
	l.cli.Set(ctx, tokens_key, float64(b), 0)

	return l
}

// Limit returns the maximum overall event rate.
func (lim *Limiter) Limit() float64 {
	return lim.limit
}

func (lim *Limiter) Burst() int {
	return lim.burst
}

// Allow reports whether an event may happen now.
func (lim *Limiter) Allow() bool {
	return lim.Allow1(time.Now().UnixNano())
}

// helper function for Allow
func (lim *Limiter) Allow1(t int64) bool {

	var ctx = context.Background()

	txf := func(tx *redis.Tx) error {
		// Get the current value or zero.
		last, err := tx.Get(ctx, lim.last_key).Int64()
		if err != nil && err != redis.Nil {
			return err
		}
		if t < last {
			last = t
		}

		lim_tokens, err := tx.Get(ctx, lim.tokens_key).Float64()
		if err != nil && err != redis.Nil {
			return err
		}

		// Calculate the new number of tokens, due to time that passed.
		elapsed := t - last
		delta := float64(elapsed) * float64(lim.limit) / 1e9
		tokens := lim_tokens + delta

		burst := float64(lim.burst)
		if tokens > burst {
			tokens = burst
		}

		// Actual operation (local in optimistic lock).
		tokens -= 1
		if tokens < 0 {
			return errors.New("Rate limit Exceeded")
		}

		// Operation is commited only if the watched keys remain unchanged.
		_, err = tx.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
			pipe.Set(ctx, lim.last_key, t, 0)
			pipe.Set(ctx, lim.tokens_key, tokens, 0)
			return nil
		})
		return err
	}

	for i := 0; i < maxRetries; i++ {
		err := lim.cli.Watch(ctx, txf, lim.last_key, lim.tokens_key)
		if err == nil {
			// Success.
			return true
		}
		if err == redis.TxFailedErr {
			// Optimistic lock lost. Retry.
			continue
		}
		return false
	}
	return false
}

func main() {
	redis_client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
		Protocol: 3,  // specify 2 for RESP 2 or 3 for RESP 3
	})
	_, err := redis_client.Ping(context.Background()).Result()
	if err != nil {
		log.Fatalln("Couldn't connect to redis.")
	}
	limiter := NewLimiter(3213, 2, 6, redis_client)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 999)
	log.Printf("Waiting %v milisecond\n", 999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 999)
	log.Printf("Waiting %v milisecond\n", 999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 1999)
	log.Printf("Waiting %v milisecond\n", 1999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
}
