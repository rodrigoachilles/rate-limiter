package database

import (
	"context"
	"github.com/rodrigoachilles/rate-limiter/internal/infra/internal_error"
	"time"
)

type Repository interface {
	AddRequest(ctx context.Context, identifier string, expiration time.Duration) (int64, *internal_error.InternalError)
	Block(ctx context.Context, identifier string, expiration time.Duration) *internal_error.InternalError
	IsBlocked(ctx context.Context, identifier string) (bool, *internal_error.InternalError)
}
