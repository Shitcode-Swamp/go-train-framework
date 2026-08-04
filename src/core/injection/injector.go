package injection

import "github.com/Shitcode-Swamp/go-train-framework/src/utils/packages"

type Injector struct {
	name           string
	token          string
	pkg            packages.Package
	typeName       string
	initFn         string
	injectedFields []Injectee
}

func NewInjector(
	token string,
	name string,
	typeName string,
	initFn string,
	pkg packages.Package,
	injectees []Injectee,
) *Injector {
	return &Injector{
		token:          token,
		name:           name,
		typeName:       typeName,
		initFn:         initFn,
		pkg:            pkg,
		injectedFields: injectees,
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

func (i *Injector) PackageName() string {
	return i.pkg.Name
}

func (i *Injector) PackagePath() string {
	return i.pkg.Path
}

func (i *Injector) TypeName() string {
	return i.typeName
}

func (i *Injector) InitFn() string {
	return i.initFn
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
