package sessions

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// A session is usually generated with the Start() function and may
// be destroyed with the Destroy() function.
//
// Sessions are uniquely identified by their session ID. This session ID is
// regenerated, i.e. exchanged, regularly to prevent others from hijacking
// sessions. This can be done explicitly with the RegenerateID() function. And
// it happens automatically based on the rules defined in this package (see
// package variables for details).

type Session struct {
	sync.RWMutex
	id          string    // The session ID. Will not be saved with the session.
	user        int64     // The session user. If nil, no user is attached to this session.
	created     time.Time // created time
	lastAccess  time.Time // The last time the session was accessed through this API.
	referenceID string    // If this session's ID was replaced, this is the ID of the newer session.
}

// Start returns a session for the given HTTP request. Because this function
// may manipulate browser cookies, it must be called before any text is written
// to the response writer.
//
// Sessions are returned from the local cache if contained.
//
// A nil value may also be returned if "createIfNew" is false and no session was
// previously assigned to this user. Note that if the user's browser rejects
// cookies, this will cause a new session to be created with every request. You
// will also want to respect any privacy laws regarding the use of cookies,
// user and session data.
//
// The following package variables influence the session handling (see their
// comments for details):
//
//   - SessionExpiry
//   - SessionIDExpiry
//   - SessionCookie
//   - NewSessionCookie
func Start(response http.ResponseWriter, request *http.Request, createIfNew bool) (*Session, error) {
	// Get the session ID from the cookie.
	var id string // The session ID. Empty if it could not be determined.
	cookie, err := request.Cookie(SessionCookie)
	if err == nil {
		id = cookie.Value //取出cookie中存的session id
	}

	//用session id从缓存map中取出该session
	var session *Session //session的值被初始化为nil
	if len(id) == 24 {
		// Get the session.
		session, err = sessions.Get(id)
		if err != nil {
			return nil, fmt.Errorf("Could not get session from cache: %s", err)
		}

		// 若该id的session找不到，仍为nil，则删除cookie
		if session == nil {
			deleteCookie(cookie, response)
		}
	}
	//如果找到了session，检查是否有效
	if session != nil {
		session.RLock()
		timeUntouched := time.Since(session.lastAccess) //距离上次访问已过去多久
		age := time.Since(session.created)              //距离创建已过去多久
		session.RUnlock()
		valid := true

		//该session是否已过期？
		if timeUntouched >= SessionExpiry {
			valid = false
		}

		if !valid {
			// 已过期，从缓存中删除无效session，并且删除cookie
			if err = session.Destroy(response, request); err != nil {
				return nil, fmt.Errorf("Could not destroy expired session: %s", err)
			}
			session = nil
		} else {
			// 并未过期，是否需要刷新？每隔1小时就刷新
			if session.referenceID == "" && age >= SessionIDExpiry {
				// 需要刷新
				err = session.RegenerateID(response)
				if err != nil {
					return nil, err
				}
			} else if age >= SessionIDExpiry+SessionIDGracePeriod {
				// Grace period expired. 宽限期已过，移除该session
				sessions.Delete(id)
			}

			// Leave the cookie for now, it may be changed by another request. If
			// not, it will be deleted with the next request. In any case, it's
			// illegal to access this session.
			return nil, errors.New("Session expired")
		}

		// If this is a reference session, get the original one.
		if session.referenceID != "" {
			// Redirect cookie to reference session.
			cookie = NewSessionCookie()
			cookie.Name = SessionCookie
			cookie.Value = session.referenceID
			http.SetCookie(response, cookie)

			// Get the referenced session.获取该用户的新session
			session, err = sessions.Get(session.referenceID)
			if err != nil {
				return nil, fmt.Errorf("Could not get referenced session: %s", err)
			}
			if session == nil {
				return nil, errors.New("Reference session not found")
			}
		}

		// We have a valid session.
		session.Lock()
		defer session.Unlock()
		session.lastAccess = time.Now()
		return session, nil
	}

	if session == nil {
		// 该用户还没有session，生成新session给他
		if !createIfNew {
			// And we don't want any.
			return nil, nil
		}

		id, err = generateSessionID()
		if err != nil {
			return nil, fmt.Errorf("Could not generate new session ID: %s", err)
		}
		session = &Session{
			id:         id,
			created:    time.Now(),
			lastAccess: time.Now(),
		}
		sessions.Set(session)

		// Also set the cookie.
		cookie = NewSessionCookie()
		cookie.Name = SessionCookie
		cookie.Value = id
		http.SetCookie(response, cookie)
	}

	return session, nil
}

// generateSessionID generates a random 128-bit(16byte=128bit), 用Base64编码成192bit(即24 byte).
// Collision probability is close to zero.
func generateSessionID() (string, error) {
	// For more on collisions:
	// https://en.wikipedia.org/wiki/Birthday_problem
	// http://www.wolframalpha.com/input/?i=1-e%5E(-1000000000*(1000000000-1)%2F(2*2%5E128))
	//Base64编码要求把3个8位字节（3*8=24）转化为4个6位的字节（4*6=24），之后在6位的前面补两个0，
	// 形成8位一个字节的形式。 如果剩下的字符不足3个字节，则用0填充。
	//16个字节，前15个字节转化为20个字节，最后1个字节用0填充成3个字节，共24个字节(byte)
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("Could not generate session ID: %s", err)
	}
	return base64.StdEncoding.EncodeToString(b), nil
}

