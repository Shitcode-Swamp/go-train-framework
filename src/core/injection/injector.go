package injection

import "github.com/Shitcode-Swamp/go-train/src/utils/packages"

type Injector struct {
	name           *string
	token          *string
	pkg            packages.Package
	typeName       *string
	initFn         *string
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

func (i *Injector) WithTypeName(typeName string) *Injector {
	i.typeName = &typeName
	return i
}

func (i *Injector) WithInitFn(initFn string) *Injector {
	i.initFn = &initFn
	return i
}

func (i *Injector) WithPackage(name, path string) *Injector {
	i.pkg = packages.Package{Name: name, Path: path}
	return i
}

func (i *Injector) Token() *string {
	return i.token
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

func (i *Injector) TypeName() *string {
	return i.typeName
}

func (i *Injector) InitFn() *string {
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
