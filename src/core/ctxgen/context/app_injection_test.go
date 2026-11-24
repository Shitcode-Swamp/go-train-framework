package injection_context

import (
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/factory"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/initialization"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/injection"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/parsers"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/render"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/lexical"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/utils/packages"
)

type sourceFile struct {
	path    string
	pkgPath string
	src     string
}

type stubPackages struct {
	pkgs map[string]string
}

func (s *stubPackages) GetPackage(filePath string) (string, bool) {
	p, ok := s.pkgs[filePath]
	return p, ok
}

func generateAppContext(t *testing.T, files []sourceFile) string {
	t.Helper()

	pkgs := make(map[string]string, len(files))
	for _, f := range files {
		pkgs[f.path] = f.pkgPath
	}

	fac := factory.NewInjectorFactory(parsers.NewLexicalParser(), &stubPackages{pkgs: pkgs})
	tree := initialization.NewInitializationTree(nil)

	for _, f := range files {
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, f.path, f.src, parser.ParseComments)
		if err != nil {
			t.Fatalf("failed to parse %s: %v", f.path, err)
		}
		ast.Inspect(file, func(n ast.Node) bool {
			if injector, ok := fac.Create(n, file, f.path); ok {
				tree.Insert(injector)
			}
			return true
		})
	}

	if err := tree.Initialize(); err != nil {
		t.Fatalf("Initialize() error: %v", err)
	}

	renderer := render.NewInjectorRenderer("AppContext", "ctx")
	builder := NewInjectionContextBuilder(renderer, "AppContext")
	if err := tree.Traverse(builder); err != nil {
		t.Fatalf("Traverse() error: %v", err)
	}
	return builder.Build()
}

func assertValidGoSource(t *testing.T, src string) {
	t.Helper()
	if _, err := format.Source([]byte(src)); err != nil {
		t.Fatalf("generated output is not valid go source: %v\n--- output ---\n%s", err, src)
	}
}

func TestGenerate_SingleInjectorNoDependencies(t *testing.T) {
	out := generateAppContext(t, []sourceFile{
		{path: "/virtual/repo/db.go", pkgPath: "example.com/repo", src: `package repo

// @Injection(db)
func NewDB() *DB {
	return nil
}
`},
	})
	assertValidGoSource(t, out)

	alias := packages.Package{Name: "repo", Path: "example.com/repo"}.Alias()
	for _, want := range []string{
		"package app_context",
		alias + ` "example.com/repo"`,
		"db *" + alias + ".DB",
		"db := " + alias + ".NewDB()",
		"db: db,",
		"return ctx.db",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("output missing %q\n--- output ---\n%s", want, out)
		}
	}
}

func TestGenerate_DependencyChainAcrossPackages(t *testing.T) {
	files := []sourceFile{
		{path: "/virtual/repo/db.go", pkgPath: "example.com/repo", src: `package repo

// @Injection(db)
func NewDB() *DB {
	return nil
}
`},
		{path: "/virtual/repo/user_repo.go", pkgPath: "example.com/repo", src: `package repo

// @Injection(userRepo)
func NewUserRepo(
	db *DB, // @Inject(db)
) *UserRepo {
	return nil
}
`},
		{path: "/virtual/service/user_service.go", pkgPath: "example.com/service", src: `package service

// @Injection(userService, no-getter)
func NewUserService(
	userRepo *UserRepo, // @Inject(userRepo)
) *UserService {
	return nil
}
`},
	}

	out := generateAppContext(t, files)
	assertValidGoSource(t, out)

	repoAlias := packages.Package{Name: "repo", Path: "example.com/repo"}.Alias()
	serviceAlias := packages.Package{Name: "service", Path: "example.com/service"}.Alias()

	if got := strings.Count(out, repoAlias+` "example.com/repo"`); got != 1 {
		t.Errorf("expected repo import exactly once, got %d\n--- output ---\n%s", got, out)
	}
	if !strings.Contains(out, serviceAlias+` "example.com/service"`) {
		t.Errorf("missing service import\n--- output ---\n%s", out)
	}

	dbIdx := strings.Index(out, "db := ")
	userRepoIdx := strings.Index(out, "userRepo := ")
	userServiceIdx := strings.Index(out, "userService := ")
	if dbIdx == -1 || userRepoIdx == -1 || userServiceIdx == -1 {
		t.Fatalf("missing expected init statements\n--- output ---\n%s", out)
	}
	if !(dbIdx < userRepoIdx && userRepoIdx < userServiceIdx) {
		t.Errorf("expected init order db, userRepo, userService; got offsets db=%d userRepo=%d userService=%d\n--- output ---\n%s",
			dbIdx, userRepoIdx, userServiceIdx, out)
	}

	if !strings.Contains(out, "userRepo := "+repoAlias+".NewUserRepo(db)") {
		t.Errorf("expected userRepo to be constructed with db\n--- output ---\n%s", out)
	}
	if !strings.Contains(out, "userService := "+serviceAlias+".NewUserService(userRepo)") {
		t.Errorf("expected userService to be constructed with userRepo\n--- output ---\n%s", out)
	}

	if strings.Contains(out, "return ctx.userService") {
		t.Errorf("expected no getter body for userService (no-getter flag)\n--- output ---\n%s", out)
	}
	if !strings.Contains(out, "return ctx.db") {
		t.Errorf("expected a getter for db\n--- output ---\n%s", out)
	}
	if !strings.Contains(out, "return ctx.userRepo") {
		t.Errorf("expected a getter for userRepo\n--- output ---\n%s", out)
	}
}

