package xtream

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/rclone/rclone/fs"
)

type Client struct {
	host, user, pass string
	userAgent        string
	http             *http.Client

	limiter *RateLimiter
	retries int
	minWait time.Duration
	maxWait time.Duration
}

func NewClient(host, user, pass, ua string,
	rps, retries int,
	minWait, maxWait time.Duration,
) *Client {

	if ua == "" {
		ua = "rclone-xtream/1.0"
	}

	return &Client{
		host:      host,
		user:      user,
		pass:      pass,
		userAgent: ua,
		http:      &http.Client{},
		limiter:   NewRateLimiter(rps),
		retries:   retries,
		minWait:   minWait,
		maxWait:   maxWait,
	}
}

func (c *Client) newRequest(method, url string) (*http.Request, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", c.userAgent)
	return req, nil
}

func (c *Client) do(ctx context.Context, req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", c.userAgent)

	var resp *http.Response
	var err error

	for attempt := 0; attempt <= c.retries; attempt++ {
		if err = c.limiter.Wait(ctx); err != nil {
			return nil, err
		}

		resp, err = c.http.Do(req)

		if !shouldRetry(err, resp) {
			return resp, err
		}

		if resp != nil {
			resp.Body.Close()
		}

		time.Sleep(backoff(attempt, c.minWait, c.maxWait))
	}

	return resp, err
}

func (c *Client) call(action string, out any) error {
	fs.Infof(c, "Api call %s", action)

	url := fmt.Sprintf(
		"%s/player_api.php?username=%s&password=%s&action=%s",
		c.host, c.user, c.pass, action,
	)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	resp, err := c.do(context.Background(), req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return json.NewDecoder(resp.Body).Decode(out)
}