// RegenerateID generates a new session ID and replaces it in the current
// session. Use this every time there is a change in user privilege level or a
// related change, e.g. when the user access rights change or when their
// password was changed.
//
// To avoid losing sessions when the network is slow or when many requests for
// the same session ID come in at the same time, the old session (with the old
// key) is turned into a reference session which will be valid for a grace
// period (defined in SessionIDGracePeriod). When that reference session is
// requested, the new session will be returned in its place.
func (s *Session) RegenerateID(response http.ResponseWriter) error {
	// Save this session under a new ID.
	oldID := s.id
	id, err := generateSessionID()
	if err != nil {
		return fmt.Errorf("Could not generate replacement session ID: %s", err)
	}
	s.Lock()
	s.id = id //将此session的id替换成刚生成的新id
	s.created = time.Now()
	s.Unlock()
	if err = sessions.Set(s); err != nil {
		return fmt.Errorf("Could not save session under new session ID: %s", err)
	}

	// 把oldID的session保存为一个暂时的reference session，以供宽限期内使用
	refSession := &Session{
		id:          oldID,
		created:     s.created,
		lastAccess:  time.Now().Add(-SessionIDExpiry), //现在的时刻-刷新间隔，若宽限期内有客户端用这个session id则一定会刷新
		referenceID: id,
	}
	if err = sessions.Set(refSession); err != nil {
		return fmt.Errorf("Could not save reference session: %s", err)
	}

	// Delete that reference session after the grace period.
	go func() {
		time.Sleep(SessionIDGracePeriod)
		sessions.Delete(oldID)
	}()

	// Change the cookie.
	cookie := NewSessionCookie()
	cookie.Name = SessionCookie
	cookie.Value = id
	http.SetCookie(response, cookie)

	return nil
}

//从缓存中删除该session，并删除客户端的cookie
// The session should not be used anymore after this call.
func (s *Session) Destroy(response http.ResponseWriter, request *http.Request) error {
	// Delete session from cache.
	sessions.Delete(s.id)
	// Get the session cookie and delete it.
	cookie, err := request.Cookie(SessionCookie)
	if err != nil {
		return fmt.Errorf("Could not retrieve session cookie: %s", err)
	}
	deleteCookie(cookie, response)

	return nil
}

// deleteCookie deletes a cookie from the user's browser.
func deleteCookie(cookie *http.Cookie, response http.ResponseWriter) {
	delCookie := *cookie
	delCookie.Value = "deleted"
	delCookie.Expires = time.Unix(0, 0)
	delCookie.MaxAge = -1
	http.SetCookie(response, &delCookie)
}

// Expired returns whether or not this session has expired. This is useful to
// frequently purge the session store.运算符优先级：&&优先于||
func (s *Session) Expired() bool {
	s.RLock()
	defer s.RUnlock()
	return s.referenceID != "" && time.Since(s.lastAccess) >= SessionIDGracePeriod ||
		time.Since(s.lastAccess) >= SessionExpiry &&
			time.Since(s.created) >= SessionIDExpiry+SessionIDGracePeriod
}

// LastAccess returns the time this session was last accessed.
func (s *Session) LastAccess() time.Time {
	s.RLock()
	defer s.RUnlock()
	return s.lastAccess
}

// LogIn assigns a user to this session, replacing any previously assigned user.
// If "exclusive" is set to true, all other sessions of this user will be
// deleted, effectively logging them out of any existing sessions first. This
// requires that Persistence.UserSessions() returns all of a user's sessions.
//
// A call to this function also causes a session ID change for security reasons.
// It must be called before any non-header content is sent to the browser.
func (s *Session) LogIn(userID int64, exclusive bool, response http.ResponseWriter) error {
	// First, log user out of existing sessions.
	if exclusive {
		if err := LogOut(userID); err != nil {
			return fmt.Errorf("Could not log user out of existing sessions: %s", err)
		}
	} else {
		s.LogOut()
	}

	// Log user into this session.
	s.Lock()
	s.user = userID
	s.Unlock()
	if err := sessions.Set(s); err != nil {
		return fmt.Errorf("Could not update session cache: %s", err)
	}

	// Switch session ID.
	if err := s.RegenerateID(response); err != nil {
		return fmt.Errorf("Could not switch session ID: %s", err)
	}

	return nil
}

// LogOut logs the currently logged in user out of this session.
//
// Note that the session will still be alive. If you want to destroy the
// current session, too, call Destroy() afterwards.
//
// If no user is logged into this session, nothing happens.
func (s *Session) LogOut() error {
	s.Lock()

	// Do we have a user at all?
	if s.user == 0 {
		s.Unlock()
		return nil
	}

	// Log user out of this session.
	s.user = 0
	s.Unlock()

	return nil //待处理的error
}

// LogOut logs the user with the given ID out of all sessions, returning all IDs of sessions that contain this user.
func LogOut(userID int64) error {
	// Get all sessions of this user.
	var sessionIDs []string
	for k, v := range sessions.sessions {
		if v.user == userID {
			sessionIDs = append(sessionIDs, k)
		}
	}
	// Unset user in each session.
	for _, sessionID := range sessionIDs {
		session, err := sessions.Get(sessionID)
		if err != nil {
			return err
		}
		session.Lock()
		session.user = 0
		session.Unlock()
		if err := sessions.Set(session); err != nil {
			return err
		}
	}

	return nil
}
