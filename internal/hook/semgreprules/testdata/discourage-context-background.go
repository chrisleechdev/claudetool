package discouragecontextbackground

import "context"

func DoWork() {
	// ruleid:discourage-context-background
	ctx := context.Background()
	_ = ctx
}

func DoWorkTODO() {
	// ruleid:discourage-context-background
	ctx := context.TODO()
	_ = ctx
}

func InlineBackground() {
	// ruleid:discourage-context-background
	Process(context.Background())
}

func InlineTODO() {
	// ruleid:discourage-context-background
	Process(context.TODO())
}

func CorrectUsage(ctx context.Context) {
	// ok:discourage-context-background
	Process(ctx)
}

func Process(ctx context.Context) {}
