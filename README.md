API rate limiters in golang  

Project provids implementation of sliding window and token bucket algorithm

## Quick start

``` Golang
func main() {
    // window size = 1 second (1000 milisecond)
    // events per window = 2
    limiter1 := slidingwindow.NewLimiter(2, time.Millisecond*1000)
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())
    log.Printf("%v\n", limiter1.Allow())

    // 2 events per second (frequency)
    // burst size = 5
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
```
Output: 
``` Console
2023/11/28 15:42:31 true
2023/11/28 15:42:31 true
2023/11/28 15:42:31 false
2023/11/28 15:42:31 false
2023/11/28 15:42:31 false
2023/11/28 15:42:31 true
2023/11/28 15:42:31 true
2023/11/28 15:42:31 true
2023/11/28 15:42:31 true
2023/11/28 15:42:31 true
2023/11/28 15:42:31 false
2023/11/28 15:42:32 true
2023/11/28 15:42:32 true
2023/11/28 15:42:32 false
2023/11/28 15:42:35 true
2023/11/28 15:42:35 true
2023/11/28 15:42:35 true
2023/11/28 15:42:35 true
2023/11/28 15:42:35 true
2023/11/28 15:42:35 false
2023/11/28 15:42:35 false

[Process exited 0]
```
