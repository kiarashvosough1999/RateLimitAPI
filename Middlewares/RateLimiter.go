//  Copyright 2023 KiarashVosough and other contributors
//
//  Permission is hereby granted, free of charge, to any person obtaining
//  a copy of this software and associated documentation files (the
//  Software"), to deal in the Software without restriction, including
//  without limitation the rights to use, copy, modify, merge, publish,
//  distribute, sublicense, and/or sell copies of the Software, and to
//  permit persons to whom the Software is furnished to do so, subject to
//  the following conditions:
//
//  The above copyright notice and this permission notice shall be
//  included in all copies or substantial portions of the Software.
//
//  THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
//  EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
//  MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
//  NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
//  LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
//  OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
//  WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

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

// NewRateLimiter create a new RateLimiter instance.
// `every` is in seconds.
// bucket determines how many attempt can client have within the `every` interval.
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

// get limiter for visitor, if the visitor is new instantiate a new visitor.
// Note that visitors that did not send any request for 3 minutes will be eliminated from the list.
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

		// for reverse proxy we should consider check this header
		//log.Print(http.Header.Get(r.Header, "X-Forwarded-For"))
		//log.Print(http.Header.Get(r.Header, "X-Real-IP"))
		if err != nil {
			log.Print(err.Error())
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// get limiter based on the client's ip
		limiter := l.getVisitor(ip)

		// check if limiter has available token
		if limiter.Allow() == false {
			http.Error(w, "Too many request, wait for 20 seconds before retrying", http.StatusTooManyRequests)
			return
		}

		// serve next handler in chain
		next.ServeHTTP(w, r)
	})
}
