package main

import (
	"log"
	"time"

	"github.com/2asm/ratelimiter/slidingWindow"
	"github.com/2asm/ratelimiter/tokenBucket"
)

func main() {
	limiter1 := slidingwindow.NewLimiter(2, time.Millisecond*1000)
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())
	log.Printf("%v\n", limiter1.Allow())

	limiter2 := tokenbucket.NewLimiter(2, 5)
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
}
