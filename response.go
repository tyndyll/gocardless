package gocardless

import (
	"net/http"
	"strconv"
	"time"
)

const (
	rateLimitHeader          = `RateLimit-Limit`
	rateLimitRemainingHeader = `RateLimit-Remaining`
	rateLimitResetHeader     = `RateLimit-Reset`
)

type Response struct {
	*http.Response
}

func (resp *Response) RateLimit() int {
	value, err := strconv.Atoi(resp.Header.Get(rateLimitHeader))
	if err != nil {
		return 0
	}
	return value
}

func (resp *Response) RateLimitRemaining() int {
	value, err := strconv.Atoi(resp.Header.Get(rateLimitRemainingHeader))
	if err != nil {
		return 0
	}
	return value
}

func (resp *Response) RateReset() time.Time {
	value, err := time.Parse(time.RFC1123, resp.Header.Get(rateLimitResetHeader))
	if err != nil {
		t := &time.Time{}
		return *t
	}
	return value
}
