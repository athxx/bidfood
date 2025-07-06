# RateLimiter Comparison

This document explains the differences between the two rate limiters implemented in `main.go`: `RateLimiter` and `Limiter`.

## 1. **RateLimiter**

### Description:

The `RateLimiter` is a custom implementation of a rate limiter using a mutex (`sync.Mutex`) and manual token counting. It operates based on the current timestamp and allows a fixed number of operations per second.

### Key Features:

- **Manual Token Management**: Tokens are manually incremented and reset based on the current second.
- **Blocking Behavior**: The `Allow` method blocks until a token is available.
- **Thread-Safety**: Uses a mutex to ensure thread-safe access to shared state (`timestamp` and `tokens`).

### How It Works:

1. The `Allow` method locks the mutex and checks the current timestamp.
2. If the timestamp has changed (new second), it resets the token count.
3. If tokens are available (less than the rate), it increments the token count and allows the operation.
4. If tokens are exhausted, it blocks until the next second.

---

## 2. **Limiter**

### Description:

The `Limiter` is a more structured implementation of a rate limiter using a buffered channel (`chan struct{}`) and a ticker (`time.Ticker`). It operates based on a fixed interval derived from the rate.

### Key Features:

- **Buffered Channel**: Tokens are managed using a channel, which simplifies concurrency handling.
- **Non-Blocking Token Generation**: A background goroutine periodically generates tokens at a fixed interval.
- **Efficient Design**: Avoids manual locking and timestamp checks.

### How It Works:

1. A ticker generates tokens at regular intervals (`time.Second / rate`) and sends them to the channel.
2. The `Allow` method consumes a token from the channel, blocking if no tokens are available.
3. The channel's buffer size ensures that excess tokens are discarded when the buffer is full.

---

## Comparison Table

| Feature               | `RateLimiter`                     | `Limiter`                          |
| --------------------- | --------------------------------- | ---------------------------------- |
| **Token Management**  | Manual (`timestamp` and `tokens`) | Automatic (`chan struct{}`)        |
| **Concurrency**       | Mutex (`sync.Mutex`)              | Channel-based                      |
| **Blocking Behavior** | Blocks until tokens are available | Blocks until tokens are available  |
| **Token Generation**  | Per second based on timestamp     | Fixed interval using `time.Ticker` |
| **Complexity**        | Higher (manual logic)             | Lower (channel simplifies logic)   |

---

## Summary

- **`RateLimiter`** is a simpler, timestamp-based implementation suitable for basic rate-limiting needs. However, it requires manual management of tokens and synchronization.
- **`Limiter`** is a more robust and efficient implementation, leveraging Go's concurrency primitives (`chan` and `time.Ticker`) for cleaner and more scalable rate-limiting.

Both implementations achieve the same goal of limiting operations per second, but `Limiter` is generally preferred for production use due to its cleaner design and better scalability.

## How to use ?

```bash
cd ratelimiter && go run ./...
```
