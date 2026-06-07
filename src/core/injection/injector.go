package injection

import (
	"fmt"
)

type Injector struct {
	name           *string
	token          *string
	packageName    *string
	injectedFields []Injectee
}

func NewInject(injectees []Injectee) *Injector {
	return &Injector{
		injectedFields: injectees,
	}
}

func (i *Injector) WithToken(token string) *Injector {
	i.token = &token
	return i
}

func (i *Injector) WithName(name string) *Injector {
	i.name = &name
	return i
}

func (i *Injector) WithPackage(packageName string) *Injector {
	i.packageName = &packageName
	return i
}

func (i *Injector) Token() *string {
	return i.token
}

func (i *Injector) Injectees() []Injectee {
	return i.injectedFields
}

func (i *Injector) Expose() {
	fmt.Printf("Injector: %s, package: %s\n", *i.name, *i.packageName)
	for _, t := range i.injectedFields {
		fmt.Printf("Injectee token: %v\n", t.token)
	}
}

type Injectee struct {
	token string
}

func (i Injectee) Token() string {
	return i.token
}

func NewInjectee(token string) Injectee {
	return Injectee{token: token}
}
