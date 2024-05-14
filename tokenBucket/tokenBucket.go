package tokenbucket

import (
	"log"
	"sync"
	"time"
)

// If no token is available, Allow returns false.
type Limiter struct {
	// id
	mu    sync.Mutex
	limit float64
	burst int
	// bucket size
	tokens float64
	// last time the limiter's tokens field was updated
	last time.Time
}

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

func (lim *Limiter) Burst() int {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.burst
}

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
