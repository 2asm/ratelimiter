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
2023/11/30 15:12:16.459216 true
2023/11/30 15:12:16.459300 true
2023/11/30 15:12:16.459304 false
2023/11/30 15:12:16.459308 false
2023/11/30 15:12:16.459311 false

token bucket
2023/11/30 15:12:16.459318 true
2023/11/30 15:12:16.459322 true
2023/11/30 15:12:16.459325 true
2023/11/30 15:12:16.459328 true
2023/11/30 15:12:16.459332 true
2023/11/30 15:12:16.459336 false
2023/11/30 15:12:17.459651 true
2023/11/30 15:12:17.459691 true
2023/11/30 15:12:17.459704 false
2023/11/30 15:12:20.460532 true
2023/11/30 15:12:20.460580 true
2023/11/30 15:12:20.460599 true
2023/11/30 15:12:20.460613 true
2023/11/30 15:12:20.460633 true
2023/11/30 15:12:20.460689 false
2023/11/30 15:12:20.460705 false

leaky bucket
2023/11/30 15:12:20.461010 ERROR: Rate limit exceeded
2023/11/30 15:12:20.461036 ERROR: Rate limit exceeded
2023/11/30 15:12:20.461054 ERROR: Rate limit exceeded
2023/11/30 15:12:20.461068 ERROR: Rate limit exceeded
2023/11/30 15:12:20.961323 Event 1 completed
2023/11/30 15:12:21.460950 Event 2 completed
2023/11/30 15:12:21.962128 Event 3 completed
2023/11/30 15:12:22.461984 Event 4 completed
2023/11/30 15:12:22.961700 Event 5 completed

[Process exited 0]
```
