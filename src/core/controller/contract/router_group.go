package contract

import "github.com/gin-gonic/gin"

type Context = gin.Context

type HandlerFunc func(ctx *Context)

type RouterGroup interface {
	GET(path string, handlers ...HandlerFunc) Docs
	POST(path string, handlers ...HandlerFunc) Docs
	PUT(path string, handlers ...HandlerFunc) Docs
	DELETE(path string, handlers ...HandlerFunc) Docs
	PATCH(path string, handlers ...HandlerFunc) Docs

	Group(relativePath string) RouterGroup
	Use(middleware ...HandlerFunc) RouterGroup
}
