# go-train-framework — Presentation Slides

---

## Slide 1 — What is Dependency Injection?

**Dependency Injection (DI)** is a pattern where an object receives its dependencies from the outside rather than creating them itself.

```go
// Without DI — service creates its own deps (hard to test, tightly coupled)
func NewUserService() *UserService {
    db := sql.Open("postgres", "...")   // created internally
    return &UserService{db: db}
}

// With DI — deps are passed in (easy to test, loosely coupled)
func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}
```

The caller is responsible for constructing and passing dependencies. In a large app this means someone must wire everything together — that is what a DI framework automates.

---

## Slide 2 — The Wiring Problem

As the app grows, the bootstrap code becomes a maintenance burden:

```go
// You write this by hand and keep it in sync manually
logger  := infra.NewLogger()
cfg     := config.New()
db      := infra.NewDatabase(cfg, logger)
cache   := infra.NewCache(cfg)
repo    := repos.NewUserRepo(db, cache)
mailer  := mail.NewMailer(cfg, logger)
svc     := services.NewUserService(repo, mailer, logger)
handler := api.NewUserHandler(svc, logger)
```

**Problems:**
- Wrong order → nil-pointer panic at startup
- Adding a new dependency anywhere means touching bootstrap code
- Hard to see the dependency structure at a glance

---

## Slide 3 — How go-train-framework Solves This

The framework is a **build-time code generator**. You annotate your constructors with special comments; the tool reads your entire codebase, discovers all annotated constructors, figures out the correct initialization order automatically, and writes the bootstrap file for you.

```
┌─ Your annotated .go files ─┐        ┌─ Generated ──────────────┐
│                            │        │                          │
│  // @Injection(userSvc)    │        │  func NewAppCtx() *App { │
│  func NewUserService(      │  ────► │    logger := NewLogger() │
│    repo *Repo, // @Inject  │        │    db     := NewDB()     │
│  ) *UserService { ... }    │        │    repo   := NewRepo(db) │
│                            │        │    svc    := NewSvc(repo)│
└────────────────────────────┘        └──────────────────────────┘
```

The core mechanism is entirely **graph-based**.

---

## Slide 4 — Graph Nodes: What Are They?

Every constructor marked with `// @Injection(token)` becomes a **node** in the dependency graph.

```go
// @Injection(logger)        ← this constructor is a node
func NewLogger() *Logger { ... }

// @Injection(database)      ← this constructor is a node
func NewDatabase(...) *sql.DB { ... }

// @Injection(userService)   ← this constructor is a node
func NewUserService(...) *UserService { ... }
```

A node carries:
- a unique **token** (the string inside `@Injection(...)`)
- the **package path** for import generation
- the **function name** to call
- a list of **injectee tokens** (its dependencies)

Tokens are the node identities — this is how nodes find and reference each other across packages.

---

## Slide 5 — Graph Edges: What Are They?

Every constructor parameter annotated with `// @Inject(token)` becomes a **directed edge** from this node to the node identified by that token.

```go
// @Injection(userService)
func NewUserService(
    repo   *UserRepo, // @Inject(userRepo)    ← edge: userService → userRepo
    mailer *Mailer,   // @Inject(mailer)      ← edge: userService → mailer
    logger *Logger,   // @Inject(logger)      ← edge: userService → logger
) *UserService { ... }
```

An edge **A → B** means: *"node A declares that it needs node B as input"* — A depends on B.

No annotation on a parameter = the framework ignores that parameter.

---

## Slide 6 — The Initial Dependency Graph

After scanning all files, the framework has a set of nodes and edges. This is the **initial dependency graph** — a Directed Acyclic Graph (DAG) where edges point from a node to what it needs.

```
        ┌───┐
        │ 6 │  (e.g. userService — needs 1, 4, 5)
        └───┘
       /  |  \
      ▼   ▼   ▼
   ┌───┐ ┌───┐ ┌───┐
   │ 1 │ │ 4 │ │ 5 │  (e.g. logger, userRepo, mailer)
   └───┘ └───┘ └───┘
           \   /
            ▼ ▼
           ┌───┐
           │ 2 │  (e.g. database — needed by userRepo and mailer)
           └───┘

        ┌───┐
        │ 7 │  (e.g. handler — needs 3)
        └───┘
          │
          ▼
        ┌───┐
        │ 3 │  (e.g. config)
        └───┘
```

*Node numbers = final initialization order. Names in practice are tokens.*

---

## Slide 7 — Goal: Find the Initialization Order

The framework must produce an ordering of all nodes such that for every edge **A → B**, node **B** appears before node **A** in the sequence.

For the graph above, a valid order is:

```
Step 1 — 1, 2, 3    (no dependencies at all)
Step 2 — 4, 5       (need 2, which is ready)
Step 3 — 7          (needs 3, which is ready)
Step 4 — 6          (needs 1, 4, 5 — all ready)
```

This is the **topological ordering** of a DAG — a classical graph algorithm problem.

---

## Slide 8 — The Challenge: Multiple Disconnected Roots

The initial dep graph often has **multiple nodes with no incoming edges** — top-level consumers that nothing else depends on (e.g. nodes `6` and `7`). The two sub-graphs are disconnected.

```
  6 ──► 1                7 ──► 3
  6 ──► 4 ──► 2          (isolated component)
  6 ──► 5 ──►

  ↑ two disconnected roots — which one do we start from?
```

Standard BFS requires a single starting node. The framework solves this with a **two-step graph transformation**.

---

## Slide 9 — Transform Step 1: Reverse All Edges

Flip every edge in the graph. An edge **A → B** (A needs B) becomes **B → A** (B enables A).

