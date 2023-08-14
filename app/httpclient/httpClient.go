package httpclient

import (
	"bufio"
	"errors"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func NewClient(cookies []*http.Cookie) *http.Client {
	jar := &MyCookieJar{}
	jar.SetCookies(&url.URL{}, cookies)
	client := &http.Client{
		Jar: jar,
	}
	return client
}

func NewHttpClientFromFile(cookieFileName string) (*http.Client, error) {
	f, err := os.Open(cookieFileName)
  if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	cookies := []*http.Cookie{}
	for scanner.Scan() {
		line := scanner.Text()
		cookies = ParseCookieString(line, cookies)
	}
	client := NewClient(cookies)
	return client, nil
}
func NewHttpClientFromString(cookiesString string) (*http.Client, error) {
	var client *http.Client
	var cookies []*http.Cookie
	if cookiesString == "" {
		return client, errors.New("Missing INSTAGRAM_COOKIES_STRING in the invoronment")
	}
	for _, cookie := range(strings.Split(cookiesString, "|,|")) {
		cookies = ParseCookieString(cookie, cookies)
	}
	return NewClient(cookies), nil
}
type MyCookieJar struct {
	cookies []*http.Cookie
}
func (jar *MyCookieJar) SetCookies(u *url.URL, cookies []*http.Cookie) {
	jar.cookies = cookies
}
func (jar *MyCookieJar) Cookies(u *url.URL) []*http.Cookie {
	return jar.cookies
}


