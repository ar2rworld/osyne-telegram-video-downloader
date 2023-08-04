package httpclient

import (
	"bufio"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
)

func NewHttpClient(cookieFileName string) (*http.Client, error) {
	f, err := os.Open(cookieFileName)
  if err != nil {
		return nil, err
	}
	scanner := bufio.NewScanner(f)
	jar := &MyCookieJar{}
	cookies := []*http.Cookie{}
	for scanner.Scan() {
		line := scanner.Text()
		log.Println(line)
		arguments := strings.Split(line, "\t")

		// 7 values in the line
		if len(arguments) == 7 {
			timestamp, err := strconv.ParseInt(arguments[4], 10, 64)
			if err != nil {
				panic(err)
			}
			timeArg := time.Unix(timestamp, 0)
			cookies = append(
				cookies,
				&http.Cookie{
					Domain: arguments[0],
					HttpOnly: arguments[1] == "TRUE",
					Path: arguments[2],
					Secure: arguments[3] == "TRUE",
					Expires: timeArg,
					Name: arguments[5],
					Value: arguments[6],
				},
			)
		}
	}
	jar.SetCookies(&url.URL{}, cookies)
	client := &http.Client{
		Jar: jar,
	}
	return client, nil
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


