package http

import (
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net/http"
	"solid-server/constants"
	"strconv"
	"strings"
	"time"
)

// https://www.youtube.com/watch?v=_sB2E1XnzOY

var (
	retryTimes int
	rawCookie  string
	userAgent  string
	refer      string
	debug      bool
)

// Options 일반적인 요청 옵션을 정의
type Options struct {
	RetryTimes int
	Cookie     string
	UserAgent  string
	Refer      string
	Debug      bool
	Silent     bool
}

// SetOptions 공통 요청 옵션을 설정
func SetOptions(opt Options) {
	retryTimes = opt.RetryTimes
	rawCookie = opt.Cookie
	userAgent = opt.UserAgent
	refer = opt.Refer
	debug = opt.Debug
}

// Request 요청
func Request(method, url string, body io.Reader, headers map[string]string) (*http.Response, error) {
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		DisableCompression:  true,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   15 * time.Minute,
	}

	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return nil, err
	}
	for k, v := range constants.FakeHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	if _, ok := headers["Referer"]; !ok {
		req.Header.Set("Referer", url)
	}

	if rawCookie != "" {
		// parse cookies in Netscape HTTP cookie format
		cookies, _ := ParseCookieString(rawCookie)
		if len(cookies) > 0 {
			for _, c := range cookies {
				req.AddCookie(c)
			}
		} else {
			// cookie is not Netscape HTTP format, set it directly
			// a=b; c=d
			req.Header.Set("Cookie", rawCookie)
		}
	}

	if userAgent != "" {
		req.Header.Set("User-Agent", userAgent)
	}

	if refer != "" {
		req.Header.Set("Referer", refer)
	}

	var (
		res          *http.Response
		requestError error
	)
	for i := 0; ; i++ {
		res, requestError = client.Do(req)
		if requestError == nil && res.StatusCode < 400 {
			break
		} else if i+1 >= retryTimes {
			var err error
			if requestError != nil {
				err = fmt.Errorf("request error: %v", requestError)
			} else {
				err = fmt.Errorf("%s request error: HTTP %d", url, res.StatusCode)
			}
			return nil, err
		}
		time.Sleep(1 * time.Second)
	}

	return res, nil
}

// Headers return the HTTP Headers of the url
func Headers(url, refer string) (http.Header, error) {
	headers := map[string]string{
		"Referer": refer,
	}

	res, err := Request(http.MethodGet, url, nil, headers)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close() // nolint
	return res.Header, nil
}

// Size get size of the url
func Size(url, refer string) (int64, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return 0, err
	}
	s := h.Get("Content-Length")
	if s == "" {
		return 0, errors.New("Content-Length is not present")
	}
	size, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return 0, err
	}
	return size, nil
}

// ContentType get Content-Type of the url
func ContentType(url, refer string) (string, error) {
	h, err := Headers(url, refer)
	if err != nil {
		return "", err
	}
	s := h.Get("Content-Type")
	// handle Content-Type like this: "text/html; charset=utf-8"
	return strings.Split(s, ";")[0], nil
}

func Cookie() {

}
