package slogfromerror

import (
	"context"

	"github.com/monzo/slog"
)

func Test() {
	ctx := context.Background()
	var err error

	// ok:slog-from-error
	slog.FromError(ctx, "message", err)
	// ok:slog-from-error
	slog.FromError(ctx, "message", err, "foo", 2)
	// ok:slog-from-error
	slog.FromError(context.Background(), "message", err, "foo", 2)

	// ruleid:slog-from-error
	slog.FromError(ctx, "message", nil)

	// ruleid:slog-from-error
	slog.FromError(ctx, "message", nil, "foo", 2)

	// ruleid:slog-from-error
	slog.FromError(context.Background(), "message", nil)
}
