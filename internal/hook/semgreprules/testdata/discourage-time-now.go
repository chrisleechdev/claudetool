package discouragetimenow

import (
	"context"
	"time"

	"github.com/monzo/wearedev/libraries/clock"
)

func GetCurrentTime() time.Time {
	// ruleid:discourage-time-now
	return time.Now()
}

func CheckDeadline(deadline time.Time) bool {
	// ruleid:discourage-time-now
	return time.Now().After(deadline)
}

func GetCurrentTimeClock(ctx context.Context) time.Time {
	// ok:discourage-time-now
	return clock.Now(ctx)
}
