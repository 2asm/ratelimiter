package tokenbucket 

import (
	"log"
	"sync"
	"time"
)

// Limit defines the maximum frequency of some events.
// Limit is represented as number of events per second.
// A zero Limit allows no events.
//type Limit float64

// A Limiter controls how frequently events are allowed to happen.
// It implements a "token bucket" of size b, initially full and refilled
// at rate r tokens per second.
// Informally, in any large enough time interval, the Limiter limits the
// rate to r tokens per second, with a maximum burst size of b events.
// See https://en.wikipedia.org/wiki/Token_bucket for more about token buckets.
//
// The zero value is a valid Limiter, but it will reject all events.
// Use NewLimiter to create non-zero Limiters.
//
// They differ in their behavior when no token is available.
// If no token is available, Allow returns false.
type Limiter struct {
    // id 
	mu     sync.Mutex
	limit  float64
	burst  int
	tokens float64
	// last is the last time the limiter's tokens field was updated
	last time.Time
}

// NewLimiter returns a new Limiter that allows events up to rate r and permits
// bursts of at most b tokens.
func NewLimiter(r float64, b int) *Limiter {
	return &Limiter{
		limit: r,
		burst: b,
	}
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

// advance calculates and returns an updated state for lim resulting from the passage of time.
// lim is not changed.
// advance requires that lim.mu is held.
func (lim *Limiter) advance(t time.Time) float64 {
	last := lim.last
    // time skew
	if t.Before(last) {
		last = t
	}

	// Calculate the new number of tokens, due to time that passed.
	elapsed := t.Sub(last)
	delta := elapsed.Seconds() * float64(lim.limit)
	tokens := lim.tokens + delta

	burst := float64(lim.burst)
	if tokens > burst {
		tokens = burst
	}
	return tokens
}

// Allow reports whether an event may happen now.
func (lim *Limiter) Allow() bool {
	return lim.Allow1(time.Now())
}

// helper function for Allow
func (lim *Limiter) Allow1(t time.Time) bool {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	tokens := lim.advance(t)

    // consume one token
	tokens -= 1

	if tokens < 0 {
		lim.tokens = 0
		return false
	} else {
		lim.last = t
		lim.tokens = tokens
		return true
	}
}

func main() {
	limiter := NewLimiter(1, 2)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Millisecond * 1999)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
}
