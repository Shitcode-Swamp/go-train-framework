package gin_adapter

import (
	"github.com/gin-gonic/gin"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/controller/core"
)

type Controller struct {
	*RouterGroup
	docs   *Docs
	engine *gin.Engine
}

func NewController() core.Controller {
	engine := gin.Default()

	return &Controller{
		RouterGroup: NewRouterGroup(&engine.RouterGroup),
		engine:      engine,
	}
}

func (c *Controller) Run(addr ...string) error {
	return c.engine.Run(addr...)
}

func (c *Controller) RunTLS(addr, certFile, keyFile string) error {
	return c.engine.RunTLS(addr, certFile, keyFile)
}

func (c *Controller) Docs(opt core.DocsOptions) core.Controller {
	c.docs = NewDocs(opt, c.RouterGroup.registry)
	c.docs.register(c.engine)

	return c
}