func TestGenerate_UntaggedParamsAreNotWired(t *testing.T) {
	out := generateAppContext(t, []sourceFile{
		{path: "/virtual/repo/db.go", pkgPath: "example.com/repo", src: `package repo

// @Injection(db)
func NewDB() *DB {
	return nil
}
`},
		{path: "/virtual/repo/cache.go", pkgPath: "example.com/repo", src: `package repo

// @Injection(cache)
func NewCache(
	ttl int, // no @Inject tag here
) *Cache {
	return nil
}
`},
	})
	assertValidGoSource(t, out)

	alias := packages.Package{Name: "repo", Path: "example.com/repo"}.Alias()
	if !strings.Contains(out, "cache := "+alias+".NewCache()") {
		t.Errorf("expected untagged param to be dropped from the constructor call\n--- output ---\n%s", out)
	}
}

func TestGenerate_MissingDependencyTokenProducesUncompilableCode(t *testing.T) {
	files := []sourceFile{
		{path: "/virtual/repo/db.go", pkgPath: "example.com/repo", src: `package repo

// @Injection(db)
func NewDB() *DB {
	return nil
}
`},
		{path: "/virtual/service/user_service.go", pkgPath: "example.com/service", src: `package service

// @Injection(userService)
func NewUserService(
	userRepo *UserRepo, // @Inject(userRepo)
) *UserService {
	return nil
}
`},
	}

	out := generateAppContext(t, files)
	assertValidGoSource(t, out)

	serviceAlias := packages.Package{Name: "service", Path: "example.com/service"}.Alias()
	if !strings.Contains(out, "userService := "+serviceAlias+".NewUserService(userRepo)") {
		t.Fatalf("expected the constructor call to reference the unresolved token verbatim, got:\n%s", out)
	}
	if strings.Contains(out, "userRepo := ") {
		t.Fatalf("expected no declaration for userRepo, so the reference above is to an undeclared variable, got:\n%s", out)
	}
}

func TestInjectionContextBuilder_Visit_PanicsOnNilInjector(t *testing.T) {
	defer func() {
		if r := recover(); r != ErrMissingInjector {
			t.Fatalf("expected panic %v, got %v", ErrMissingInjector, r)
		}
	}()
	NewInjectionContextBuilder(noopRenderer{}, "AppContext").Visit(nil)
}

func TestInjectionContextBuilder_Visit_PanicsOnMissingToken(t *testing.T) {
	defer func() {
		if r := recover(); r != ErrMissingToken {
			t.Fatalf("expected panic %v, got %v", ErrMissingToken, r)
		}
	}()
	pkg := packages.Package{Name: "p", Path: "example.com/p"}
	NewInjectionContextBuilder(noopRenderer{}, "AppContext").
		Visit(injection.NewInjector("", "New", "Type", "New", pkg, nil, nil))
}

func TestInjectionContextBuilder_Visit_PanicsOnMissingTypeName(t *testing.T) {
	defer func() {
		if r := recover(); r != ErrMissingTypeName {
			t.Fatalf("expected panic %v, got %v", ErrMissingTypeName, r)
		}
	}()
	pkg := packages.Package{Name: "p", Path: "example.com/p"}
	NewInjectionContextBuilder(noopRenderer{}, "AppContext").
		Visit(injection.NewInjector("token", "New", "", "New", pkg, nil, nil))
}

type noopRenderer struct{}

func (noopRenderer) Import(*injection.Injector) lexical.Import { return lexical.Import{} }
func (noopRenderer) Field(*injection.Injector) lexical.Field   { return lexical.Field{} }
func (noopRenderer) Init(*injection.Injector) string           { return "" }
func (noopRenderer) Assign(*injection.Injector) string         { return "" }
func (noopRenderer) Getter(*injection.Injector) string         { return "" }
