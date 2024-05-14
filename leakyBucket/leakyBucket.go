package leakyBucket

import (
	"log"
	"sync"
	"time"
)

// rate events per second
type Limiter struct {
	mu   sync.Mutex
	rate float64

	// wait duration for a event
	per_event   time.Duration
	bucket_size int64

	// completion time of last event in the queue
	last time.Time
}

func NewLimiter(r float64, b int64) *Limiter {
	if r == 0.0 {
		log.Fatalf("Rate can't be zero")
	}
	per := time.Duration(1.0 / r * float64(time.Second))
	return &Limiter{
		rate:        r,
		bucket_size: b,
		per_event:   per,
	}
}

func (lim *Limiter) Rate() float64 {
	lim.mu.Lock()
	defer lim.mu.Unlock()
	return lim.rate
}

func (lim *Limiter) Allow() (time.Duration, bool) {
	return lim.Allow1(time.Now())
}

// helper function for Allow
func (lim *Limiter) Allow1(t time.Time) (time.Duration, bool) {
	lim.mu.Lock()
	defer lim.mu.Unlock()

	last_time := lim.last
	if last_time.Before(t) {
		last_time = t
	}
	queue_size := last_time.Sub(t).Nanoseconds() / lim.per_event.Nanoseconds()
	if queue_size >= lim.bucket_size {
		return 0, false
	} else {
		// put event in queue
		completion_time := last_time.Add(lim.per_event)
		lim.last = completion_time
		d := completion_time.Sub(t)
		return d, true
	}
}

func main() {
	wg := sync.WaitGroup{}
	limiter := NewLimiter(10, 4)
	for i := 1; i < 20; i++ {
		if i == 8 {
			time.Sleep(1 * time.Second)
		}
		t, ok := limiter.Allow()
		if ok {
			wg.Add(1)
			go func(i int, t time.Duration) {
				defer wg.Done()
				time.Sleep(t)
				log.Printf("Event %d completed\n", i)
			}(i, t)
		} else {
			log.Printf("ERROR: Rate limit exceeded\n")
		}
	}
	wg.Wait()
}
