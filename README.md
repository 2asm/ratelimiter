# API rate limiters in golang  

Project provides implementation of sliding window, token bucket and leaky bucket algorithm

## Quick start

``` Golang
func main() {
    // window size = 1 second (1000 milisecond)
    // events per window = 2
    limiter1 := slidingwindow.NewLimiter(2, time.Millisecond*1000)
    fmt.Printf("sliding bucket\n")
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())

    // 2 events per second (frequency)
    // burst size = 5
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

    // 2 events per second
    // bucket size = 4 
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
```
Output: 
``` Console
sliding bucket
2023/11/30 14:52:59 true
2023/11/30 14:52:59 true
2023/11/30 14:52:59 false
2023/11/30 14:52:59 false
2023/11/30 14:52:59 false

token bucket
2023/11/30 14:52:59 true
2023/11/30 14:52:59 true
2023/11/30 14:52:59 true
2023/11/30 14:52:59 true
2023/11/30 14:52:59 true
2023/11/30 14:52:59 false
2023/11/30 14:53:00 true
2023/11/30 14:53:00 true
2023/11/30 14:53:00 false
2023/11/30 14:53:03 true
2023/11/30 14:53:03 true
2023/11/30 14:53:03 true
2023/11/30 14:53:03 true
2023/11/30 14:53:03 true
2023/11/30 14:53:03 false
2023/11/30 14:53:03 false

leaky bucket
2023/11/30 14:53:03 ERROR: Rate limit exceeded
2023/11/30 14:53:03 ERROR: Rate limit exceeded
2023/11/30 14:53:03 ERROR: Rate limit exceeded
2023/11/30 14:53:03 ERROR: Rate limit exceeded
2023/11/30 14:53:03 Event 1 completed
2023/11/30 14:53:04 Event 2 completed
2023/11/30 14:53:04 Event 3 completed
2023/11/30 14:53:05 Event 4 completed
2023/11/30 14:53:05 Event 5 completed

[Process exited 0]
```
