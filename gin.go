package svrlessgin

import (
	"net/http"
	"sync"

	"github.com/Just-maple/svc2handler"
	"github.com/gin-gonic/gin"
)

type (
	GinSvcHandler func(svc interface{}) gin.HandlerFunc

	GinIOController interface {
		Response(c *gin.Context, ret interface{}, err error)

		ParamHandler(c *gin.Context, params []interface{}) bool
	}

	GinIOWrapper struct {
		GinIOController
		HandlerFunc http.HandlerFunc
		ServiceFunc interface{}

		c *gin.Context
	}
)

var (
	_ svc2handler.IOController = &GinIOWrapper{}
)

func (io *GinIOWrapper) Response(w http.ResponseWriter, ret interface{}, err error) {
	io.GinIOController.Response(io.c, ret, err)
}

func (io *GinIOWrapper) ParamHandler(w http.ResponseWriter, r *http.Request, params []interface{}) (ok bool) {
	return io.GinIOController.ParamHandler(io.c, params)
}

func NewWithController(ginIO GinIOController) GinSvcHandler {
	return func(svc interface{}) gin.HandlerFunc {
		handlerPool := sync.Pool{New: func() interface{} {
			r := &GinIOWrapper{GinIOController: ginIO, ServiceFunc: svc}
			r.HandlerFunc = svc2handler.HandleSvcWithIO(r, r.ServiceFunc)
			return r
		}}
		return func(c *gin.Context) {
			r := handlerPool.Get().(*GinIOWrapper)
			r.c = c
			r.HandlerFunc(c.Writer, c.Request)
			handlerPool.Put(r)
		}
	}
}
