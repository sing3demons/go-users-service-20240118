package router

import (
	"github.com/gin-gonic/gin"
	"github.com/sing3demons/users/constant"
)

type IContext interface {
	QueryString(name string) string
	Param(key string) string

	JSON(code int, obj any)
	Body(obj any) error
	ReadBodyJSON(obj any) error

	SetAuthorization(value string)
	Set(key string, value any)
	Get(key string) (value any, exists bool)
	GetSessionId() string
	GetAuthorization() string
	AbortWithStatusJSON(code int, msg any)
	Next()
}

func (c *HTTPContext) GetSessionId() string {
	return c.Context.Writer.Header().Get(constant.XSessionId)
}

func (c *HTTPContext) GetHeader(key string) string {
	return c.Context.GetHeader(key)
}

func (c *HTTPContext) Set(key string, value any) {
	c.Context.Set(key, value)
}

func (c *HTTPContext) Get(key string) (value any, exists bool) {
	return c.Context.Get(key)
}

type HTTPContext struct {
	*Microservice
	*gin.Context
}

func NewContext(ms *Microservice, c *gin.Context) *HTTPContext {
	return &HTTPContext{ms, c}
}

func (c *HTTPContext) JSON(code int, obj any) {
	c.Context.JSON(code, obj)
}

func (ctx *HTTPContext) Body(obj any) error {
	err := ctx.Context.ShouldBind(&obj)
	if err != nil {
		return err
	}

	return nil
}
func (ctx *HTTPContext) ReadBodyJSON(obj any) error {
	err := ctx.Context.ShouldBindJSON(&obj)
	if err != nil {
		return err
	}
	return nil
}

func (ctx *HTTPContext) SetAuthorization(value string) {
	ctx.Context.Request.Header.Set("Authorization", "Bearer "+value)
}

func (ctx *HTTPContext) GetAuthorization() string {
	return ctx.Context.Request.Header.Get("Authorization")
}

func (c *HTTPContext) QueryString(name string) string {
	return c.Context.Query(name)
}
func (c *HTTPContext) Param(key string) string {
	return c.Context.Param(key)
}

func (c *HTTPContext) Next() {
	c.Context.Next()
}

func (c *HTTPContext) AbortWithStatusJSON(code int, msg any) {
	c.Context.AbortWithStatusJSON(code, msg)
}
