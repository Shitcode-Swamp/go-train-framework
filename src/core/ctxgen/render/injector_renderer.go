package render

import (
	"fmt"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/ctxgen/injection"
	"github.com/Shitcode-Swamp/go-train-framework/src/lexical"
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

func (r *InjectorRenderer) Import(injector *injection.Injector) lexical.Import {
	return lexical.NewImport(injector.Package().Alias(), injector.Package().Path)
}

func (r *InjectorRenderer) Field(injector *injection.Injector) lexical.Field {
	return lexical.NewField(injector.Token(), r.typeName(injector))
}

func (r *InjectorRenderer) Init(injector *injection.Injector) string {
	initFn := lexical.NewFunc(fmt.Sprintf("%s.%s", injector.Package().Alias(), injector.InitFn()), nil, nil)
	return fmt.Sprintf("%s := %s", injector.Token(), initFn.StringCall(r.constructorArgs(injector)...))
}

func (r *InjectorRenderer) Assign(injector *injection.Injector) string {
	return fmt.Sprintf("%s: %s,", injector.Token(), injector.Token())
}

func (r *InjectorRenderer) Getter(injector *injection.Injector) string {
	if injector.HasFlag(injection.NoGetterFlag) {
		return ""
	}
	fn := lexical.NewFunc(
		cases.Title(language.English).String(injector.Token()),
		nil,
		[]string{r.typeName(injector)},
		fmt.Sprintf("return %s.%s", r.receiverVar, injector.Token()),
	).WithReceiver(lexical.NewReceiver(r.receiverVar, "*"+r.receiverType))
	return fn.String()
}

func (r *InjectorRenderer) typeName(injector *injection.Injector) string {
	return fmt.Sprintf("*%s.%s", injector.Package().Alias(), injector.TypeName())
}

func (r *InjectorRenderer) constructorArgs(injector *injection.Injector) []string {
	injectees := injector.Injectees()
	args := make([]string, 0, len(injectees))
	for _, param := range injectees {
		args = append(args, param.Token())
	}
	return args
}
