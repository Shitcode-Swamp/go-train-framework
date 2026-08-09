package gin_adapter

import (
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/swag"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/controller/core"
)

const swaggerRoute = "/docs/*any"

var docsInstanceSeq atomic.Uint64

type Docs struct {
	options      core.DocsOptions
	instanceName string
	registry     *routeRegistry
}

func NewDocs(opt core.DocsOptions, registry *routeRegistry) *Docs {
	d := &Docs{
		options:      opt,
		instanceName: fmt.Sprintf("%s-%d", swag.Name, docsInstanceSeq.Add(1)),
		registry:     registry,
	}

	swag.Register(d.instanceName, d)
	return d
}

func (d *Docs) ReadDoc() string {
	spec, _ := json.Marshal(map[string]any{
		"swagger": "2.0",
		"info": map[string]string{
			"title":       d.options.Title,
			"description": d.options.Description,
			"version":     d.options.Version,
		},
		"paths": d.buildPaths(),
	})

	return string(spec)
}

func (d *Docs) buildPaths() map[string]any {
	paths := map[string]any{}

	for _, route := range d.registry.routes {
		path := swaggerPath(route.path)

		methods, ok := paths[path].(map[string]any)
		if !ok {
			methods = map[string]any{}
			paths[path] = methods
		}

		methods[strings.ToLower(route.method)] = route.operation()
	}

	return paths
}

func (d *Docs) register(engine *gin.Engine) {
	handler := ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.InstanceName(d.instanceName))

	if d.options.User != "" {
		engine.GET(swaggerRoute, gin.BasicAuth(gin.Accounts{d.options.User: d.options.Password}), handler)
		return
	}

	engine.GET(swaggerRoute, handler)
}
