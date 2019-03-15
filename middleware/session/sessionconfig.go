package sessions

import (
	"math"
	"net/http"
	"time"
)

var (
	// SessionExpiry is the maximum time which may pass before a session that
	// has not been accessed will be destroyed, hence logging a user out.
	SessionExpiry time.Duration = math.MaxInt64

	// SessionIDExpiry is the maximum duration a session ID can be used before it
	// is changed to a new session ID. This helps prevent session hijacking. It
	// may be set to 0, leading to a session ID change with every request.
	// However, this will increase the load on the session persistence layer
	// considerably.
	//
	// Note that expired session IDs will remain active for the duration of
	// SessionIDGracePeriod (leading to session ID overlaps) to avoid race
	// conditions when multiple requests are issued at nearly the same time.
	//每隔1小时刷新session，防止被盗取
	SessionIDExpiry = time.Hour

	// SessionIDGracePeriod is the duration for a replaced (old) session ID to
	// remain active so multiple concurrent requests from the browser don't
	// accidentally lead to session loss. While the default of five minutes may
	// appear long, in a mobile context or other slow networks, it is a reasonable
	// time.
	SessionIDGracePeriod = 5 * time.Minute

	// SessionCookie is the name of the session cookie that will contain the
	// session ID.
	SessionCookie = "id"

	// NewSessionCookie is used to create new session cookies or to renew them.
	// The "Name" and "Value" fields need not be set. It is recommended that you
	// overwrite the default implementation with your specific defaults,
	// especially the "Domain", "Path", and "Secure" fields. Be sure to set
	// "Secure" to true when using TLS (HTTPS). For more information on cookies,
	// refer to:
	//
	//     - https://tools.ietf.org/html/rfc6265
	//     - https://en.wikipedia.org/wiki/HTTP_cookie#Cookie_attributes
	NewSessionCookie = func() *http.Cookie {
		return &http.Cookie{ // Default lifetime is 10 years (i.e. forever).
			Expires:  time.Now().Add(10 * 365 * 24 * time.Hour), // For IE, other browsers will use MaxAge.
			MaxAge:   10 * 365 * 24 * 60 * 60,
			HttpOnly: true,

			// Uncomment and edit the following fields for production use:
			//Domain: "www.example.com",
			//Path:   "/",
			//Secure: true,
		}
	}

	// MaxSessionCacheSize is the maximum size of the local sessions cache. If
	// this value is 0, nothing is cached. If this value is negative, the cache
	// may expand indefinitely. When the maximum size is reached, sessions with
	// the oldest access time are discarded. They are also removed from the cache
	// when their age exceeds SessionCacheExpiry. (This is checked whenever the
	// cache is accessed.)
	MaxSessionCacheSize = 1024 * 1024

	// SessionCacheExpiry is the maximum duration an inactive session will remain
	// in the local cache.
	SessionCacheExpiry = time.Hour
)
