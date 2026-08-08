package injection

import (
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/packages"
	"github.com/Shitcode-Swamp/go-train-framework/src/utils/sets"
)

const NoGetterFlag = "no-getter"

type Injector struct {
	name           string
	token          string
	pkg            packages.Package
	typeName       string
	initFn         string
	injectedFields []Injectee
	flags          sets.Set[string]
}

func NewInjector(
	token string,
	name string,
	typeName string,
	initFn string,
	pkg packages.Package,
	injectees []Injectee,
	flags []string,
) *Injector {
	return &Injector{
		token:          token,
		name:           name,
		typeName:       typeName,
		initFn:         initFn,
		pkg:            pkg,
		injectedFields: injectees,
		flags:          sets.NewSet(flags...),
	}
}

func (i *Injector) Token() string {
	return i.token
}

func (i *Injector) Name() string {
	return i.name
}

func (i *Injector) Injectees() []Injectee {
	return i.injectedFields
}

func (i *Injector) Package() packages.Package {
	return i.pkg
}

func (i *Injector) TypeName() string {
	return i.typeName
}

func (i *Injector) InitFn() string {
	return i.initFn
}

func (i *Injector) HasFlag(flag string) bool {
	_, ok := i.flags[flag]
	return ok
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
