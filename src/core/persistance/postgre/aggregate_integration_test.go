package postgre_test

import (
	"context"
	"testing"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/filters"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/postgre/aggregate"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/persistance/store"
)

func TestAggregateGroupWithHavingAndSort(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	_, err := tbl.BulkInsert(ctx, []product{
		{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true},
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 8, Active: true},
		{ID: "p4", Name: "Doohickey", Price: 4.99, Stock: 7, Active: true},
	})
	if err != nil {
		t.Fatalf("seed BulkInsert: %v", err)
	}

	agg := aggregate.New[product, stockSummary](pool, productsTable)
	productStock := filters.Field[product, int]("stock")
	itemCount := filters.Field[product, int64]("item_count")
	totalQty := filters.Field[product, int64]("total_qty")

	pipeline := []store.Stage{
		store.Match(productActive.Is(filters.Eq(true))),
		store.Group([]string{productActive.String()}, store.Count(itemCount.String()), store.Sum(productStock.String(), totalQty.String())),
		store.Match(totalQty.Is(filters.Gt(int64(0)))),
		store.Sort(store.Desc(totalQty.String())),
		store.Limit(5),
	}

	results, err := agg.Aggregate(ctx, pipeline)
	if err != nil {
		t.Fatalf("Aggregate: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("expected exactly 1 group, got %d: %+v", len(results), results)
	}
	got := results[0]
	if !got.Active || got.Count != 3 || got.TotalQty != 25 {
		t.Fatalf("unexpected aggregate result: %+v", got)
	}
}

func TestAggregateJoin(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	tbl := postgre.New[product](pool, productsTable)

	_, err := tbl.BulkInsert(ctx, []product{
		{ID: "p1", Name: "Widget", Price: 9.99, Stock: 10, Active: true},
		{ID: "p2", Name: "Gadget", Price: 19.99, Stock: 8, Active: true},
	})
	if err != nil {
		t.Fatalf("seed products: %v", err)
	}

	orderStore := postgre.New[order](pool, ordersTable)
	_, err = orderStore.BulkInsert(ctx, []order{
		{OrderID: "o1", ProductID: "p1", Quantity: 3},
		{OrderID: "o2", ProductID: "p1", Quantity: 2},
		{OrderID: "o3", ProductID: "p2", Quantity: 4},
	})
	if err != nil {
		t.Fatalf("seed orders: %v", err)
	}

	orderProductID := filters.Field[order, string]("product_id")
	joinedQuantity := filters.Field[order, int]("o.quantity")
	totalOrdered := filters.Field[product, int64]("total_ordered")

	agg := aggregate.New[product, orderedTotal](pool, productsTable)
	pipeline := []store.Stage{
		store.Match(productID.Is(filters.In("p1", "p2"))),
		store.Join(ordersTable, productID.String(), orderProductID.String(), "o"),
		store.Group([]string{productID.String()}, store.Sum(joinedQuantity.String(), totalOrdered.String())),
	}
	results, err := agg.Aggregate(ctx, pipeline)
	if err != nil {
		t.Fatalf("Aggregate join: %v", err)
	}

	totals := map[string]int64{}
	for _, r := range results {
		totals[r.ID] = r.TotalOrdered
	}
	if totals["p1"] != 5 {
		t.Fatalf("expected p1 total_ordered=5, got %d", totals["p1"])
	}
	if totals["p2"] != 4 {
		t.Fatalf("expected p2 total_ordered=4, got %d", totals["p2"])
	}
}
