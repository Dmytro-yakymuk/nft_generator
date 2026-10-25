package auth

import (
	"net/http"
	"time"
)

// SessionSettings variable cookie settings.
type SessionSettings struct {
	Name string
	Path string
}

// SessionCookie handles cookie authorization.
type SessionCookie struct {
	settings SessionSettings
}

// NewSessionCookie create new cookie authorization with provided settings.
func NewSessionCookie(settings SessionSettings) *SessionCookie {
	return &SessionCookie{
		settings: settings,
	}
}

// GetToken retrieves token from request.
func (cookieAuth *SessionCookie) GetToken(r *http.Request) (string, error) {
	cookie, err := r.Cookie(cookieAuth.settings.Name)
	if err != nil {
		return "", err
	}

	return cookie.Value, nil
}

// SetTokenCookie sets parametrized token cookie that is not accessible from js.
// The cookie is marked Secure when the request came over HTTPS.
func (cookieAuth *SessionCookie) SetTokenCookie(w http.ResponseWriter, r *http.Request, token string) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieAuth.settings.Name,
		Value:    token,
		Path:     cookieAuth.settings.Path,
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: true,
		Secure:   isHTTPS(r),
		SameSite: http.SameSiteStrictMode,
	})
}

// isHTTPS reports whether the request came over HTTPS directly or through
// a TLS-terminating proxy such as ngrok.
func isHTTPS(r *http.Request) bool {
	return r.TLS != nil || r.Header.Get("X-Forwarded-Proto") == "https"
}

// RemoveTokenCookie removes auth cookie that is not accessible from js.
// The path and attributes match SetTokenCookie, otherwise the browser keeps the original cookie.
func (cookieAuth *SessionCookie) RemoveTokenCookie(w http.ResponseWriter, r *http.Request) {
	http.SetCookie(w, &http.Cookie{
		Name:     cookieAuth.settings.Name,
		Value:    "",
		Path:     cookieAuth.settings.Path,
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		HttpOnly: true,
		Secure:   isHTTPS(r),
		SameSite: http.SameSiteStrictMode,
	})
}
