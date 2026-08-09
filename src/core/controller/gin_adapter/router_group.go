package gin_adapter

import (
	"net/http"
	gopath "path"

	"github.com/gin-gonic/gin"

	"github.com/Shitcode-Swamp/go-train-framework/src/core/controller/contract"
)

type RouterGroup struct {
	group    *gin.RouterGroup
	registry *routeRegistry
}

func NewRouterGroup(group *gin.RouterGroup) *RouterGroup {
	return newRouterGroup(group, &routeRegistry{})
}

func newRouterGroup(group *gin.RouterGroup, registry *routeRegistry) *RouterGroup {
	return &RouterGroup{group: group, registry: registry}
}

func (r *RouterGroup) fullPath(relativePath string) string {
	return gopath.Join(r.group.BasePath(), relativePath)
}

func toGinHandlers(handlers []contract.HandlerFunc) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, len(handlers))
	for i, h := range handlers {
		ginHandlers[i] = gin.HandlerFunc(h)
	}

	return ginHandlers
}

func (r *RouterGroup) GET(path string, handlers ...contract.HandlerFunc) contract.Docs {
	r.group.GET(path, toGinHandlers(handlers)...)

	return r.registry.add(http.MethodGet, r.fullPath(path))
}

func (r *RouterGroup) POST(path string, handlers ...contract.HandlerFunc) contract.Docs {
	r.group.POST(path, toGinHandlers(handlers)...)

	return r.registry.add(http.MethodPost, r.fullPath(path))
}

func (r *RouterGroup) PUT(path string, handlers ...contract.HandlerFunc) contract.Docs {
	r.group.PUT(path, toGinHandlers(handlers)...)

	return r.registry.add(http.MethodPut, r.fullPath(path))
}

func (r *RouterGroup) DELETE(path string, handlers ...contract.HandlerFunc) contract.Docs {
	r.group.DELETE(path, toGinHandlers(handlers)...)

	return r.registry.add(http.MethodDelete, r.fullPath(path))
}

func (r *RouterGroup) Group(relativePath string) contract.RouterGroup {
	return newRouterGroup(r.group.Group(relativePath), r.registry)
}

func (r *RouterGroup) Use(middleware ...contract.HandlerFunc) contract.RouterGroup {
	r.group.Use(toGinHandlers(middleware)...)

	return r
}
