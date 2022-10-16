package middleware

import (
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func LimitRate(next http.Handler) http.Handler {
	type client struct {
		limiter *rate.Limiter
		lastSeen time.Time 
	}

	var (
		mu sync.Mutex
		clients = make(map[string]*client)
	)

	// Remove old entries from the clients map every minute.
	go func() {
		time.Sleep(time.Minute)

		// Lock the mutex to prevent any new requests from being handled while we're
		// cleaning up the map.
		mu.Lock()

		// Loop through the map, removing any entries that haven't been seen in the
		// last minute.
		for ip, client := range clients {
			if time.Since(client.lastSeen) > time.Minute {
				delete(clients, ip)
			}
		}

		// Unlock the mutex
		mu.Unlock()
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract the IP address from the request.
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		// Lock the mutex to prevent the map from being accessed by multiple goroutines at the same time.
		mu.Lock()

		// If we don't have a limiter for this IP address, create one.
		if _, ok := clients[ip]; !ok {
			clients[ip] = &client{limiter: rate.NewLimiter(2, 4)}
		}

		// Update the last seen time for the client.
		clients[ip].lastSeen = time.Now()

		// Check if the client is allowed to make another request.
		if !clients[ip].limiter.Allow() {
			mu.Unlock()
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}

		// Unlock the mutex.
		mu.Unlock()

		next.ServeHTTP(w, r)
	})
}
