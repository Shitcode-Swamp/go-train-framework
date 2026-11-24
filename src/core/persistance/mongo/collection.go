package mongo

import (
	"context"
	"errors"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/mongo/internal"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type Collection[T any] struct {
	collection *mongo.Collection
}

func New[T any](collection *mongo.Collection) store.Store[T] {
	return &Collection[T]{
		collection: collection,
	}
}

func (c *Collection[T]) Insert(ctx context.Context, doc T, options ...store.InsertOptions) (any, error) {
	res, err := c.collection.InsertOne(ctx, doc, internal.ToInsertMongoOpt(options)...)
	if err != nil {
		return nil, internal.FromMongoErr(err)
	}
	return res.InsertedID, nil
}

func (c *Collection[T]) BulkInsert(ctx context.Context, docs []T, options ...store.InsertOptions) ([]any, error) {
	results, err := c.collection.InsertMany(ctx, docs, internal.ToInsertManyMongoOpt(options)...)
	if err != nil {
		return nil, internal.FromMongoErr(err)
	}
	return results.InsertedIDs, nil
}

func (c *Collection[T]) Update(ctx context.Context, filter filters.Filtered[T], update T, options ...store.UpdateOptions) (any, error) {
	return c.replace(ctx, filter, update, internal.UpsertRequested(options))
}

func (c *Collection[T]) BulkWrite(ctx context.Context, ops []store.BulkUpdateOp[T]) ([]any, error) {
	ids := make([]any, len(ops))
	for i, op := range ops {
		id, err := c.write(ctx, op)
		if err != nil {
			return nil, err
		}
		ids[i] = id
	}
	return ids, nil
}

func (c *Collection[T]) write(ctx context.Context, op store.BulkUpdateOp[T]) (any, error) {
	if op.Insert {
		return c.Insert(ctx, op.Update)
	}
	return c.replace(ctx, op.Filter, op.Update, op.Upsert)
}

func (c *Collection[T]) replace(ctx context.Context, filter filters.Filtered[T], update T, upsert bool) (any, error) {
	res := c.collection.FindOneAndReplace(ctx, internal.ToMongoFilter(filter), update, internal.ToReplaceMongoOpt(upsert))
	var replaced internal.IDOnly
	if err := res.Decode(&replaced); err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return nil, nil
		}
		return nil, internal.FromMongoErr(err)
	}
	return replaced.ID, nil
}

func (c *Collection[T]) One(ctx context.Context, filter filters.Filtered[T]) (T, error) {
	res := c.collection.FindOne(ctx, internal.ToMongoFilter(filter))
	var doc T
	err := res.Err()
	if err != nil {
		return doc, internal.FromMongoErr(err)
	}
	if err = res.Decode(&doc); err != nil {
		return doc, internal.FromMongoErr(err)
	}
	return doc, nil
}

func (c *Collection[T]) Many(ctx context.Context, filter filters.Filtered[T]) ([]T, error) {
	cursor, err := c.collection.Find(ctx, internal.ToMongoFilter(filter))
	if err != nil {
		return nil, internal.FromMongoErr(err)
	}
	docs := make([]T, 0)
	if err = cursor.All(ctx, &docs); err != nil {
		return nil, internal.FromMongoErr(err)
	}
	return docs, nil
}

func (c *Collection[T]) DeleteOne(ctx context.Context, filter filters.Filtered[T]) error {
	_, err := c.collection.DeleteOne(ctx, internal.ToMongoFilter(filter))
	if err != nil {
		return internal.FromMongoErr(err)
	}
	return nil
}

func (c *Collection[T]) DeleteMany(ctx context.Context, filter filters.Filtered[T]) error {
	_, err := c.collection.DeleteMany(ctx, internal.ToMongoFilter(filter))
	if err != nil {
		return internal.FromMongoErr(err)
	}
	return nil
}
