package injection_context

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	ErrMissingInjector = fmt.Errorf("missing injector")
	ErrMissingToken    = fmt.Errorf("missing token for injector")
	ErrMissingTypeName = fmt.Errorf("missing typeName for injector")
)

type InjectionContextBuilder struct {
	injectees []*injection.Injector
}

func NewInjectionContextBuilder() *InjectionContextBuilder {
	return &InjectionContextBuilder{}
}

func (icb *InjectionContextBuilder) Visit(injector *injection.Injector) {
	if injector == nil {
		panic(ErrMissingInjector)
	}
	if injector.Token() == nil {
		panic(ErrMissingToken)
	}
	if injector.TypeName() == nil {
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
	leanImports := make(map[string]struct{}, len(icb.injectees))
	for _, injectee := range icb.injectees {
		leanImports[injectee.PackagePath()] = struct{}{}
	}
	imports := make([]string, 0, len(leanImports))
	for imp := range leanImports {
		imports = append(imports, fmt.Sprintf("\"%s\"", imp))
	}
	return strings.Join(imports, "\n\t")
}

func (icb *InjectionContextBuilder) buildAppCtx() string {
	appCtx := strings.Builder{}
	appCtx.WriteString("type AppContext struct {\n")
	appCtx.WriteString(icb.buildFields())
	appCtx.WriteString("}\n\n")
	appCtx.WriteString(icb.buildCtxInit())
	appCtx.WriteString(icb.buildGetters())
	return appCtx.String()
}

func (icb *InjectionContextBuilder) buildFields() string {
	fields := strings.Builder{}
	for _, injectee := range icb.injectees {
		fields.WriteString(
			fmt.Sprintf("\t%s *%s.%s\n", 
				*injectee.Token(), 
				injectee.PackageName(), 
				*injectee.TypeName(),
			),
		)
	}
	return fields.String()
}

func (icb *InjectionContextBuilder) buildCtxInit() string {
	ctxInit := strings.Builder{}
	ctxInit.WriteString("func NewAppContext() *AppContext {\n")
	for _, injectee := range icb.injectees {
		constructorParams := icb.injecteeContructorParams(injectee)
		ctxInit.WriteString(fmt.Sprintf(
			"\t%s := %s.%s(%s)\n",
			*injectee.Token(),
			injectee.PackageName(),
			*injectee.InitFn(),
			constructorParams,
		))
	}
	ctxInit.WriteString("\treturn &AppContext{\n")
	for _, injectee := range icb.injectees {
		ctxInit.WriteString(fmt.Sprintf("\t\t%s: %s,\n", *injectee.Token(), *injectee.Token()))
	}
	ctxInit.WriteString("\t}\n")
	ctxInit.WriteString("}\n")
	return ctxInit.String()
}

func (icb *InjectionContextBuilder) injecteeContructorParams(injectee *injection.Injector) string {
	params := strings.Builder{}
	if injectee.Injectees() == nil || len(injectee.Injectees()) == 0 {
		return ""
	}
	for idx, param := range injectee.Injectees() {
		if idx == len(injectee.Injectees())-1 {
			params.WriteString(fmt.Sprintf("%s", param.Token()))
		} else {
			params.WriteString(fmt.Sprintf("%s, ", param.Token()))
		}
	}
	return params.String()
}

func (icb *InjectionContextBuilder) buildGetters() string {
	getters := strings.Builder{}
	for _, injectee := range icb.injectees {
		getters.WriteString(
			fmt.Sprintf("func (ctx *AppContext) %s() *%s.%s {\n",
				cases.Title(language.English).String(*injectee.Token()),
				injectee.PackageName(),
				*injectee.TypeName()),
		)
		getters.WriteString(
			fmt.Sprintf("\treturn ctx.%s\n", *injectee.Token()),
		)
		getters.WriteString("}\n")
	}
	return getters.String()
}
