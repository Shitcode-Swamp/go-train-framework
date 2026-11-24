package postgre_test

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"testing"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/postgre"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/postgre/internal/pgsql"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
)

func TestWithTransactionCommit(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	productStore := postgre.New[product](pool, productsTable)
	db := postgre.NewDatabase(pool)

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
		t.Fatalf("expected both inserts to be committed, got %d rows: %+v", len(got), got)
	}
}

func TestWithTransactionRollback(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	productStore := postgre.New[product](pool, productsTable)
	db := postgre.NewDatabase(pool)

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
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	productStore := postgre.New[product](pool, productsTable)
	db := postgre.NewDatabase(pool)

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
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	productStore := postgre.New[product](pool, productsTable)
	db := postgre.NewDatabase(pool)

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
			errs[i] = decrementStockWithRetry(ctx, db, productStore, "p1", store.TransactionOptions{IsolationLevel: store.IsolationSerializable})
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

func TestWithTransactionIsolationLevelIsApplied(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	db := postgre.NewDatabase(pool)

	cases := []struct {
		name     string
		options  []store.TransactionOptions
		expected string
	}{
		{"default is read committed", nil, "read committed"},
		{"explicit serializable", []store.TransactionOptions{{IsolationLevel: store.IsolationSerializable}}, "serializable"},
		{"repeatable read", []store.TransactionOptions{{IsolationLevel: store.IsolationRepeatableRead}}, "repeatable read"},
		{"read committed", []store.TransactionOptions{{IsolationLevel: store.IsolationReadCommitted}}, "read committed"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			var got string
			err := db.WithTransaction(ctx, func(txCtx context.Context) error {
				tx, ok := pgsql.TxFromContext(txCtx)
				if !ok {
					t.Fatal("expected an ambient tx in context")
				}
				return tx.QueryRow(txCtx, "SELECT current_setting('transaction_isolation')").Scan(&got)
			}, tc.options...)
			if err != nil {
				t.Fatalf("WithTransaction: %v", err)
			}
			if got != tc.expected {
				t.Fatalf("expected isolation %q, got %q", tc.expected, got)
			}
		})
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
