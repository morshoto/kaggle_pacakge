package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

const userAgent = "Mozilla/5.0 (compatible; KaggleDiscussionDownloader/1.0)"
const (
	maxRetries    = 5
	baseBackoff   = 1 * time.Second
	maxBackoffCap = 10 * time.Second
)

// Client wraps an http.Client and its cookie jar together with helpers.
type Client struct {
	http    *http.Client
	cookies map[string]string
	verbose bool
}

func NewClient(verbose bool) *Client {
	jar := &simpleCookieJar{cookies: map[string]string{}}
	return &Client{
		http:    &http.Client{Jar: jar, Timeout: 30 * time.Second},
		cookies: jar.cookies,
		verbose: verbose,
	}
}

func (c *Client) Get(rawURL string, params url.Values) (*http.Response, error) {
	if len(params) > 0 {
		rawURL = rawURL + "?" + params.Encode()
	}
	req, err := http.NewRequest(http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	if xsrf, ok := c.cookies["XSRF-TOKEN"]; ok && xsrf != "" {
		req.Header.Set("X-XSRF-TOKEN", xsrf)
	}
	return c.http.Do(req)
}

func (c *Client) PostJSON(rawURL string, body any) (*http.Response, error) {
	data, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, rawURL, bytes.NewReader(data))
	if err != nil {
		return nil, err
	}
	req.Header.Set("User-Agent", userAgent)
	req.Header.Set("Content-Type", "application/json")
	if xsrf, ok := c.cookies["XSRF-TOKEN"]; ok && xsrf != "" {
		req.Header.Set("X-XSRF-TOKEN", xsrf)
	}
	return c.http.Do(req)
}

func (c *Client) FetchBody(rawURL string, params url.Values) ([]byte, error) {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := c.Get(rawURL, params)
		if err != nil {
			lastErr = err
			c.maybeSleep(nil, attempt)
			continue
		}
		if shouldRetry(resp.StatusCode) {
			lastErr = fmt.Errorf("HTTP %d for %s", resp.StatusCode, rawURL)
			delay := c.maybeSleep(resp, attempt)
			_ = resp.Body.Close()
			if delay > 0 {
				continue
			}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return nil, fmt.Errorf("HTTP %d for %s", resp.StatusCode, rawURL)
		}
		return io.ReadAll(resp.Body)
	}
	return nil, lastErr
}

func (c *Client) FetchJSON(rawURL string, params url.Values, dest any) error {
	data, err := c.FetchBody(rawURL, params)
	if err != nil {
		return err
	}
	return json.Unmarshal(data, dest)
}

func (c *Client) PostJSONDecode(rawURL string, body any, dest any) error {
	var lastErr error
	for attempt := 0; attempt <= maxRetries; attempt++ {
		resp, err := c.PostJSON(rawURL, body)
		if err != nil {
			lastErr = err
			c.maybeSleep(nil, attempt)
			continue
		}
		if shouldRetry(resp.StatusCode) {
			lastErr = fmt.Errorf("HTTP %d for %s", resp.StatusCode, rawURL)
			delay := c.maybeSleep(resp, attempt)
			_ = resp.Body.Close()
			if delay > 0 {
				continue
			}
		}
		defer resp.Body.Close()
		if resp.StatusCode >= 400 {
			return fmt.Errorf("HTTP %d for %s", resp.StatusCode, rawURL)
		}
		return json.NewDecoder(resp.Body).Decode(dest)
	}
	return lastErr
}

func (c *Client) LogInfo(format string, args ...any) {
	if c.verbose {
		log.Printf("[info] "+format, args...)
	}
}

func shouldRetry(status int) bool {
	return status == http.StatusTooManyRequests || status >= http.StatusInternalServerError
}

func (c *Client) maybeSleep(resp *http.Response, attempt int) time.Duration {
	if attempt >= maxRetries {
		return 0
	}
	delay := retryDelay(resp, attempt)
	if delay > 0 {
		c.LogInfo("rate limited, retrying in %s", delay)
		time.Sleep(delay)
	}
	return delay
}

func retryDelay(resp *http.Response, attempt int) time.Duration {
	if resp != nil {
		if v := resp.Header.Get("Retry-After"); v != "" {
			if secs, err := strconv.Atoi(v); err == nil && secs > 0 {
				return time.Duration(secs) * time.Second
			}
			if t, err := http.ParseTime(v); err == nil {
				if d := time.Until(t); d > 0 {
					return d
				}
			}
		}
	}
	backoff := baseBackoff << attempt
	if backoff > maxBackoffCap {
		return maxBackoffCap
	}
	return backoff
}

// simpleCookieJar is a minimal CookieJar that records cookies into a shared map.
type simpleCookieJar struct {
	cookies map[string]string
}

func (j *simpleCookieJar) SetCookies(_ *url.URL, cookies []*http.Cookie) {
	for _, ck := range cookies {
		j.cookies[ck.Name] = ck.Value
	}
}

func (j *simpleCookieJar) Cookies(u *url.URL) []*http.Cookie {
	out := make([]*http.Cookie, 0, len(j.cookies))
	for name, value := range j.cookies {
		out = append(out, &http.Cookie{Name: name, Value: value})
	}
	return out
}
