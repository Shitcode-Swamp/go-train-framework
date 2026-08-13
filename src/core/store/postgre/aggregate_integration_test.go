package postgre_test

import (
	"context"
	"testing"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/contract"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/filters"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/postgre"
	"github.com/Shitcode-Swamp/go-train-framework/src/core/store/postgre/aggregate"
)

func TestAggregateGroupWithHavingAndSort(t *testing.T) {
	pool := requirePool(t)
	ctx := context.Background()
	resetSchema(t, ctx, pool)
	store := postgre.New[product](pool, productsTable)

	_, err := store.BulkInsert(ctx, []product{
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

	pipeline := []contract.Stage{
		contract.Match(productActive.Is(filters.Eq(true))),
		contract.Group([]string{productActive.String()}, contract.Count(itemCount.String()), contract.Sum(productStock.String(), totalQty.String())),
		contract.Match(totalQty.Is(filters.Gt(int64(0)))),
		contract.Sort(contract.Desc(totalQty.String())),
		contract.Limit(5),
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
	store := postgre.New[product](pool, productsTable)

	_, err := store.BulkInsert(ctx, []product{
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
	pipeline := []contract.Stage{
		contract.Match(productID.Is(filters.In("p1", "p2"))),
		contract.Join(ordersTable, productID.String(), orderProductID.String(), "o"),
		contract.Group([]string{productID.String()}, contract.Sum(joinedQuantity.String(), totalOrdered.String())),
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
