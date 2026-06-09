package main

import (
	"net/http"
	"sync"
	"time"
)

// --- ADMIN SESSION ---

var (
	sessions   = make(map[string]time.Time)
	sessionsMu sync.Mutex
)

func adminSession(r *http.Request) (string, bool) {
	c, err := r.Cookie("admin")
	if err != nil || c.Value == "" {
		return "", false
	}
	sessionsMu.Lock()
	expiry, ok := sessions[c.Value]
	if !ok || time.Now().After(expiry) {
		delete(sessions, c.Value)
		sessionsMu.Unlock()
		return "", false
	}
	sessionsMu.Unlock()
	return c.Value, true
}
