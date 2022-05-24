package auth

import (
	"net/http"
	"strings"
)

const (
	HeaderToken        = "token"
	HeaderAuth         = "Authorization"
	HeaderBearer       = "BEARER"
	AccessToken = "accessToken"
)

type TokenLocation int

const (
	TokenLocationNotFound TokenLocation = iota
	TokenLocationHeader
	TokenLocationCookie
)

func (tl TokenLocation) String() string {
	switch tl {
	case TokenLocationNotFound:
		return "Not Found"
	case TokenLocationHeader:
		return "Header"
	case TokenLocationCookie:
		return "Cookie"
	default:
		return "Unknown"
	}
}

func ParseAuthTokenFromRequest(r *http.Request) (string, TokenLocation) {
	authHeader := r.Header.Get(HeaderAuth)

	// Attempt to parse the token from the cookie
	if cookie, err := r.Cookie(AccessToken); err == nil {
		return cookie.Value, TokenLocationCookie
	}

	// Parse the token from the header
	if len(authHeader) > 6 && strings.ToUpper(authHeader[0:6]) == HeaderBearer {
		// Default session token
		return authHeader[7:], TokenLocationHeader
	}

	if len(authHeader) > 5 && strings.ToLower(authHeader[0:5]) == HeaderToken {
		// OAuth token
		return authHeader[6:], TokenLocationHeader
	}

	return "", TokenLocationNotFound
}
