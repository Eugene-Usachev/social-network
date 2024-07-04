package cache

import "context"

type Cache interface {
	IsNegativeCase(ctx context.Context, key string) bool

	GetString(ctx context.Context, key string) (res string, isExist bool)
	GetBytes(ctx context.Context, key string) (res []byte, isExist bool)

	SetNegativeCase(ctx context.Context, key string)
	SetString(ctx context.Context, key string, value string)
	SetBytes(ctx context.Context, key string, value []byte)

	Delete(ctx context.Context, key string) error
}
