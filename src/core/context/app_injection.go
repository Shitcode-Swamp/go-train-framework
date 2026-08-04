package injection_context

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/injection"
)

var (
	ErrMissingInjector = fmt.Errorf("missing injector")
	ErrMissingToken    = fmt.Errorf("missing token for injector")
	ErrMissingTypeName = fmt.Errorf("missing typeName for injector")
)

type Renderer interface {
	Import(injector *injection.Injector) string
	Field(injector *injection.Injector) string
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
	sb.WriteString("// FILE IS AUTO GENERATED\n")
	sb.WriteString("// DO NOT EDIT MANUALLY\n")
	sb.WriteString("package app_context\n")
	sb.WriteString(fmt.Sprintf("import (\n\t%s\n)\n", icb.buildImports()))
	sb.WriteString(icb.buildAppCtx())
	return sb.String()
}

func (icb *InjectionContextBuilder) buildImports() string {
	seen := make(map[string]struct{}, len(icb.injectees))
	imports := make([]string, 0, len(icb.injectees))
	for _, injectee := range icb.injectees {
		imp := icb.renderer.Import(injectee)
		if _, ok := seen[imp]; ok {
			continue
		}
		seen[imp] = struct{}{}
		imports = append(imports, imp)
	}
	return strings.Join(imports, "\n\t")
}

func (icb *InjectionContextBuilder) buildAppCtx() string {
	appCtx := strings.Builder{}
	appCtx.WriteString(fmt.Sprintf("type %s struct {\n", icb.contextType))
	appCtx.WriteString(icb.buildFields())
	appCtx.WriteString("}\n\n")
	appCtx.WriteString(icb.buildCtxInit())
	appCtx.WriteString(icb.buildGetters())
	return appCtx.String()
}

func (icb *InjectionContextBuilder) buildFields() string {
	fields := strings.Builder{}
	for _, injectee := range icb.injectees {
		fields.WriteString(icb.renderer.Field(injectee))
	}
	return fields.String()
}

func (icb *InjectionContextBuilder) buildCtxInit() string {
	ctxInit := strings.Builder{}
	ctxInit.WriteString(fmt.Sprintf("func New%s() *%s {\n", icb.contextType, icb.contextType))
	for _, injectee := range icb.injectees {
		ctxInit.WriteString(icb.renderer.Init(injectee))
	}
	ctxInit.WriteString(fmt.Sprintf("\treturn &%s{\n", icb.contextType))
	for _, injectee := range icb.injectees {
		ctxInit.WriteString(icb.renderer.Assign(injectee))
	}
	ctxInit.WriteString("\t}\n")
	ctxInit.WriteString("}\n")
	return ctxInit.String()
}

func (icb *InjectionContextBuilder) buildGetters() string {
	getters := strings.Builder{}
	for _, injectee := range icb.injectees {
		getters.WriteString(icb.renderer.Getter(injectee))
	}
	return getters.String()
}
