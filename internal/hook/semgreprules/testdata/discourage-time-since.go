package discouragetimenow

import (
	"context"
	"time"

	"github.com/monzo/wearedev/libraries/clock"
)

func GetTimeSince(ctx context.Context) time.Duration {
	loc, _ := time.LoadLocation("Europe/Madrid")
	ts := time.Date(2020, 1, 5, 12, 5, 0, 0, loc)
	// ruleid:discourage-time-since
	return time.Since(ts)
}

func GetTimeSinceClock(ctx context.Context) time.Duration {
	loc, _ := time.LoadLocation("Europe/Madrid")
	ts := time.Date(2020, 1, 5, 12, 5, 0, 0, loc)
	// ok:discourage-time-since
	return clock.Since(ctx, ts)
}
