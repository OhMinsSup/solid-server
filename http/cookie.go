package http

import (
	"bufio"
	"bytes"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

func ParseCookieString(s string) ([]*http.Cookie, error) {
	return ParserCookie(bytes.NewReader([]byte(s)))
}

// ParserCookie cookie 값을 구문 분석한다.
func ParserCookie(r io.Reader) ([]*http.Cookie, error) {
	scanner := bufio.NewScanner(r)

	cookies := []*http.Cookie{}

	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "#") || line == "" {
			// 주석 및 빈 줄 무시
			continue
		}

		split := strings.Split(line, "\t")
		if len(split) < 7 {
			// 충분히 길지 않은 줄은 무시
			continue
		}

		expiresSplit := strings.Split(split[4], ".")

		expiresSec, err := strconv.Atoi(expiresSplit[0])
		if err != nil {
			return nil, err
		}

		expiresNsec := 0
		if len(expiresSplit) > 1 {
			expiresNsec, err = strconv.Atoi(expiresSplit[1])
			if err != nil {
				expiresNsec = 0
			}
		}

		cookie := &http.Cookie{
			Name:     split[5],
			Value:    split[6],
			Path:     split[2],
			Domain:   split[0],
			Expires:  time.Unix(int64(expiresSec), int64(expiresNsec)),
			Secure:   strings.ToLower(split[3]) == "true",
			HttpOnly: strings.ToLower(split[1]) == "true",
		}
		cookies = append(cookies, cookie)
	}

	return cookies, nil
}
