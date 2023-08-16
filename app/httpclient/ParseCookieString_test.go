package httpclient

import (
	"net/http"
	"testing"
	"time"
)

func TestParseCookieString(t *testing.T) {
	t.Run("Parsing cookie in a string", func(t *testing.T) {
		var cookies []*http.Cookie
		cookieString := ".domain.com	TRUE	/path	TRUE	1725656330	name	value"
		parsedCookies := ParseCookieString(cookieString, cookies)
		targetCookie := &http.Cookie{
			Domain:   ".domain.com",
			HttpOnly: true,
			Path:     "/path",
			Secure:   true,
			Expires:  time.Unix(1725656330, 0),
			Name:     "name",
			Value:    "value",
		}
		parsedCookie := parsedCookies[0]
		if parsedCookie.Domain != targetCookie.Domain &&
			parsedCookie.HttpOnly != targetCookie.HttpOnly &&
			parsedCookie.Path != targetCookie.Path &&
			parsedCookie.Secure != targetCookie.Secure &&
			parsedCookie.Expires != targetCookie.Expires &&
			parsedCookie.Name != targetCookie.Name &&
			parsedCookie.Value != targetCookie.Value {
			t.Errorf("Couldn't be parsed \"%s\" into cookie", cookieString)
		}
	})
}
