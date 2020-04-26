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
	ginS.GET("add", easy(func(ctx context.Context, param Param) (int, error) {
		return param.A + param.B, nil
	}))
	ginS.GET("dec", easy(func(ctx context.Context, param Param) (int, error) {
		return param.A - param.B, nil
	}))
	ginS.GET("multi", easy(func(ctx context.Context, param Param) (int, error) {
		return param.A * param.B, nil
	}))
	ginS.GET("divide", easy(func(ctx context.Context, param Param) (int, error) {
		if param.B == 0 {
			return 0, errors.New("b cannot be zero")
		}
		return param.A / param.B, nil
	}))

	panic(ginS.Run(":80"))
}
