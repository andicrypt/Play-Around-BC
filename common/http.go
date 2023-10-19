package common

import (
	"github.com/beefsack/go-rate"
	"math"
	"net/http"
	"time"
)

type HTTP struct {
	MaxConnection            int `yaml:"max_connection" mapstructure:"max_connection"`
	MaxIdleConnectionPerHost int `yaml:"max_idle_connection_per_host" mapstructure:"max_idle_connection_per_host"`
	MaxIdleConnection        int `yaml:"max_idle_connection" mapstructure:"max_idle_connection"`
	Timeout                  int `yaml:"timeout" mapstructure:"timeout"`
	RateLimit                int `yaml:"rate_limit" mapstructure:"rate_limit"`
}

func (h *HTTP) MakeClient(hook Hook) *http.Client {
	t := http.DefaultTransport.(*http.Transport).Clone()
	t.MaxIdleConns = 300
	t.MaxConnsPerHost = 100
	t.MaxIdleConnsPerHost = 200

	if h.MaxConnection > 0 {
		t.MaxConnsPerHost = h.MaxConnection
	}

	if h.MaxIdleConnection > 0 {
		t.MaxIdleConns = h.MaxIdleConnection
	}

	if h.MaxIdleConnectionPerHost > 0 {
		t.MaxIdleConnsPerHost = h.MaxIdleConnectionPerHost
	}
	if h.Timeout == 0 {
		h.Timeout = 60
	}
	rateLimit := math.MaxInt
	if h.RateLimit > 0 {
		rateLimit = h.RateLimit
	}
	return &http.Client{
		Timeout:   time.Duration(h.Timeout) * time.Second,
		Transport: NewThrottledTransport(rateLimit, t, hook),
	}
}

type Hook interface {
	Before(req *http.Request)
	After(res *http.Response)
	OnError(res *http.Response, err error)
}

type defaultHook struct{}

func (d defaultHook) Before(req *http.Request)              {}
func (d defaultHook) After(res *http.Response)              {}
func (d defaultHook) OnError(res *http.Response, err error) {}

// ThrottledTransport Rate Limited HTTP Client
type ThrottledTransport struct {
	roundTripperWrap http.RoundTripper
	ratelimiter      *rate.RateLimiter
	hook             Hook
}

func (c *ThrottledTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	c.ratelimiter.Wait()
	c.hook.Before(r)
	res, err := c.roundTripperWrap.RoundTrip(r)
	if err != nil {
		c.hook.OnError(res, err)
	}
	c.hook.After(res)
	return res, err
}

// NewThrottledTransport wraps an HTTP transport with a rate limiter
// example usage:
// transport := http.DefaultTransport.(*http.Transport).Clone()
// client.Transport = NewThrottledTransport(60, transport) allows 60 requests per second
func NewThrottledTransport(limit int, transportWrap http.RoundTripper, hook Hook) http.RoundTripper {
	if hook == nil {
		hook = &defaultHook{}
	}
	return &ThrottledTransport{
		roundTripperWrap: transportWrap,
		ratelimiter:      rate.New(limit, time.Second),
		hook:             hook,
	}
}
