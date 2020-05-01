package svrlessgin

import (
	"net/http"

	"github.com/Just-maple/svc2handler"
	"github.com/gin-gonic/gin"
)

type (
	GinSvcHandler func(svc interface{}) gin.HandlerFunc

	GinIOController interface {
		Response(c *gin.Context, ret interface{}, err error)

		ParamHandler(c *gin.Context, params []interface{}) bool
	}

	wrapGinIO struct {
		GinIOController
		c *gin.Context
	}
)

var (
	_ svc2handler.IOController = &wrapGinIO{}
)

func (io wrapGinIO) Response(w http.ResponseWriter, ret interface{}, err error) {
	io.GinIOController.Response(io.c, ret, err)
}

func (io wrapGinIO) ParamHandler(w http.ResponseWriter, r *http.Request, params []interface{}) (ok bool) {
	return io.GinIOController.ParamHandler(io.c, params)
}

func CreateGinIOController(ginIO GinIOController) GinSvcHandler {
	return func(svc interface{}) gin.HandlerFunc {
		return HandleSvcWithGinIO(ginIO, svc)
	}
}

func HandleSvcWithGinIO(io GinIOController, svc interface{}) gin.HandlerFunc {
	return func(c *gin.Context) {
		svc2handler.HandleSvcWithIO(&wrapGinIO{GinIOController: io, c: c}, svc)(c.Writer, c.Request)
	}
}
