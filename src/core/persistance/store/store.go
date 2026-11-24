package store

import (
	"context"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
)

type InsertOptions struct{}
type UpdateOptions struct{ Upsert *bool }

type Store[T any] interface {
	Insert(ctx context.Context, doc T, options ...InsertOptions) (id any, err error)
	BulkInsert(ctx context.Context, docs []T, options ...InsertOptions) (ids []any, err error)

	Update(ctx context.Context, filter filters.Filtered[T], update T, options ...UpdateOptions) (any, error)
	BulkWrite(ctx context.Context, ops []BulkUpdateOp[T]) ([]any, error)

	One(ctx context.Context, filter filters.Filtered[T]) (T, error)
	Many(ctx context.Context, filter filters.Filtered[T]) ([]T, error)

	DeleteOne(ctx context.Context, filter filters.Filtered[T]) error
	DeleteMany(ctx context.Context, filter filters.Filtered[T]) error
}

type BulkUpdateOp[T any] struct {
	Filter filters.Filtered[T]
	Update T
	Insert  bool
	Upsert bool
}

type AggregateStore[T, U any] interface {
	Aggregate(ctx context.Context, pipeline []Stage) ([]U, error)
}
