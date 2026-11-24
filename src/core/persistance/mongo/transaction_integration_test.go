package mongo_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/mongo"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
)

func TestWithTransactionCommit(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	productStore := mongo.New[product](coll)
	db := mongo.NewDatabase(testClient)

	err := db.WithTransaction(ctx, func(txCtx context.Context) error {
		if _, err := productStore.Insert(txCtx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true}); err != nil {
			return err
		}
		if _, err := productStore.Insert(txCtx, product{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true}); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithTransaction: %v", err)
	}

	got, err := productStore.Many(ctx, productActive.Is(filters.Eq(true)))
	if err != nil {
		t.Fatalf("Many: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected both inserts to be committed, got %d documents: %+v", len(got), got)
	}
}

func TestWithTransactionRollback(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	productStore := mongo.New[product](coll)
	db := mongo.NewDatabase(testClient)

	sentinel := errors.New("boom")
	err := db.WithTransaction(ctx, func(txCtx context.Context) error {
		if _, err := productStore.Insert(txCtx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true}); err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected WithTransaction to propagate the callback error, got %v", err)
	}

	_, err = productStore.One(ctx, productID.Is(filters.Eq("p1")))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected insert to be rolled back, got %v", err)
	}
}

func TestWithTransactionBulkInsertParticipatesInAmbientTx(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	productStore := mongo.New[product](coll)
	db := mongo.NewDatabase(testClient)

	sentinel := errors.New("boom")
	err := db.WithTransaction(ctx, func(txCtx context.Context) error {
		_, err := productStore.BulkInsert(txCtx, []product{
			{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true},
			{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true},
		})
		if err != nil {
			return err
		}
		return sentinel
	})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected WithTransaction to propagate the callback error, got %v", err)
	}

	got, err := productStore.Many(ctx, productActive.Is(filters.Eq(true)))
	if err != nil {
		t.Fatalf("Many: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected BulkInsert to roll back with the outer transaction, got %+v", got)
	}
}

func TestWithTransactionConcurrentDecrementsDoNotLoseUpdates(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	productStore := mongo.New[product](coll)
	db := mongo.NewDatabase(testClient)

	const initialStock = 50
	const workers = 10

	if _, err := productStore.Insert(ctx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: initialStock, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	var wg sync.WaitGroup
	errs := make([]error, workers)
	for i := 0; i < workers; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			errs[i] = decrementStockWithRetry(ctx, db, productStore, "p1")
		}(i)
	}
	wg.Wait()

	for i, err := range errs {
		if err != nil {
			t.Fatalf("worker %d: %v", i, err)
		}
	}

	got, err := productStore.One(ctx, productID.Is(filters.Eq("p1")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Stock != initialStock-workers {
		t.Fatalf("expected stock %d after %d concurrent decrements, got %d (lost update)", initialStock-workers, workers, got.Stock)
	}
}

func TestWithTransactionRunsInsideSession(t *testing.T) {
	requireCollection(t)
	ctx := context.Background()
	db := mongo.NewDatabase(testClient)

	err := db.WithTransaction(ctx, func(txCtx context.Context) error {
		if mongodriver.SessionFromContext(txCtx) == nil {
			return errors.New("expected an ambient session in context")
		}
		return nil
	})
	if err != nil {
		t.Fatalf("WithTransaction: %v", err)
	}
}

func decrementStockWithRetry(ctx context.Context, db store.Transactor, productStore store.Store[product], id string, options ...store.TransactionOptions) error {
	const maxAttempts = 50
	for attempt := 0; attempt < maxAttempts; attempt++ {
		err := db.WithTransaction(ctx, func(txCtx context.Context) error {
			current, err := productStore.One(txCtx, productID.Is(filters.Eq(id)))
			if err != nil {
				return err
			}
			current.Stock--
			_, err = productStore.Update(txCtx, productID.Is(filters.Eq(id)), current)
			return err
		}, options...)
		if err == nil {
			return nil
		}
		if !errors.Is(err, store.ErrConflict) {
			return err
		}
	}
	return fmt.Errorf("exceeded %d retry attempts on store.ErrConflict", maxAttempts)
}
