package mongo_test

import (
	"context"
	"os"
	"testing"

	tcmongo "github.com/testcontainers/testcontainers-go/modules/mongodb"
	"go.mongodb.org/mongo-driver/v2/bson"
	mongodriver "go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/filters"
)

const (
	databaseName       = "storetest"
	productsCollection = "products"
)

type product struct {
	ID     string  `bson:"_id"`
	Name   string  `bson:"name"`
	Price  float64 `bson:"price"`
	Stock  int     `bson:"stock"`
	Active bool    `bson:"active"`
}

var (
	productID     = filters.Field[product, string]("_id")
	productActive = filters.Field[product, bool]("active")
)

var (
	testClient   *mongodriver.Client
	containerErr error
)

func TestMain(m *testing.M) {
	os.Exit(runIntegrationTests(m))
}

func runIntegrationTests(m *testing.M) int {
	ctx := context.Background()

	container, err := tcmongo.Run(ctx, "mongo:7", tcmongo.WithReplicaSet("rs0"))
	if err != nil {
		containerErr = err
		return m.Run()
	}
	defer func() { _ = container.Terminate(ctx) }()

	uri, err := container.ConnectionString(ctx)
	if err != nil {
		containerErr = err
		return m.Run()
	}

	// The replica set advertises the container's internal hostname, which is unreachable from the host.
	client, err := mongodriver.Connect(options.Client().ApplyURI(uri + "&directConnection=true"))
	if err != nil {
		containerErr = err
		return m.Run()
	}
	defer func() { _ = client.Disconnect(ctx) }()

	if err := client.Database(databaseName).CreateCollection(ctx, productsCollection); err != nil {
		containerErr = err
		return m.Run()
	}

	testClient = client
	return m.Run()
}

func requireCollection(t *testing.T) *mongodriver.Collection {
	t.Helper()
	if containerErr != nil {
		t.Skipf("mongo testcontainer unavailable: %v", containerErr)
	}
	return testClient.Database(databaseName).Collection(productsCollection)
}

func resetCollection(t *testing.T, ctx context.Context, coll *mongodriver.Collection) {
	t.Helper()
	if _, err := coll.DeleteMany(ctx, bson.M{}); err != nil {
		t.Fatalf("reset collection: %v", err)
	}
}
