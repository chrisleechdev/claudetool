package discouragetimenow

import (
	"context"
	"testing"
	"time"

	"github.com/monzo/wearedev/libraries/clock"
	"github.com/monzo/wearedev/libraries/test"
)

func TestBadExample_Direct(t *testing.T) {
	// ruleid:discourage-time-now-in-tests
	now := time.Now()

	_ = now
}

func TestBadExample_Inline(t *testing.T) {
	// ruleid:discourage-time-now-in-tests
	timestamp := time.Now().Unix()

	_ = timestamp
}

func TestBadExample_Comparison(t *testing.T) {
	deadline := time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC)

	// ruleid:discourage-time-now-in-tests
	if time.Now().After(deadline) {
		t.Error("too late")
	}
}

func TestBadExample_Format(t *testing.T) {
	// ruleid:discourage-time-now-in-tests
	formatted := time.Now().Format(time.RFC3339)

	_ = formatted
}

func TestGoodExample(t *testing.T) {
	// ok:discourage-time-now-in-tests
	m, ctx := test.NewEmptyMock(t)
	fixedTime := time.Date(2020, time.January, 1, 0, 0, 0, 0, time.UTC)
	clock.SetMockTime(m, fixedTime)

	result := DoSomething(ctx)

	if !result.Timestamp.Equal(fixedTime) {
		t.Error("unexpected timestamp")
	}
}

func TestGoodExample_ClockNow(t *testing.T) {
	// ok:discourage-time-now-in-tests
	m, ctx := test.NewEmptyMock(t)
	clock.SetMockTime(m, time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC))

	now := clock.Now(ctx)
	_ = now
}

func DoSomething(ctx context.Context) struct{ Timestamp time.Time } {
	return struct{ Timestamp time.Time }{Timestamp: clock.Now(ctx)}
}
