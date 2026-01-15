package xtream

import (
	"errors"
	"math/rand"
	"net"
	"net/http"
	"time"
)

func shouldRetry(err error, resp *http.Response) bool {
	if err != nil {
		var netErr net.Error
		if errors.As(err, &netErr) {
			return true
		}
		return true
	}

	if resp == nil {
		return true
	}

	if resp.StatusCode == 429 {
		return true
	}

	return resp.StatusCode >= 500
}

func backoff(attempt int, min, max time.Duration) time.Duration {
	d := min * (1 << attempt)
	if d > max {
		d = max
	}
	// jitter ±25 %
	jitter := time.Duration(rand.Int63n(int64(d/2))) - d/4
	return d + jitter
}
