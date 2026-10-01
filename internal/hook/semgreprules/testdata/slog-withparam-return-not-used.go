//lint:file-ignore Ignore all the checks
package testdata

import (
	"context"

	"github.com/monzo/slog"
)

func exampleIncorrectUsage() {
	ctx := context.Background()

	// ruleid:slog-withparam-return-not-used
	slog.WithParam(ctx, "key", "value")

	// ruleid:slog-withparam-return-not-used
	slog.WithParams(ctx, map[string]string{
		"key1": "value1",
		"key2": "value2",
	})

	// ok:slog-withparam-return-not-used
	foo(slog.WithParam(ctx, "key", "value"), "bar")

	// ok:slog-withparam-return-not-used
	foo(slog.WithParams(ctx, map[string]string{
		"key": "value",
	}), "bar")
}

func foo(ctx context.Context, params ...any) {}
