package main

import (
	"net"
	"sync"
	"time"
)

type Limiter struct {
	mu       sync.Mutex
	requests map[string][]time.Time
}

func NewLimiter() *Limiter {
	return &Limiter{
		requests: make(map[string][]time.Time),
	}
}

func (l *Limiter) Allow(clientIP string, s *Server) bool {
	l.mu.Lock()
	defer l.mu.Unlock()

	now := time.Now()
	l.requests[clientIP] = cleanupRequests(l.requests[clientIP], now)

	if len(l.requests[clientIP]) >= s.config.MaxRequestsPerSecond {
		return false
	}

	l.requests[clientIP] = append(l.requests[clientIP], now)
	return true
}

func cleanupRequests(requests []time.Time, now time.Time) []time.Time {
	var cleaned []time.Time
	for _, t := range requests {
		if now.Sub(t) <= time.Second {
			cleaned = append(cleaned, t)
		}
	}
	return cleaned
}

func extractClientIP(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		return ""
	}
	return host
}
