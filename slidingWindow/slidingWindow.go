package slidingwindow

import (
	"log"
	"sync"
	"time"
)

type Limiter struct {
	mu     sync.Mutex
    // Events allowed per window
	limit  float64
    // tokens keep count of how many events are left in the window
	tokens float64
    // Window size
	size   time.Duration 
	// firstInWindow is the firstInWindow time the limiter's tokens field was updated 
    // while the request was in the same window
    firstInWindow time.Time
}

func NewLimiter(r float64, d time.Duration) *Limiter {
	return &Limiter{
		limit: r,
        size: d,
	}
}

func (lim *Limiter) Limit() float64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.limit
}

// advance require mutex to be held
func (lim *Limiter) advance(t time.Time) float64 {
	firstInWindow := lim.firstInWindow
	if t.Before(firstInWindow) {
		firstInWindow = t
	}

	// Calculate the new number of tokens, due to time that passed.
	elapsed := t.Sub(firstInWindow)
    
    tokens := lim.tokens
	if elapsed > lim.size {
		tokens = lim.limit 
    }
	return tokens
}

func (lim *Limiter) Allow() bool {
	return lim.Allow1(time.Now())
}

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
        elapsed := t.Sub(lim.firstInWindow)
        if elapsed > lim.size {
            lim.firstInWindow = t;
        }
		lim.tokens = tokens
		return true
	}
}

func main() {
	limiter := NewLimiter(3, time.Second*1)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Second)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Second)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	time.Sleep(time.Second)
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
	log.Printf("%v\n", limiter.Allow())
}
