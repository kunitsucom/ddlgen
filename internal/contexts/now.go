package contexts

import (
	"context"
	"time"
)

type contextKeyNow struct{}

func Now(ctx context.Context) time.Time {
	if now, ok := ctx.Value(contextKeyNow{}).(time.Time); ok {
		return now
	}

	return time.Now()
}

func WithNow(ctx context.Context, now time.Time) context.Context {
	return context.WithValue(ctx, contextKeyNow{}, now)
}

func WithNowString(ctx context.Context, layout string, value string) context.Context {
	v, err := time.Parse(layout, value)
	if err != nil {
		v = time.Now()
	}

	return WithNow(ctx, v)
}
