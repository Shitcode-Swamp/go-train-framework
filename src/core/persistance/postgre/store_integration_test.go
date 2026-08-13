package postgre_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

func TestInsertAndOne(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	_, err := tbl.Insert(ctx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}

	got, err := tbl.One(ctx, productID.Is(filters.Eq("p1")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Name != "Widget" || got.Stock != 10 {
		t.Fatalf("unexpected product: %+v", got)
	}
}

func TestInsertDuplicateKey(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	if _, err := tbl.Insert(ctx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	_, err := tbl.Insert(ctx, product{ID: "p1", Name: "Widget2", Price: 1, Stock: 1, Active: true})
	if !errors.Is(err, store.ErrDuplicateKey) {
		t.Fatalf("expected ErrDuplicateKey, got %v", err)
	}
}

func TestOneNotFound(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	_, err := tbl.One(ctx, productID.Is(filters.Eq("missing")))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBulkInsertAndMany(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	_, err := tbl.BulkInsert(ctx, []product{
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true},
		{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 0, Active: false},
	})
	if err != nil {
		t.Fatalf("BulkInsert: %v", err)
	}

	got, err := tbl.Many(ctx, productActive.Is(filters.Eq(true)))
	if err != nil {
		t.Fatalf("Many: %v", err)
	}
	if len(got) != 1 || got[0].ID != "p2" {
		t.Fatalf("expected only p2 active, got %+v", got)
	}
}

func TestUpdate(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	if _, err := tbl.Insert(ctx, product{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	err := tbl.Update(ctx, productID.Is(filters.Eq("p2")), product{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 3, Active: true})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	got, err := tbl.One(ctx, productID.Is(filters.Eq("p2")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Stock != 3 {
		t.Fatalf("expected stock 3, got %d", got.Stock)
	}
}

func TestUpdateWithUpsertFallback(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	upsert := true
	err := tbl.Update(ctx, productID.Is(filters.Eq("p4")),
		product{ID: "p4", Name: "Doohickey", Price: 4.99, Stock: 7, Active: true},
		store.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		t.Fatalf("Update with upsert: %v", err)
	}
	got, err := tbl.One(ctx, productID.Is(filters.Eq("p4")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Name != "Doohickey" {
		t.Fatalf("expected upserted row, got %+v", got)
	}
}

func TestBulkUpdate(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	_, err := tbl.BulkInsert(ctx, []product{
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true},
		{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 0, Active: false},
	})
	if err != nil {
		t.Fatalf("seed BulkInsert: %v", err)
	}

	ops := []store.BulkUpdateOp[product]{
		{Filter: productID.Is(filters.Eq("p2")), Update: product{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 8, Active: true}},
		{Filter: productID.Is(filters.Eq("p3")), Update: product{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 2, Active: true}},
	}
	if err := tbl.BulkUpdate(ctx, ops); err != nil {
		t.Fatalf("BulkUpdate: %v", err)
	}

	p2, err := tbl.One(ctx, productID.Is(filters.Eq("p2")))
	if err != nil {
		t.Fatalf("One p2: %v", err)
	}
	if p2.Stock != 8 {
		t.Fatalf("expected p2 stock 8, got %d", p2.Stock)
	}

	p3, err := tbl.One(ctx, productID.Is(filters.Eq("p3")))
	if err != nil {
		t.Fatalf("One p3: %v", err)
	}
	if p3.Stock != 2 || !p3.Active {
		t.Fatalf("expected p3 stock 2 active true, got %+v", p3)
	}
}

func TestDelete(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	if _, err := tbl.Insert(ctx, product{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 2, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	if err := tbl.Delete(ctx, productID.Is(filters.Eq("p3"))); err != nil {
		t.Fatalf("Delete: %v", err)
	}
	_, err := tbl.One(ctx, productID.Is(filters.Eq("p3")))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}
