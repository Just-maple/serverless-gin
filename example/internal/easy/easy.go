package main

import (
	"context"
	"errors"
	"net/http"

	svrlessgin "github.com/Just-maple/serverless-gin"
	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/ginS"
)

var (
	easy = svrlessgin.NewEasyWrapper()
)

type Param struct {
	A int `form:"a"`
	B int `form:"b"`
}

func main() {
	// raw gin
	ginS.GET("add/raw", func(c *gin.Context) {
		type ret struct {
			Data int `json:"data"`
		}
		var param Param
		err := c.Bind(&param)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, ret{Data: param.A + param.B})
	})

	// easy
	ginS.GET("add", easy(func(param Param) (int, error) {
		return param.A + param.B, nil
	}))
	ginS.GET("dec", easy(func(param Param) (int, error) {
		return param.A - param.B, nil
	}))
	ginS.GET("multi", easy(func(param Param) (int, error) {
		return param.A * param.B, nil
	}))
	// if first param is context.Context
	// fill from gin.Context.Request.Context()
	ginS.GET("divide", easy(func(ctx context.Context, param Param) (float64, error) {
		if param.B == 0 {
			return 0, errors.New("b cannot be zero")
		}
		return float64(param.A) / float64(param.B), nil
	}))
	panic(ginS.Run(":80"))
}
