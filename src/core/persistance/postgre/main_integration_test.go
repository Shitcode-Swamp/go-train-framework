package postgre_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	tcpostgres "github.com/testcontainers/testcontainers-go/modules/postgres"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
)

const (
	productsTable = "products"
	ordersTable   = "orders"
)

type product struct {
	ID     string  `store:"id,pk"`
	Name   string  `store:"name"`
	Price  float64 `store:"price"`
	Stock  int     `store:"stock"`
	Active bool    `store:"active"`
}

type stockSummary struct {
	Active   bool  `store:"active"`
	Count    int64 `store:"item_count"`
	TotalQty int64 `store:"total_qty"`
}

type order struct {
	OrderID   string `store:"order_id,pk"`
	ProductID string `store:"product_id"`
	Quantity  int    `store:"quantity"`
}

type orderedTotal struct {
	ID           string `store:"id,pk"`
	TotalOrdered int64  `store:"total_ordered"`
}

var (
	productID     = filters.Field[product, string]("id")
	productActive = filters.Field[product, bool]("active")
)

var (
	testPool     *pgxpool.Pool
	containerErr error
)

func TestMain(m *testing.M) {
	os.Exit(runIntegrationTests(m))
}

func runIntegrationTests(m *testing.M) int {
	ctx := context.Background()

	container, err := tcpostgres.Run(ctx, "postgres:16-alpine",
		tcpostgres.WithDatabase("storetest"),
		tcpostgres.WithUsername("postgres"),
		tcpostgres.WithPassword("postgres"),
		tcpostgres.BasicWaitStrategies(),
	)
	if err != nil {
		containerErr = err
		return m.Run()
	}
	defer func() { _ = container.Terminate(ctx) }()

	dsn, err := container.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		containerErr = err
		return m.Run()
	}

	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		containerErr = err
		return m.Run()
	}
	defer pool.Close()

	schema := fmt.Sprintf(`
CREATE TABLE %[1]s (
	id     TEXT PRIMARY KEY,
	name   TEXT NOT NULL,
	price  DOUBLE PRECISION NOT NULL,
	stock  INT NOT NULL,
	active BOOLEAN NOT NULL
);
CREATE TABLE %[2]s (
	order_id   TEXT PRIMARY KEY,
	product_id TEXT NOT NULL REFERENCES %[1]s(id),
	quantity   INT NOT NULL
);
`, productsTable, ordersTable)
	if _, err := pool.Exec(ctx, schema); err != nil {
		containerErr = err
		return m.Run()
	}

	testPool = pool
	return m.Run()
}

func requirePool(t *testing.T) *pgxpool.Pool {
	t.Helper()
	if containerErr != nil {
		t.Skipf("postgres testcontainer unavailable: %v", containerErr)
	}
	return testPool
}

func resetSchema(t *testing.T, ctx context.Context, pool *pgxpool.Pool) {
	t.Helper()
	if _, err := pool.Exec(ctx, fmt.Sprintf("TRUNCATE %s, %s", ordersTable, productsTable)); err != nil {
		t.Fatalf("truncate: %v", err)
	}
}
