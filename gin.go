package svrlessgin

import (
	"context"
	"net/http"
	"sync"

	"github.com/Just-maple/svc2handler"
	"github.com/gin-gonic/gin"
)

type (
	ctxKey string

	GinSvcHandler func(svc interface{}) gin.HandlerFunc

	GinIOController interface {
		Response(c *gin.Context, ret interface{}, err error)

		ParamHandler(c *gin.Context, params []interface{}) bool
	}

	ginIOWrapper struct {
		GinIOController
		serviceFunc interface{}
		ginContext  *gin.Context
		handlerFunc http.HandlerFunc
	}
)

const ctxKeyServiceFunc = ctxKey("func")

var (
	_ svc2handler.IOController = &ginIOWrapper{}
)

func (io *ginIOWrapper) Response(w http.ResponseWriter, ret interface{}, err error) {
	io.GinIOController.Response(io.ginContext, ret, err)
}

func (io *ginIOWrapper) ParamHandler(w http.ResponseWriter, r *http.Request, params []interface{}) (ok bool) {
	return io.GinIOController.ParamHandler(io.ginContext, params)
}

func newGinIOWrapperPool(ginIO GinIOController, svc interface{}) sync.Pool {
	// try patch
	_ = svc2handler.HandleSvcWithIO(&ginIOWrapper{}, svc)
	return sync.Pool{
		New: func() interface{} {
			r := &ginIOWrapper{GinIOController: ginIO, serviceFunc: svc}
			r.handlerFunc = svc2handler.HandleSvcWithIO(r, r.serviceFunc)
			return r
		},
	}
}

func (io *ginIOWrapper) injectServiceFunc() {
	request := io.ginContext.Request
	*request = *request.WithContext(context.WithValue(request.Context(), ctxKeyServiceFunc, io.serviceFunc))
}

var serviceFuncInject = false

func SetServiceFuncInject(enable bool) {
	serviceFuncInject = enable
}

func GetServiceFunc(ctx context.Context) (ret interface{}, ok bool) {
	ret, ok = ctx.Value(ctxKeyServiceFunc).(interface{})
	return
}

func NewWithController(ginIO GinIOController) GinSvcHandler {
	return func(svc interface{}) gin.HandlerFunc {
		if f, ok := svc.(gin.HandlerFunc); ok {
			return f
		}
		wrapperPool := newGinIOWrapperPool(ginIO, svc)
		return func(c *gin.Context) {
			wrapper := wrapperPool.Get().(*ginIOWrapper)
			wrapper.ginContext = c
			if serviceFuncInject {
				wrapper.injectServiceFunc()
			}
			wrapper.handlerFunc(c.Writer, c.Request)
			wrapperPool.Put(wrapper)
		}
	}
}
