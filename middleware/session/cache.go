package sessions

import (
	"sync"
	"time"
)

// a simple cache for user sessions. It is used implicitly by all sessions functions.
// Member functions should not be called while sessions are locked.
type cache struct {
	sync.Mutex
	sessions map[string]*Session
}

// sessions is the global sessions cache.
var sessions *cache

// initCache initializes the global sessions cache.
func initCache() {
	sessions = &cache{
		sessions: make(map[string]*Session),
	}
}

// Get returns a session with the given ID from the cache.  If no such session exists,
// a nil session may be returned. This function does not update the session's last access date.
func (c *cache) Get(id string) (*Session, error) {
	c.Lock()
	defer c.Unlock()
	// Do we have a cached session?
	session := c.sessions[id]
	return session, nil
}

// Set inserts or updates a session in the cache.
func (c *cache) Set(session *Session) error {
	c.Lock()
	defer c.Unlock()
	session.Lock()
	session.lastAccess = time.Now()
	id := session.id
	session.Unlock()

	// 抖抖麻袋，顶上空出1个位置
	var requiredSpace int
	if _, ok := c.sessions[id]; !ok {
		requiredSpace = 1
	}
	c.compact(requiredSpace)

	// Save in cache.
	if MaxSessionCacheSize != 0 {
		c.sessions[id] = session
	}
	return nil
}

// Delete deletes a session. A logged-in user will be logged out.
func (c *cache) Delete(id string) {
	c.Lock()
	defer c.Unlock()
	// Remove from cache.
	delete(c.sessions, id)
}

// compact drops sessions from the cache to make space for the given number
// of sessions. It also drops sessions that have been in the cache longer than
// SessionCacheExpiry. The number of dropped sessions are returned.
//
// This function does not synchronize concurrent access to the cache.
func (c *cache) compact(requiredSpace int) (int, error) {
	// Cache may still grow.
	if len(c.sessions)+requiredSpace <= MaxSessionCacheSize {
		return 0, nil
	}

	// Drop the oldest sessions.
	var dropped int
	if requiredSpace > MaxSessionCacheSize {
		requiredSpace = MaxSessionCacheSize // We can't request more than is allowed.
	}
	for len(c.sessions)+requiredSpace > MaxSessionCacheSize {
		// Find oldest sessions and delete them.
		var (
			oldestAccessTime time.Time
			oldestSessionID  string
		)
		for id, session := range c.sessions {
			session.RLock()
			before := session.lastAccess.Before(oldestAccessTime)
			session.RUnlock()
			if oldestSessionID == "" || before {
				oldestSessionID = id
				oldestAccessTime = session.lastAccess
			}
		}
		delete(c.sessions, oldestSessionID)
		dropped++
	}

	return dropped, nil
}
