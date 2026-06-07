package injection

import (
	"fmt"
	"strings"
)

type Package struct {
	name string
	path string
}

type Injector struct {
	name           *string
	token          *string
	pkg            Package
	typeName	   *string
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

func (i *Injector) WithPackage(name, path string) *Injector {
	i.pkg = Package{name: name, path: path}
	return i
}

func (i *Injector) Token() *string {
	return i.token
}

func (i *Injector) Injectees() []Injectee {
	return i.injectedFields
}

func (i *Injector) Import() string {
	return i.pkg.path
}

func (i *Injector) TypeName() *string {
	return i.typeName
}

func (i *Injector) Init() string {
	return fmt.Sprintf(
		"%s := %s.%s(%s)",
		*i.token, i.pkg.name, *i.name, i.buildDeps())
}

func (i *Injector) buildDeps() string {
	depsSb := strings.Builder{}
	if len(i.injectedFields) == 0 {
		return ""
	}
	for _, injectee := range i.injectedFields {
		depsSb.WriteString("\n")
		depsSb.WriteString(injectee.Token())
		depsSb.WriteString(",")
	}
	depsSb.WriteString("\n")
	return depsSb.String()
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
