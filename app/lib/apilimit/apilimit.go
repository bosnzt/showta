package apilimit

import (
    "sync"
    "time"
)

type ApiLimit struct {
    MaxCount int
    Interval time.Duration
    Tokens   int
    LastTick time.Time
    Mutex    sync.Mutex
}

type ApiRateLimiter struct {
    limits map[string]*ApiLimit
    Mutex  sync.Mutex
}

func (l *ApiRateLimiter) SetLimit(list map[string]ApiLimit) {
    l.Mutex.Lock()
    defer l.Mutex.Unlock()
    for api, v := range list {
        l.limits[api] = &ApiLimit{
            MaxCount: v.MaxCount,
            Interval: v.Interval,
            Tokens:   v.MaxCount,
            LastTick: time.Now(),
        }
    }
}

func (l *ApiRateLimiter) Allow(api string) bool {
    l.Mutex.Lock()
    defer l.Mutex.Unlock()

    limit, ok := l.limits[api]
    if !ok {

        return true
    }

    now := time.Now()
    if now.Sub(limit.LastTick) >= limit.Interval {
        limit.Tokens = limit.MaxCount
        limit.LastTick = now
    }

    if limit.Tokens > 0 {
        limit.Tokens--
        return true
    }
    return false
}

func NewApiRateLimiter(list map[string]ApiLimit) *ApiRateLimiter {
    limiter := &ApiRateLimiter{
        limits: make(map[string]*ApiLimit),
    }
    limiter.SetLimit(list)
    return limiter
}
