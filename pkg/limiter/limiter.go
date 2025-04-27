package limiter

import (
	"context"
	"github.com/cloudwego/hertz/pkg/app"
	"sync"
	"time"
)

const (
	defaultRate            = 5
	defaultWindow          = time.Second
	defaultCleanupInterval = time.Minute * 10
	defaultBanDuration     = time.Hour
)

type Limiter struct {
	mu          sync.Mutex
	requests    map[string][]int64
	banned      map[string]time.Time
	opts        options
	lastCleanup time.Time
}

var defaultOpts = options{
	rate:   defaultRate,
	window: defaultWindow,
	keyFunc: func(ctx context.Context, c *app.RequestContext) string {
		key := string(c.GetHeader("X-Request-Key"))
		if key == "" {
			key = c.Query("key")
			if key == "" {
				key = c.ClientIP()
			}
		}
		return key
	},
	cleanupInterval: defaultCleanupInterval,
	banDuration:     defaultBanDuration,
}

type options struct {
	rate            int
	window          time.Duration
	keyFunc         func(ctx context.Context, c *app.RequestContext) string
	cleanupInterval time.Duration
	banDuration     time.Duration
}

type Option func(options *options)

// NewLimiter creates a key-based request rate limiter
func NewLimiter(opts ...Option) *Limiter {
	o := defaultOpts
	for _, opt := range opts {
		opt(&o)
	}
	if o.rate <= 0 {
		o.rate = defaultRate
	}
	if o.window <= 0 {
		o.window = defaultWindow
	}
	return &Limiter{
		requests:    make(map[string][]int64),
		opts:        o,
		lastCleanup: time.Now(),
	}
}

// WithRate set maximum request rate
func WithRate(rate int) Option {
	return func(o *options) {
		o.rate = rate
	}
}

// WithWindow defines time duration per window
func WithWindow(window time.Duration) Option {
	return func(o *options) {
		o.window = window
	}
}

// WithKeyFunc set function to get key
func WithKeyFunc(keyFunc func(ctx context.Context, c *app.RequestContext) string) Option {
	return func(o *options) {
		o.keyFunc = keyFunc
	}
}

// WithCleanupInterval set the interval for global cleanup
func WithCleanupInterval(interval time.Duration) Option {
	return func(o *options) {
		o.cleanupInterval = interval
	}
}

// WithBanDuration set the duration for banning users who exceed rate limit
func WithBanDuration(duration time.Duration) Option {
	return func(o *options) {
		o.banDuration = duration
	}
}
func (l *Limiter) GetKey(ctx context.Context, c *app.RequestContext) string {
	return l.opts.keyFunc(ctx, c)
}
func (l *Limiter) IsAllowed(key string) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	now := time.Now()
	if banUntil, isBanned := l.banned[key]; isBanned {
		if now.Before(banUntil) {
			return false
		}
		delete(l.banned, key)
	}
	nowNano := now.UnixNano()
	l.cleanupOldRequests(key, now)
	reqs, exists := l.requests[key]
	if !exists {
		reqs = make([]int64, 0, l.opts.rate)
		l.requests[key] = reqs
	}
	if len(reqs) >= l.opts.rate {
		l.requests[key] = reqs[1:]
		l.banned[key] = now.Add(l.opts.banDuration)
		return false
	}
	l.requests[key] = append(l.requests[key], nowNano)
	return true
}
func (l *Limiter) StartCleanup() {
	go func() {
		ticker := time.NewTicker(l.opts.cleanupInterval)
		defer ticker.Stop()
		for {
			select {
			case <-ticker.C:
				l.periodicCleanup()
			}
		}
	}()
}

// cleanupOldRequests clearing outdated request records
func (l *Limiter) cleanupOldRequests(key string, now time.Time) {
	if reqs, exists := l.requests[key]; exists {
		cutoff := now.Add(-l.opts.window).UnixNano()
		validIdx := 0
		for _, reqTime := range reqs {
			if reqTime >= cutoff {
				reqs[validIdx] = reqTime
				validIdx++
			}
		}
		if validIdx == 0 {
			delete(l.requests, key)
		} else if validIdx < len(reqs) {
			l.requests[key] = reqs[:validIdx]
		}
	}
}
func (l *Limiter) periodicCleanup() {
	now := time.Now()
	if now.Sub(l.lastCleanup) >= l.opts.cleanupInterval {
		l.mu.Lock()
		defer l.mu.Unlock()
		for key, banUntil := range l.banned {
			if now.After(banUntil) {
				delete(l.banned, key)
			}
		}
		for key := range l.requests {
			l.cleanupOldRequests(key, now)
		}
		l.lastCleanup = now
	}
}
