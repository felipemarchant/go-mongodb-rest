package utils

import (
	"context"
	"time"
)

func ContextWithTimeout() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), 100*time.Second)
}