```
  BEFORE (A → deps of A)       AFTER (reversed)

  6 ──► 1                      1 ──► 6
  6 ──► 4                      4 ──► 6
  6 ──► 5                      5 ──► 6
  4 ──► 2                      2 ──► 4
  5 ──► 2                      2 ──► 5
  7 ──► 3                      3 ──► 7
```

After reversal:
- An edge **A → B** now means *"initialize A before B"*
- Nodes **1, 2, 3** have no incoming edges — they are pure providers with no dependencies of their own

---

## Slide 10 — Transform Step 2: Add a Ghost Root Node

Create a virtual **root** node and draw an edge from it to every node that has **no incoming edges** in the reversed graph (i.e. the pure providers found in step 1).

```
              root  (virtual — never emitted to code)
            /  |  \
           ▼   ▼   ▼
         ┌─┐ ┌─┐ ┌─┐
         │1│ │2│ │3│
         └─┘ └─┘ └─┘
          │   │\    \
          │   ▼ ▼    ▼
          │  ┌─┐ ┌─┐ ┌─┐
          │  │4│ │5│ │7│
          │  └─┘ └─┘ └─┘
          │    \  /
          └─────▼▼
               ┌─┐
               │6│
               └─┘
```

Now there is a single entry point. The graph is ready for BFS.

---

## Slide 11 — Before and After: The Full Transformation

```
   INITIAL GRAPH              REVERSED + GHOST ROOT
   (from annotations)         (ready for traversal)

   6 → {1, 4, 5}     ───►    root → {1, 2, 3}
   4 → {2}                    1    → {6}
   5 → {2}                    2    → {4, 5}
   7 → {3}                    3    → {7}
   1 → {}                     4    → {6}
   2 → {}                     5    → {6}
   3 → {}
```

The ghost `root` acts only as a traversal anchor — it is never written to the generated file.

---

## Slide 12 — Topological Sort: Kahn's Algorithm

**Kahn's algorithm** uses BFS and tracks the **in-degree** (number of incoming edges) for each node. A node is ready to visit once its in-degree drops to zero — meaning all its predecessors have been processed.

```
In-degrees (reversed + root graph):
  root=0  1=1  2=1  3=1  4=1  5=1  6=3  7=1

Queue  → [root]
Visit root   → decrement 1, 2, 3  →  queue: [1, 2, 3]
Visit 1      → decrement 6 (3→2)  →  queue: [2, 3]
Visit 2      → decrement 4, 5     →  queue: [3, 4, 5]
Visit 3      → decrement 7        →  queue: [4, 5, 7]
Visit 4      → decrement 6 (2→1)  →  queue: [5, 7]
Visit 5      → decrement 6 (1→0)  →  queue: [7, 6]
Visit 7      → (no outgoing)      →  queue: [6]
Visit 6      → done               →  queue: []
```

**Initialization order: 1 → 2 → 3 → 4 → 5 → 7 → 6** ✓

---

## Slide 13 — Why BFS Gives the Correct Order

Every time a node is dequeued:
- All nodes it depended on (its predecessors in the reversed graph = its dependents in the original) have already been dequeued
- Its in-degree is zero, confirming this

Therefore every constructor call in the generated file can safely reference variables declared on earlier lines — the graph guarantees it.

```
logger  := NewLogger()          // step 1 (no deps)
db      := NewDatabase()        // step 2 (no deps)
cfg     := NewConfig()          // step 3 (no deps)
repo    := NewUserRepo(db)      // step 4 (db ready ✓)
mailer  := NewMailer(db)        // step 5 (db ready ✓)
handler := NewHandler(cfg)      // step 7 (cfg ready ✓)
svc     := NewUserService(      // step 6
               logger, repo, mailer) // all ready ✓
```

---

## Slide 14 — Cycle Detection

If the dependency graph contains a **cycle**, no valid initialization order exists. Kahn's algorithm detects this naturally: some nodes will never reach in-degree zero and will never be enqueued.

```
  A ──► B
  ▲     │       A needs B, B needs C, C needs A.
  └──── C       None of them can be initialized first.
```

After the algorithm finishes, if the output contains fewer nodes than the graph, a cycle exists. The framework returns `ErrGraphHasCycle` immediately — the build fails before any code is written.

```go
var ErrGraphHasCycle = errors.New("Graph has cycle")
```

No runtime nil-pointer panics. No half-generated files.

---

## Slide 15 — From Graph to Generated AppCtx

The sorted node list is handed to a **Visitor** that walks it in order and emits the `AppCtx` file:

```go
for _, node := range sortedNodes {   // topological order
    visitor.Visit(node)              // emit one constructor call
}
```

Each `Visit` call emits:
- an **import** line for the node's package
- a **variable declaration** calling the constructor with already-declared dependencies
- a **struct field** and **getter** on `AppCtx`

Because the order is guaranteed by the graph, the generated file is always valid Go — no further analysis needed.

---

## Slide 16 — Full Pipeline: Graph Is the Core

```
  Annotated .go files
          │
          ▼
  ┌───────────────────┐
  │  Scan + Parse     │  @Injection → nodes
  │  (AST traversal)  │  @Inject    → edges
  └────────┬──────────┘
           │  initial dep graph
           ▼
  ┌───────────────────┐
  │  Reverse edges    │  A→B  becomes  B→A
  │  Add ghost root   │  root → {all in-degree-0 nodes}
  └────────┬──────────┘
           │  transformed graph
           ▼
  ┌───────────────────┐
  │  Kahn's BFS       │  topological order
  │  Cycle check      │  fail fast on cycle
  └────────┬──────────┘
           │  ordered node list
           ▼
  ┌───────────────────┐
  │  Visitor codegen  │  emit imports, fields, calls
  └────────┬──────────┘
           │
           ▼
     app_context.go ✓
```

Three graph operations — **build, transform, sort** — is all the framework does.
