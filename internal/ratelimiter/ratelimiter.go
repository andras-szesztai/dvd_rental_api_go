package ratelimiter

import "time"

type Limiter interface {
	Allow(ip string) (bool, time.Duration)
}

type Config struct {
	RequestPerTimeFrame int
	TimeFrame           time.Duration
}

func (rl *FixedWindowLimiter) Allow(ip string) (bool, time.Duration) {
	rl.RLock()
	count, exists := rl.clients[ip]
	rl.RUnlock()

	if !exists || count < rl.limit {
		rl.Lock()
		if !exists {
			go rl.resetCount(ip)
		}

		rl.clients[ip]++
		rl.Unlock()
		return true, 0
	}

	return false, rl.window
}

func (rl *FixedWindowLimiter) resetCount(ip string) {
	time.Sleep(rl.window)
	rl.Lock()
	defer rl.Unlock()
	delete(rl.clients, ip)
}
