package injection

import (
	"errors"
	"fmt"
)

var ErrNoInjectToken = errors.New("params in injection initializer must contain @Inject(<token>) doc string")

type Injector struct {
	name           *string
	packageName    *string
	injectedFields []Injectee
}

func NewInject(injectees []Injectee) *Injector {
	return &Injector{
		injectedFields: injectees,
	}
}

func (i *Injector) WithName(name string) *Injector {
	i.name = &name
	return i
}

func (i *Injector) WithPackage(packageName string) *Injector {
	i.packageName = &packageName
	return i
}

func (i *Injector) Expose() {
	fmt.Printf("Injector: %s, package: %s\n", *i.name, *i.packageName)
	for _, t := range i.injectedFields {
		fmt.Printf("token: %v\n", t.token)
	}
}

type Injectee struct {
	token string
}

func NewInjectee(token string) Injectee {
	return Injectee{token: token}
}
