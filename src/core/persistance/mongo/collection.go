package mongo

import (
	"context"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/array"
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
	res, err := c.collection.InsertOne(ctx, doc, toInsertMongoOpt(options)...)
	if err != nil {
		return nil, fromMongoErr(err) 
	}
	return res.InsertedID, nil
}

func (c *Collection[T]) BulkInsert(ctx context.Context, docs []T, options ...store.InsertOptions) ([]any, error) {
	results, err := c.collection.InsertMany(ctx, docs, toInsertManyMongoOpt(options)...)
	if err != nil {
		return nil, fromMongoErr(err) 
	}
	return results.InsertedIDs , nil
}

func (c *Collection[T]) Update(ctx context.Context, filter filters.Filtered[T], update T, options ...store.UpdateOptions) (any , error) {
	res, err := c.collection.UpdateOne(ctx, filter, update, toUpdateOneMongoOpt(options)...)
	if err != nil {
		return nil, fromMongoErr(err) 
	}	
	if array.Some(options, func(e store.UpdateOptions) bool { return *e.Upsert }) {
		return res.UpsertedID, err
	}
	return nil, nil
}

func (c *Collection[T]) BulkUpdate(ctx context.Context, ops []store.BulkUpdateOp[T]) ([]any, error) {
	// Implementation for bulk updating items in the MongoDB collection
	return nil, nil
}

func (c *Collection[T]) One(ctx context.Context, filter filters.Filtered[T]) (T, error) {
	// Implementation for finding items in the MongoDB collection based on a filter
	var a T
	return a, nil
}

func (c *Collection[T]) Many(ctx context.Context, filter filters.Filtered[T]) ([]T, error) {
	// Implementation for finding items in the MongoDB collection based on a filter
	return nil, nil
}

func (c *Collection[T]) Delete(ctx context.Context, filter filters.Filtered[T]) error {
	// Implementation for deleting an item from the MongoDB collection
	return nil
}
