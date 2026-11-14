package mongo

import (
	"context"

	"go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/mongo/internal"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

type Database struct {
	client *mongo.Client
}

func NewDatabase(client *mongo.Client) store.Transactor {
	return &Database{client: client}
}

func (d *Database) WithTransaction(ctx context.Context, fn func(ctx context.Context) error, options ...store.TransactionOptions) error {
	session, err := d.client.StartSession()
	if err != nil {
		return internal.FromMongoErr(err)
	}
	defer session.EndSession(ctx)

	_, err = session.WithTransaction(ctx, func(sessCtx context.Context) (any, error) {
		return nil, fn(sessCtx)
	})
	return internal.FromMongoErr(err)
}

var _ store.Transactor = (*Database)(nil)
