package render

import (
	"fmt"
	"strings"

	"github.com/Shitcode-Swamp/go-train/src/core/injection"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type InjectorRenderer struct {
	receiverType string
	receiverVar  string
}

func NewInjectorRenderer(receiverType, receiverVar string) *InjectorRenderer {
	return &InjectorRenderer{
		receiverType: receiverType,
		receiverVar:  receiverVar,
	}
}

func (r *InjectorRenderer) Import(injector *injection.Injector) string {
	return fmt.Sprintf("\"%s\"", injector.PackagePath())
}

func (r *InjectorRenderer) Field(injector *injection.Injector) string {
	return fmt.Sprintf("\t%s *%s.%s\n", injector.Token(), injector.PackageName(), injector.TypeName())
}

func (r *InjectorRenderer) Init(injector *injection.Injector) string {
	return fmt.Sprintf(
		"\t%s := %s.%s(%s)\n",
		injector.Token(),
		injector.PackageName(),
		injector.InitFn(),
		r.constructorParams(injector),
	)
}

func (r *InjectorRenderer) Assign(injector *injection.Injector) string {
	return fmt.Sprintf("\t\t%s: %s,\n", injector.Token(), injector.Token())
}

func (r *InjectorRenderer) Getter(injector *injection.Injector) string {
	getter := strings.Builder{}
	getter.WriteString(fmt.Sprintf(
		"func (%s *%s) %s() *%s.%s {\n",
		r.receiverVar,
		r.receiverType,
		cases.Title(language.English).String(injector.Token()),
		injector.PackageName(),
		injector.TypeName(),
	))
	getter.WriteString(fmt.Sprintf("\treturn %s.%s\n", r.receiverVar, injector.Token()))
	getter.WriteString("}\n")
	return getter.String()
}

func (r *InjectorRenderer) constructorParams(injector *injection.Injector) string {
	injectees := injector.Injectees()
	if len(injectees) == 0 {
		return ""
	}
	tokens := make([]string, 0, len(injectees))
	for _, param := range injectees {
		tokens = append(tokens, param.Token())
	}
	return strings.Join(tokens, ", ")
}
