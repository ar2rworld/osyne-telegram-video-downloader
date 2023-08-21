package httpclient

import (
	"net/http"
	"strconv"
	"strings"
	"time"
)

const NumberOfArguments = 7

func ParseCookieString(s string, cookies []*http.Cookie) []*http.Cookie {
	arguments := strings.Split(s, "\t")

	// 7 values in the netscape formatted line
	if len(arguments) == NumberOfArguments {
		timestamp, err := strconv.ParseInt(arguments[4], 10, 64)
		if err != nil {
			panic(err)
		}
		timeArg := time.Unix(timestamp, 0)
		cookies = append(
			cookies,
			&http.Cookie{
				Domain:   arguments[0],
				HttpOnly: arguments[1] == "TRUE",
				Path:     arguments[2],
				Secure:   arguments[3] == "TRUE",
				Expires:  timeArg,
				Name:     arguments[5],
				Value:    arguments[6],
			},
		)
	}
	return cookies
}
