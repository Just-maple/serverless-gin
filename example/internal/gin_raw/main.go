package main

import (
	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
)

var (
	svr = gin.Default()
)

func main() {
	RegisterComputeService(svr.Group("compute"), common.ComputeSvc{})
	RegisterAccountService(svr.Group("account"), common.AccountSvc{})
	RegisterOrderService(svr.Group("order"), common.OrderSvc{})

	panic(svr.Run(":80"))
}
