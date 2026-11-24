package injection_context

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/v3/src/core/ctxgen/injection"
	"github.com/Shitcode-Swamp/go-train-framework/v3/src/lexical"
)

const (
	FileGenerated = "FILE IS AUTO GENERATED"
	DoNotEdit     = "DO NOT EDIT MANUALLY"
)

var (
	ErrMissingInjector = fmt.Errorf("missing injector")
	ErrMissingToken    = fmt.Errorf("missing token for injector")
	ErrMissingTypeName = fmt.Errorf("missing typeName for injector")
)

type Renderer interface {
	Import(injector *injection.Injector) lexical.Import
	Field(injector *injection.Injector) lexical.Field
	Init(injector *injection.Injector) string
	Assign(injector *injection.Injector) string
	Getter(injector *injection.Injector) string
}

type InjectionContextBuilder struct {
	renderer    Renderer
	contextType string
	injectees   []*injection.Injector
}

func NewInjectionContextBuilder(renderer Renderer, contextType string) *InjectionContextBuilder {
	return &InjectionContextBuilder{
		renderer:    renderer,
		contextType: contextType,
	}
}

func (icb *InjectionContextBuilder) Visit(injector *injection.Injector) {
	if injector == nil {
		panic(ErrMissingInjector)
	}
	if injector.Token() == "" {
		panic(ErrMissingToken)
	}
	if injector.TypeName() == "" {
		panic(ErrMissingTypeName)
	}
	icb.injectees = append(icb.injectees, injector)
}

func (icb *InjectionContextBuilder) Build() string {
	sb := strings.Builder{}
	sb.WriteString(lexical.NewComment(FileGenerated, DoNotEdit).String())
	sb.WriteString(lexical.NewPackage("app_context").String())
	sb.WriteString(fmt.Sprintf("import (\n%s\n)\n", icb.buildImports()))
	sb.WriteString(icb.buildAppCtx())
	return sb.String()
}

func (icb *InjectionContextBuilder) buildImports() string {
	seen := make(map[string]struct{}, len(icb.injectees))
	imports := make([]string, 0, len(icb.injectees))
	for _, injectee := range icb.injectees {
		imp := icb.renderer.Import(injectee).String()
		if _, ok := seen[imp]; ok {
			continue
		}
		seen[imp] = struct{}{}
		imports = append(imports, imp)
	}
	return strings.Join(imports, "\n")
}

func (icb *InjectionContextBuilder) buildAppCtx() string {
	appCtx := strings.Builder{}
	appCtx.WriteString(lexical.NewStruct(icb.contextType, icb.buildFields()...).String())
	appCtx.WriteString("\n")
	appCtx.WriteString(icb.buildCtxInit())
	appCtx.WriteString(icb.buildGetters())
	return appCtx.String()
}

func (icb *InjectionContextBuilder) buildFields() []lexical.Field {
	fields := make([]lexical.Field, 0, len(icb.injectees))
	for _, injectee := range icb.exposed() {
		fields = append(fields, icb.renderer.Field(injectee))
	}
	return fields
}

func (icb *InjectionContextBuilder) buildCtxInit() string {
	body := make([]string, 0, len(icb.injectees)*2+2)
	for _, injectee := range icb.injectees {
		body = append(body, icb.renderer.Init(injectee))
	}
	body = append(body, icb.discardUnused()...)
	body = append(body, fmt.Sprintf("return &%s{", icb.contextType))
	for _, injectee := range icb.exposed() {
		body = append(body, icb.renderer.Assign(injectee))
	}
	body = append(body, "}")
	fn := lexical.NewFunc(fmt.Sprintf("New%s", icb.contextType), nil, []string{"*" + icb.contextType}, body...)
	return fn.String()
}

func (icb *InjectionContextBuilder) exposed() []*injection.Injector {
	exposed := make([]*injection.Injector, 0, len(icb.injectees))
	for _, injectee := range icb.injectees {
		if injectee.HasGetter() {
			exposed = append(exposed, injectee)
		}
	}
	return exposed
}

// A no-getter dependency nobody injects is still constructed for its side
// effects, but its local variable would otherwise fail "declared and not used".
func (icb *InjectionContextBuilder) discardUnused() []string {
	injected := make(map[string]struct{}, len(icb.injectees))
	for _, injector := range icb.injectees {
		for _, injectee := range injector.Injectees() {
			injected[injectee.Token()] = struct{}{}
		}
	}
	discards := make([]string, 0)
	for _, injector := range icb.injectees {
		if _, ok := injected[injector.Token()]; ok || injector.HasGetter() {
			continue
		}
		discards = append(discards, fmt.Sprintf("_ = %s", injector.Token()))
	}
	return discards
}

func (icb *InjectionContextBuilder) buildGetters() string {
	getters := strings.Builder{}
	for _, injectee := range icb.injectees {
		getters.WriteString(icb.renderer.Getter(injectee))
	}
	return getters.String()
}
