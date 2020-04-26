package main

import (
	"context"
	"errors"

	svrlessgin "github.com/Just-maple/serverless-gin"
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
