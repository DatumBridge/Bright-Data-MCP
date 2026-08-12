package mcp

import (
	"crypto/rand"
	"encoding/hex"
	"sync"
	"time"
)

const SessionTTL = 30 * time.Minute

// SessionStore tracks Streamable HTTP MCP sessions (Mcp-Session-Id) after initialize.
type SessionStore struct {
	mu sync.Mutex
	m  map[string]time.Time
}

func NewSessionStore() *SessionStore {
	return &SessionStore{m: make(map[string]time.Time)}
}

func (s *SessionStore) pruneLocked() {
	cutoff := time.Now().Add(-SessionTTL)
	for id, t := range s.m {
		if t.Before(cutoff) {
			delete(s.m, id)
		}
	}
}

func randomSessionID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// Create allocates a new session id.
func (s *SessionStore) Create() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked()
	id := randomSessionID()
	s.m[id] = time.Now()
	return id
}

// Valid reports whether the session id is active (slides expiry on use).
func (s *SessionStore) Valid(id string) bool {
	if id == "" {
		return false
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.pruneLocked()
	t, ok := s.m[id]
	if !ok {
		return false
	}
	if time.Since(t) > SessionTTL {
		delete(s.m, id)
		return false
	}
	s.m[id] = time.Now()
	return true
}
