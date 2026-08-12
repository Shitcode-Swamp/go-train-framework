package contract

import (
	"context"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/filters"
)

type InsertOptions struct{}
type UpdateOptions struct {
	Upsert *bool
}

type Store[T any] interface {
	Insert(ctx context.Context, doc T, options ...InsertOptions) (id any, err error)
	BulkInsert(ctx context.Context, docs []T, options ...InsertOptions) (ids []any, err error)

	Update(ctx context.Context, filter filters.Filtered[T], update T, options ...UpdateOptions) error
	BulkUpdate(ctx context.Context, ops []BulkUpdateOp[T]) error

	One(ctx context.Context, filter filters.Filtered[T]) (T, error)
	Many(ctx context.Context, filter filters.Filtered[T]) ([]T, error)

	Delete(ctx context.Context, filter filters.Filtered[T]) error
}

type BulkUpdateOp[T any] struct {
	Filter filters.Filtered[T]
	Update T
	Upsert bool
}

type AggregateStore[T, U any] interface {
	Aggregate(ctx context.Context, pipeline []Stage) ([]U, error)
}
