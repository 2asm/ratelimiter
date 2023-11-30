package main

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/2asm/ratelimiter/leakyBucket"
	"github.com/2asm/ratelimiter/slidingWindow"
	"github.com/2asm/ratelimiter/tokenBucket"
)

func main() {
	limiter1 := slidingwindow.NewLimiter(2, time.Millisecond*1000)
    fmt.Printf("sliding bucket\n")
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())

	limiter2 := tokenbucket.NewLimiter(2, 5)
    fmt.Printf("\ntoken bucket\n")
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
    time.Sleep(time.Second*1)
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
    time.Sleep(time.Second*3)
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())
	log.Printf("%v\n", limiter2.Allow())

	wg := sync.WaitGroup{}
	limiter := leakyBucket.NewLimiter(2, 4)
    fmt.Printf("\nleaky bucket\n")
	for i := 1; i < 10; i++ {
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
