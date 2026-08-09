package core

type Controller interface {
	RouterGroup
	Docs(opt DocsOptions) Controller
	Run(addr ...string) error
	RunTLS(addr, certFile, keyFile string) error
}

