package main

import (
	"net/http"

	svrlessgin "github.com/Just-maple/serverless-gin"
	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var _ svrlessgin.IO = MyIO{}

// must implement svrlessgin.IO
// all service request wrap with io will pass this two method
// ParamHandler define how your service function get param from http
// Response define how your service function response param to http
type MyIO struct{}

// param handler define how your application fill in empty values and do some global action before param reach service
// params return ptr of empty values you define in service function except context.Context
// and you can use switch params[i].(type) or c.Bind(params[i]) to unmarshal the http request content to your param
// for some examples:
//
// func (this Compute) Add (ctx context.Context,param ParamAdd) (total int,err error)
// params values will be []interface{*paramAdd}
//
// func (this Compute) Add (ctx context.Context,a ParamA,b ParamB) (total int,err error)
// params values will be []interface{*ParamA,*ParamB}
//
// func (this Compute) Add (ctx context.Context,a *Param) (total int,err error)
// params values will be []interface{**Param}
//
// func (this Compute) Add (ctx context.Context,a *Param,userID UserID) (total int,err error)
// params values will be []interface{**Param,*UserID}
//
// func (this Compute) Nothing (ctx context.Context) (err error)
// params values will be []interface{}
func (l MyIO) ParamHandler(c *gin.Context, params []interface{}) {
	paramLen := len(params)
	switch {
	default:
		fallthrough
	case paramLen >= 1:
		// param one always be main param from http body or query
		// so i need to bind auto by gin
		if err := c.Bind(params[0]); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		fallthrough
	case paramLen >= 2:
		for i := 1; i < paramLen; i++ {
			// param two can be some extra info such as ip or user_id from cookies or session or jwt
			switch typ := params[i].(type) {
			case *common.UserID:
				get, ok := c.Get("user_id")
				if !ok {
					c.AbortWithStatus(http.StatusUnauthorized)
					return
				}
				*typ = common.UserID(get.(string))
			case *common.IP:
				*typ = common.IP(c.ClientIP())
			}
		}
	}
}

// response define how your application handle the return value and error from service
// for some examples:
//
// func (this Compute) Add (ctx context.Context,param ParamAdd) (total int,err error)
// ret values will be int from Compute.Add
//
// func (this Compute) Add (ctx context.Context,a ParamA,b ParamB) (total ResCompute,err error)
// ret values will be ResCompute from Compute.Add
//
// func (this Compute) Add (ctx context.Context,a *Param) (err error)
// ret values will be nil from Compute.Add
//
// func (this Compute) Nothing (ctx context.Context) (err error)
// params values will be nil from Compute.Nothing
func (l MyIO) Response(c *gin.Context, ret interface{}, err error) {
	if err != nil {
		c.JSON(http.StatusBadRequest, err.Error())
		return
	}
	// you can use type assert to define how these type response
	switch r := ret.(type) {
	case Render:
		c.Render(r.Code, r.Render)
	default:
		c.JSON(http.StatusOK, ret)
	}
}

// define some type that will cause diff response and use it on IO.Response
type (
	Render struct {
		Render render.Render
		Code   int
	}
)