package main

import (
	svrlessgin "github.com/Just-maple/serverless-gin"
	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
)

var (
	svr     = gin.New()
	wrapper = svrlessgin.CreateIOWrapper(MyIO{})
)

func main() {
	svr.Use(gin.Logger())
	RegisterComputeService(svr.Group("compute"), common.ComputeSvc{})
	RegisterAccountService(svr.Group("account"), common.AccountSvc{})
	RegisterOrderService(svr.Group("order"), common.OrderSvc{})

	panic(svr.Run(":80"))
}
