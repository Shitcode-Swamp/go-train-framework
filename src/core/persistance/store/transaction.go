package store

import "context"

type IsolationLevel string

const (
	IsolationSerializable   IsolationLevel = "serializable"
	IsolationRepeatableRead IsolationLevel = "repeatable_read"
	IsolationReadCommitted  IsolationLevel = "read_committed"
)

type TransactionOptions struct {
	IsolationLevel IsolationLevel
}

type Transactor interface {
	WithTransaction(ctx context.Context, fn func(ctx context.Context) error, options ...TransactionOptions) error
}
