package main

import (
	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
)

// as your project comes bigger and bigger
// you don't need to repeat your param binding and response
// it did reduce your works time on writing bind code or auth user
// just define your io rule and use it anywhere
//
// also those members be apart in project can focus on the service only
// you don't need to care any about code style in router or api returns standardize
// cause it's depends on your io.Response

func RegisterComputeService(group gin.IRoutes, svc common.Compute) {
	group.GET("/add", wrapper(svc.Add))
	group.GET("/dec", wrapper(svc.Dec))
}

func RegisterAccountService(group gin.IRoutes, svc common.Account) {
	group.GET("/new", wrapper(svc.NewAccount))
	group.POST("/edit", wrapper(svc.EditAccountPassword))
}

func RegisterOrderService(group gin.IRoutes, svc common.Order) {
	group.PUT("/new", wrapper(svc.New))
}
