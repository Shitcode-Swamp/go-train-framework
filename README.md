# go-train-framework

![gopher-train](./gopher-train.png)

[![Latest release](https://img.shields.io/github/v/release/Shitcode-Swamp/go-train-framework)](https://github.com/Shitcode-Swamp/go-train-framework/releases/latest)

[![Go version](https://img.shields.io/github/go-mod/go-version/Shitcode-Swamp/go-train-framework)](go.mod)

Custom framework made in educational puposes.

## Roadmap

| Status     | Feature              | Description                                                                                |
| ---------- | -------------------- | ------------------------------------------------------------------------------------------ |
| ✅ Done    | Dependency Injection | Auto-wires `AppContext` fields, constructors, and getters from `@Injection`/`@Inject` tags |
|  ✅ Done | Controllers codegen  | Controller abstraction + docs                                            |
| ✅ Done    | gRPC codegen         | Generate gRPC service stubs from a shared spec repo                                        |
| ✅ Done    | DB abstraction       | `Store[T]`/`AggregateStore[T,U]` over Postgres and MongoDB, shared filter DSL, transactions |
| 🔜 Planned | Message broker and MQ abstractions      |  Kafka and RabitMQ |

## Overall description

Framework automaticaly builds AppCtx with needed dependencies
To use the framework build tool you might run binary via
`chmod +x ./build/go-train && ./build/go-train` or `run.sh`

You must describe what directory to watch via flag `--root=path`

How the framework is built is descrived in the `structure.excalidraw`

## DI tags

Two comment annotations drive the generator:

- `@Injection(token)` — goes on a constructor function's doc comment. Marks it
  as a dependency to include in `AppContext`, keyed by `token`.
- `@Inject(token)` — goes on a constructor parameter (as an inline `/* ... */`
  comment). Marks that parameter as needing another `@Injection`-tagged
  dependency, resolved by its `token`.

```go
// @Injection(testInjectee)
func InitInjectee() *TestStruct {
	return &TestStruct{}
}

// @Injection(test)
func Init(
	testField interface{} /* @Inject(testInjectee) */,
) *TestStruct {
	return &TestStruct{
		TestField: testField,
	}
}
```

Running the tool over a directory containing both produces something like:

```go
type AppContext struct {
	testInjectee *pkg.TestStruct
	test         *pkg.TestStruct
}

func NewAppContext() *AppContext {
	testInjectee := pkg.InitInjectee()
	test := pkg.Init(testInjectee)
	return &AppContext{
		testInjectee: testInjectee,
		test:         test,
	}
}
```

Dependency order is resolved automatically — you don't need to declare `Init`
before `InitInjectee`, or worry about wiring them together by hand.

### Flags

`@Injection(token, flag-one, flag-two)` accepts optional comma-separated
flags after the token to opt out of parts of the generated code for that
dependency:

- `no-getter` — don't generate an `AppContext` getter method for it. The
  field is still populated and wired into the constructor, it's just not
  exposed via a getter.

```go
// @Injection(internalCache, no-getter)
func NewCache() *Cache {
	return &Cache{}
}
```

## gRPC codegen

`go-train-grpc` reads the `grpcConfig` section of your config file and turns
a shared repo of `.proto` files into generated Go/gRPC code:

```yaml
grpcConfig:
  output: ./out/grpc
  sourceRegistry:
    sourceRegistry: your-org/proto-specs
    isGithub: true
```

- `sourceRegistry.sourceRegistry` — where the `.proto` files live. With
  `isGithub: true` this is an `owner/repo` shorthand resolved against
  `github.com`; otherwise it's used as-is, so it can be any git remote
  (`https://`, `ssh://`) or a local path.
- `sourceRegistry.isGithub` — toggles the `owner/repo` shorthand above.
- `output` — where the generated `.go` files are written.

It clones the registry, then shells out to `protoc` with the `go` and
`go-grpc` plugins, so `protoc`, `protoc-gen-go`, and `protoc-gen-go-grpc`
must be on your `PATH`.

```sh
go install github.com/Shitcode-Swamp/go-train-framework/v3/cmd/go-train-grpc@latest
go-train-grpc -config=./go-train.yaml
```

For a private registry over `https://`, set `GO_TRAIN_GRPC_TOKEN` to a
personal access token:

```sh
GO_TRAIN_GRPC_TOKEN=<token> go-train-grpc -config=./go-train.yaml
```

## DB abstraction

`src/core/persistance` gives you one `Store[T]` / `AggregateStore[T,U]` API
that runs against either Postgres or MongoDB. Application code is written
against the interfaces in the `store` package; which backend is actually
behind it is decided once, at construction time.

```text
src/core/persistance/
├── store/     contracts: Store[T], AggregateStore[T,U], Transactor, Stage, sentinel errors
├── filters/   the Field[O,V]/Filter[V]/Filtered[O] DSL used to build filters
├── postgre/   Postgres adapter (pgx/v5): Table[T], Database
│   └── aggregate/   Pipeline[T,U], Postgres-only AggregateStore[T,U]
└── mongo/     MongoDB adapter (mongo-driver/v2): Collection[T], Database
```

### Defining an entity

```go
type Product struct {
	ID     string  `store:"id,pk"`
	Name   string  `store:"name"`
	Price  float64 `store:"price"`
	Stock  int     `store:"stock"`
	Active bool    `store:"active"`
}
```

The `store:"column[,pk]"` tag drives the Postgres adapter's column mapping
(reflection-based today; a codegen step that emits typed field constants
from these tags is planned but not built yet). The Mongo adapter marshals
structs with the driver's own `bson` rules — add `bson:"..."` tags if you
need different field names there.

### Constructing a store

```go
// Postgres
pool, _ := pgxpool.New(ctx, dsn)
products := postgre.New[Product](pool, "products") // store.Store[Product]

// MongoDB - import our adapter aliased, since you also need the driver's
// own "mongo" package (go.mongodb.org/mongo-driver/v2/mongo) to build the
// *mongo.Client/*mongo.Collection in the first place:
//   mongodb "github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/mongo"
client, _ := mongo.Connect(options.Client().ApplyURI(uri))
collection := client.Database("shop").Collection("products")
products := mongodb.New[Product](collection) // store.Store[Product]
```

Both return a `store.Store[T]`, so code written against that interface
doesn't know or care which backend it's talking to.

### Filters

Build filters with `filters.Field[O, V]` — `O` is the entity type, `V` the
field's value type:

```go
id     := filters.Field[Product, string]("id")
active := filters.Field[Product, bool]("active")
stock  := filters.Field[Product, int]("stock")

filter := id.Is(filters.Eq("p1"))
filter = active.Is(filters.Eq(true)).And(stock.Is(filters.Gt(0)))
filter = id.Is(filters.Not(filters.In("p1", "p2")))
```

Operators: `Eq`, `Ne`, `Gt`, `Gte`, `Lt`, `Lte`, `In`, `Not`. `.And`/`.Or`
live on the resulting `Filtered[O]` (not on `Filter[V]`) since combining
conditions across different fields needs to erase each field's own value
type — a range on one field is just `Gt(...).And(Lt(...))` on the same
`Field`. A zero-value `filters.Filtered[T]{}` matches everything, which is
what you want for "all rows" in `Many`/`DeleteMany`.

### CRUD

```go
id, err := products.Insert(ctx, Product{ID: "p1", Name: "Widget", Stock: 10, Active: true})
ids, err := products.BulkInsert(ctx, []Product{ /* ... */ })

_, err = products.Update(ctx, id.Is(filters.Eq("p1")), Product{ /* ... */ })

upsert := true
_, err = products.Update(ctx, id.Is(filters.Eq("p4")), Product{ /* ... */ },
	store.UpdateOptions{Upsert: &upsert})

_, err = products.BulkWrite(ctx, []store.BulkUpdateOp[Product]{
	{Filter: id.Is(filters.Eq("p1")), Update: Product{ /* ... */ }},
	{Insert: true, Update: Product{ID: "p2" /* ... */}},
})

one, err := products.One(ctx, id.Is(filters.Eq("p1")))
many, err := products.Many(ctx, active.Is(filters.Eq(true)))

err = products.DeleteOne(ctx, id.Is(filters.Eq("p1")))
err = products.DeleteMany(ctx, active.Is(filters.Eq(false)))
```

### Errors

Both adapters translate driver-specific errors into one shared sentinel
set in `store`: `ErrNotFound`, `ErrDuplicateKey`, `ErrConstraintViolation`,
`ErrConflict`, `ErrInvalidFilter`, `ErrInvalidDocument`,
`ErrUnsupportedOperation`, `ErrConnection`.

```go
_, err := products.One(ctx, filter)
if errors.Is(err, store.ErrNotFound) {
	// ...
}
```

The original driver error stays in the chain too, if you need the
backend-specific code:

```go
var pgErr *pgconn.PgError
errors.As(err, &pgErr) // pgErr.Code == "23505"
```

### Aggregation pipelines

`AggregateStore[T, U]` runs a `[]store.Stage` pipeline, where `T` is the
source entity and `U` is the shape results get scanned into. Only the
Postgres adapter implements it today, via `postgre/aggregate`:

```go
import pgaggregate "github.com/Shitcode-Swamp/go-train-framework/v3/src/core/persistance/postgre/aggregate"

type StockSummary struct {
	Active   bool  `store:"active"`
	Count    int64 `store:"item_count"`
	TotalQty int64 `store:"total_qty"`
}

agg := pgaggregate.New[Product, StockSummary](pool, "products")

totalQty := filters.Field[Product, int64]("total_qty")
pipeline := []store.Stage{
	store.Match(active.Is(filters.Eq(true))),
	store.Group([]string{"active"}, store.Count("item_count"), store.Sum("stock", "total_qty")),
	store.Match(totalQty.Is(filters.Gt(int64(0)))), // HAVING, since it's after Group
	store.Sort(store.Desc("total_qty")),
	store.Limit(5),
}

results, err := agg.Aggregate(ctx, pipeline)
```

Note the package itself is named `postgre_aggregate` (its import path is
`.../postgre/aggregate`), so either alias the import as above or write
`postgre_aggregate.New[...]`.

Stages: `Match` (filter — `WHERE` before a `Group`, `HAVING` after),
`Convert` (cast a field), `Group` + `Sum`/`Avg`/`Min`/`Max`/`Count`,
`Sort` + `Asc`/`Desc`, `Limit`, `Skip`, `Join`, `Project` + `As`.

Two things worth knowing before you reach for `Join` or `Match`:

- `store.Match` always needs `filters.Field[T, V]` typed against the
  pipeline's *source* entity `T` — even when filtering an aggregate's own
  output, like `total_qty` above. The generic type is only there so the
  pipeline's stage dispatch recognizes it; the field name is just a string.
- On Postgres, `Join` compiles to a plain SQL `LEFT JOIN`, which multiplies
  rows on a 1:many match. That's different from MongoDB's `$lookup`, which
  nests matches into an array — but since Mongo has no aggregation adapter
  here yet, there's nothing to be inconsistent with in practice.

### Transactions

```go
db := postgre.NewDatabase(pool) // or mongodb.NewDatabase(client)

err := db.WithTransaction(ctx, func(txCtx context.Context) error {
	if _, err := products.Insert(txCtx, Product{ /* ... */ }); err != nil {
		return err
	}
	return orders.Insert(txCtx, Order{ /* ... */ })
})
```

Keep using the exact same `Store[T]` instances you already built
(`products`, `orders` above) — they detect the ambient transaction from
`txCtx` and route through it automatically. There's no separate
"transactional store" type to construct.

On Postgres you can pick the isolation level per call:

```go
err := db.WithTransaction(ctx, fn, store.TransactionOptions{
	IsolationLevel: store.IsolationSerializable,
})
```

Default (no options) is `READ COMMITTED`, matching plain Postgres. That's
fine for atomicity (all-or-nothing), but **not** safe for "read a value,
compute a new value, write it back" under concurrency — two concurrent
transactions can both read the same value and both write back the same
result, silently losing an update. Use `IsolationSerializable` for that
pattern, and retry on `store.ErrConflict`:

```go
for {
	err := db.WithTransaction(ctx, func(txCtx context.Context) error {
		current, err := products.One(txCtx, id.Is(filters.Eq("p1")))
		if err != nil {
			return err
		}
		current.Stock--
		_, err = products.Update(txCtx, id.Is(filters.Eq("p1")), current)
		return err
	}, store.TransactionOptions{IsolationLevel: store.IsolationSerializable})
	if err == nil {
		break
	}
	if !errors.Is(err, store.ErrConflict) {
		return err // a real error, not a retryable conflict
	}
	// retry
}
```

On MongoDB, `WithTransaction` requires a replica set — a standalone
`mongod` doesn't support transactions at all — and `IsolationLevel` is
currently a no-op there: Mongo transactions are always snapshot-isolated
by the driver's own session mechanics, so there's no equivalent knob to
forward yet.

## Using it as a Go tool in other repos

The command lives at `cmd/go-train` so it can be installed like any other Go tool.

**Global install** (puts `go-train` on your `$GOBIN`/`$GOPATH/bin`):

```sh
go install github.com/Shitcode-Swamp/go-train-framework/v3/cmd/go-train@latest
go-train -root=./path/to/watch -ctxPath=./path/to/output
```

**As a per-module tool dependency** (Go 1.24+, pins the version in the consuming
repo's `go.mod` via a `tool` directive):

```sh
go get -tool github.com/Shitcode-Swamp/go-train-framework/v3/cmd/go-train@latest
go tool go-train -root=./path/to/watch -ctxPath=./path/to/output
```

While developing locally without a published module, point a consuming repo's
`go.mod` at your local checkout instead:

```go
replace github.com/Shitcode-Swamp/go-train-framework/v3 => /path/to/go-train-framework
```
