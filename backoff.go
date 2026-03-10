package logs

import "time"

// Default maximum delay between messages when backing off.
const DefaultMaxDelay = time.Hour

// Default initial backoff after the first failure.
const DefaultInitialBackoff = time.Second

// Backoff is a helper for backing off repeated log messages.
// The zero object is ready to use, with defaults according to the constants DefaultMaxDelay and DefaultInitialBackoff.
// Example usage:
//	backoff := logs.MakeBackoff(time.Second)
//	for {
//		if err := doSomething(); err != nil {
//			if backoff.MustSend() {
//				log.Errorf("Something went wrong: %v", err)
//			}
//		} else {
//			backoff.Success()
//		}
//	}

type Backoff struct {
	FailureCount   int           // Number of consecutive failures
	MaxDelay       time.Duration // Maximum delay between messages. If zero, DefaultMaxDelay is used
	InitialBackoff time.Duration // Initial backoff after the first failure
	LastSent       time.Time     // Time when we last sent a message
}

// Return a new Backoff with the given initial backoff duration, and maximum delay between messages.
func MakeBackoff(initialBackoff, maxDelay time.Duration) Backoff {
	return Backoff{
		InitialBackoff: initialBackoff,
		MaxDelay:       maxDelay,
	}
}

// Mark failure, and return true if we should send a message now.
func (b *Backoff) MustSend() bool {
	b.FailureCount++
	maxDelay := b.MaxDelay
	if maxDelay == 0 {
		maxDelay = DefaultMaxDelay
	}
	initialBackoff := b.InitialBackoff
	if initialBackoff == 0 {
		initialBackoff = DefaultInitialBackoff
	}
	failureCount := min(b.FailureCount, 16) // prevent overflow
	backoff := initialBackoff * (1 << (failureCount - 1))
	backoff = min(backoff, maxDelay)
	now := time.Now()
	if now.Sub(b.LastSent) >= backoff {
		b.LastSent = now
		return true
	}
	return false
}

func (b *Backoff) Success() {
	b.FailureCount = 0
	b.LastSent = time.Time{}
}
