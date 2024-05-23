package backoff

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

// Config configures a Backoff
type Config struct {
	Min     time.Duration // start backoff at this level
	Max     time.Duration // increase exponentially to this level
	Retries int           // give up after this many; zero means infinite retries
}

// Backoff implements exponential backoff with randomized wait times
type Backoff struct {
	cfg          Config
	ctx          context.Context
	retries      int
	nextDelayMin time.Duration
	nextDelayMax time.Duration
}

// New creates a Backoff object. Pass a Context that can also terminate the operation.
func New(ctx context.Context, cfg Config) *Backoff {
	return &Backoff{
		cfg:          cfg,
		ctx:          ctx,
		nextDelayMin: cfg.Min,
		nextDelayMax: doubleDuration(cfg.Min, cfg.Max),
	}
}

// Reset the Backoff back to its initial condition
func (b *Backoff) Reset() {
	b.retries = 0
	b.nextDelayMin = b.cfg.Min
	b.nextDelayMax = doubleDuration(b.cfg.Min, b.cfg.Max)
}

// Ongoing returns true if caller should keep going
func (b *Backoff) Ongoing() bool {
	// Stop if Context has errored or max retry count is exceeded
	return b.ctx.Err() == nil && (b.cfg.Max == 0 || b.retries < b.cfg.Retries)
}

// Err returns the reason for terminating the backoff, or nil if it didn't terminate
func (b *Backoff) Err() error {
	if b.ctx.Err() != nil {
		return b.ctx.Err()
	}
	if b.cfg.Max != 0 && b.retries >= b.cfg.Retries {
		return fmt.Errorf("terminated after %d retries", b.retries)
	}
	return nil
}

// Retries returns the number of retries so far
func (b *Backoff) Retries() int {
	return b.retries
}

// Wait sleeps for the backoff time then increases the retry count and backoff time
// Returns immediately if Context is terminated
func (b *Backoff) Wait() {
	// Increase the number of retries and get the next delay
	sleepTime := b.NextDelay()

	if b.Ongoing() {
		select {
		case <-b.ctx.Done():
		case <-time.After(sleepTime):
		}
	}
}

func (b *Backoff) NextDelay() time.Duration {
	b.retries++

	// Handle the edge case where the min and max have the same value
	// (or due to some misconfig max is < min)
	if b.nextDelayMin >= b.nextDelayMax {
		return b.nextDelayMin
	}

	// Add a jitter within the next exponential backoff range
	sleepTime := b.nextDelayMin + time.Duration(rand.Int63n(int64(b.nextDelayMax-b.nextDelayMin)))

	// Apply the exponential backoff to calculate the next jitter
	// range, unless we've already reached the max
	if b.nextDelayMax < b.cfg.Max {
		b.nextDelayMin = doubleDuration(b.nextDelayMin, b.cfg.Max)
		b.nextDelayMax = doubleDuration(b.nextDelayMax, b.cfg.Max)
	}

	return sleepTime
}

func doubleDuration(value time.Duration, max time.Duration) time.Duration {
	value = value * 2

	if value <= max {
		return value
	}

	return max
}
