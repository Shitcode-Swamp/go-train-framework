package core

type DocsOptions struct {
	Title       string
	Description string
	Version     string
	User        string
	Password    string
}


type Docs interface {
	RouteDocs() RouteDocs
}

type RouteDocs interface {
	Query(query any) RouteDocs
	Body(body any) RouteDocs
	Path(path any) RouteDocs
}