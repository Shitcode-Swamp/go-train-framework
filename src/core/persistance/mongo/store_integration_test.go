package mongo_test

import (
	"context"
	"errors"
	"testing"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/mongo"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/store"
)

func TestInsertAndOne(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	_, err := products.Insert(ctx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true})
	if err != nil {
		t.Fatalf("Insert: %v", err)
	}

	got, err := products.One(ctx, productID.Is(filters.Eq("p1")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Name != "Widget" || got.Stock != 10 {
		t.Fatalf("unexpected product: %+v", got)
	}
}

func TestInsertDuplicateKey(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	if _, err := products.Insert(ctx, product{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	_, err := products.Insert(ctx, product{ID: "p1", Name: "Widget2", Price: 1, Stock: 1, Active: true})
	if !errors.Is(err, store.ErrDuplicateKey) {
		t.Fatalf("expected ErrDuplicateKey, got %v", err)
	}
}

func TestOneNotFound(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	_, err := products.One(ctx, productID.Is(filters.Eq("missing")))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound, got %v", err)
	}
}

func TestBulkInsertAndMany(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	_, err := products.BulkInsert(ctx, []product{
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true},
		{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 0, Active: false},
	})
	if err != nil {
		t.Fatalf("BulkInsert: %v", err)
	}

	got, err := products.Many(ctx, productActive.Is(filters.Eq(true)))
	if err != nil {
		t.Fatalf("Many: %v", err)
	}
	if len(got) != 1 || got[0].ID != "p2" {
		t.Fatalf("expected only p2 active, got %+v", got)
	}
}

func TestUpdate(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	if _, err := products.Insert(ctx, product{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 5, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	id, err := products.Update(ctx, productID.Is(filters.Eq("p2")), product{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 3, Active: true})
	if err != nil {
		t.Fatalf("Update: %v", err)
	}
	if id != "p2" {
		t.Fatalf("expected returned id p2, got %v", id)
	}
	got, err := products.One(ctx, productID.Is(filters.Eq("p2")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Stock != 3 {
		t.Fatalf("expected stock 3, got %d", got.Stock)
	}
}

func TestUpdateWithUpsertFallback(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	upsert := true
	id, err := products.Update(ctx, productID.Is(filters.Eq("p4")),
		product{ID: "p4", Name: "Doohickey", Price: 4.99, Stock: 7, Active: true},
		store.UpdateOptions{Upsert: &upsert},
	)
	if err != nil {
		t.Fatalf("Update with upsert: %v", err)
	}
	if id != "p4" {
		t.Fatalf("expected returned id p4, got %v", id)
	}
	got, err := products.One(ctx, productID.Is(filters.Eq("p4")))
	if err != nil {
		t.Fatalf("One: %v", err)
	}
	if got.Name != "Doohickey" {
		t.Fatalf("expected upserted document, got %+v", got)
	}
}

func TestBulkUpdate(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	_, err := products.BulkInsert(ctx, []product{
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
	ids, err := products.BulkWrite(ctx, ops)
	if err != nil {
		t.Fatalf("BulkUpdate: %v", err)
	}
	if len(ids) != 2 || ids[0] != "p2" || ids[1] != "p3" {
		t.Fatalf("expected returned ids [p2 p3], got %v", ids)
	}

	p2, err := products.One(ctx, productID.Is(filters.Eq("p2")))
	if err != nil {
		t.Fatalf("One p2: %v", err)
	}
	if p2.Stock != 8 {
		t.Fatalf("expected p2 stock 8, got %d", p2.Stock)
	}

	p3, err := products.One(ctx, productID.Is(filters.Eq("p3")))
	if err != nil {
		t.Fatalf("One p3: %v", err)
	}
	if p3.Stock != 2 || !p3.Active {
		t.Fatalf("expected p3 stock 2 active true, got %+v", p3)
	}
}

func TestDeleteOne(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	if _, err := products.Insert(ctx, product{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 2, Active: true}); err != nil {
		t.Fatalf("seed Insert: %v", err)
	}

	if err := products.DeleteOne(ctx, productID.Is(filters.Eq("p3"))); err != nil {
		t.Fatalf("DeleteOne: %v", err)
	}
	_, err := products.One(ctx, productID.Is(filters.Eq("p3")))
	if !errors.Is(err, store.ErrNotFound) {
		t.Fatalf("expected ErrNotFound after delete, got %v", err)
	}
}

func TestDeleteOneLimitsToSingleDocument(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	if _, err := products.BulkInsert(ctx, []product{
		{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true},
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 8, Active: true},
	}); err != nil {
		t.Fatalf("seed BulkInsert: %v", err)
	}

	if err := products.DeleteOne(ctx, productActive.Is(filters.Eq(true))); err != nil {
		t.Fatalf("DeleteOne: %v", err)
	}

	got, err := products.Many(ctx, productActive.Is(filters.Eq(true)))
	if err != nil {
		t.Fatalf("Many: %v", err)
	}
	if len(got) != 1 {
		t.Fatalf("expected exactly 1 document left after DeleteOne, got %d: %+v", len(got), got)
	}
}

func TestDeleteMany(t *testing.T) {
	coll := requireCollection(t)
	ctx := context.Background()
	resetCollection(t, ctx, coll)
	products := mongo.New[product](coll)

	if _, err := products.BulkInsert(ctx, []product{
		{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true},
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 8, Active: true},
		{ID: "p3", Name: "Gizmo", Price: 29.99, Stock: 0, Active: false},
	}); err != nil {
		t.Fatalf("seed BulkInsert: %v", err)
	}

	if err := products.DeleteMany(ctx, productActive.Is(filters.Eq(true))); err != nil {
		t.Fatalf("DeleteMany: %v", err)
	}

	got, err := products.Many(ctx, productActive.Is(filters.Eq(true)))
	if err != nil {
		t.Fatalf("Many: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("expected no active documents left after DeleteMany, got %d: %+v", len(got), got)
	}

	remaining, err := products.One(ctx, productID.Is(filters.Eq("p3")))
	if err != nil {
		t.Fatalf("One p3: %v", err)
	}
	if remaining.ID != "p3" {
		t.Fatalf("expected p3 to remain untouched, got %+v", remaining)
	}
}
