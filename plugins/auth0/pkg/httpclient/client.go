package httpclient

import (
	"net/http"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

type RateLimit struct {
	Remaining int
	ResetTime time.Time
}

func (rl RateLimit) Wait() {
	if rl.Remaining < 1 {
		if rl.ResetTime.Before(time.Now()) {
			return
		}

		duration := time.Until(rl.ResetTime)
		time.Sleep(duration)
	}
}

type Transport struct {
	roundTripperWrap http.RoundTripper
	rateLimiter      *RateLimit
}

func (c *Transport) RoundTrip(r *http.Request) (*http.Response, error) {
	if c.rateLimiter != nil {
		c.rateLimiter.Wait()
	}

	resp, err := c.roundTripperWrap.RoundTrip(r)
	if err != nil || resp.StatusCode >= 400 {
		return resp, err
	}

	rl, err := parseRateLimit(resp)
	if err != nil {
		return resp, err
	}

	c.rateLimiter = rl

	return resp, err
}

func NewTransport(transportWrap http.RoundTripper) http.RoundTripper {
	return &Transport{
		roundTripperWrap: transportWrap,
	}
}

func parseRateLimit(resp *http.Response) (*RateLimit, error) {
	var rl RateLimit

	if remaining := resp.Header.Get("X-Ratelimit-Remaining"); remaining != "" {
		rl.Remaining, _ = strconv.Atoi(remaining)
	}

	reset, err := strconv.Atoi(resp.Header.Get("X-Ratelimit-Reset"))
	if err != nil {
		return nil, errors.Wrapf(err, "failed to parse X-RateLimit-Reset header")
	}

	rl.ResetTime = time.Unix(int64(reset), 0)

	return &rl, nil
}
