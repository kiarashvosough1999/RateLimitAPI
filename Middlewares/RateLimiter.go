package Middlewares

import (
	"golang.org/x/time/rate"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

// Visitor Create a custom Visitor struct which holds the rate limiter for each
// Visitor and the last time that the Visitor was seen.
type Visitor struct {
	limiter  *rate.Limiter
	lastSeen time.Time
}

type RateLimiter struct {
	visitors map[string]*Visitor
	mu       sync.Mutex
	every    rate.Limit
	bucket   int
}

// NewRateLimiter create a new RateLimiter
// every: is in seconds
func NewRateLimiter(every time.Duration, bucket int) *RateLimiter {
	limiter := RateLimiter{
		make(map[string]*Visitor),
		sync.Mutex{},
		rate.Every(every * time.Second),
		bucket,
	}
	go limiter.cleanupVisitors()
	return &limiter
}

// get limiter for visitor, if the visitor is new instantiate a new visitor
// Note that visitors that did not send request for 3 minutes will be eliminated from the list.
func (l *RateLimiter) getVisitor(ip string) *rate.Limiter {
	l.mu.Lock()
	defer l.mu.Unlock()

	v, exists := l.visitors[ip]
	if !exists {

		limiter := rate.NewLimiter(l.every, l.bucket)
		l.visitors[ip] = &Visitor{limiter, time.Now()}
		return limiter
	}

	// Update the last seen time for the Visitor.
	v.lastSeen = time.Now()
	return v.limiter
}

// Every minute check the map for visitors that haven't been seen for
// more than 3 minutes and delete the entries.
func (l *RateLimiter) cleanupVisitors() {
	for {
		time.Sleep(time.Minute)

		l.mu.Lock()
		for ip, v := range l.visitors {
			if time.Since(v.lastSeen) > 3*time.Minute {
				delete(l.visitors, ip)
			}
		}
		l.mu.Unlock()
	}
}

// Limiter Middleware which check for visitor and apply rate limit
func (l *RateLimiter) Limiter(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// find ip of client
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		//log.Print(http.Header.Get(r.Header, "X-Forwarded-For"))
		//log.Print(http.Header.Get(r.Header, "X-Real-IP"))
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// get limiter based on the client ip
		limiter := l.getVisitor(ip)

		// check if limiter has available token
		if limiter.Allow() == false {
			http.Error(w, "Too many request", http.StatusTooManyRequests)
			return
		}

		// serve next handler in chain
		next.ServeHTTP(w, r)
	})
}
