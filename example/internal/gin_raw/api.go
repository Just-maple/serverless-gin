package main

import (
	"net/http"

	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
)

func RegisterComputeService(group gin.IRoutes, svc common.Compute) {
	group.GET("/add", func(c *gin.Context) {
		var param common.ParamIn
		if err := c.Bind(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ret, err := svc.Add(c.Request.Context(), param)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, ret)
	})

	group.GET("/dec", func(c *gin.Context) {
		var param common.ParamIn
		if err := c.Bind(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		ret, err := svc.Dec(c.Request.Context(), param)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, ret)
	})
}

func RegisterAccountService(group gin.IRoutes, svc common.Account) {
	group.GET("/new", func(c *gin.Context) {
		var param common.ParamAccount
		if err := c.Bind(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err := svc.NewAccount(c.Request.Context(), param, common.IP(c.ClientIP()))
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, nil)
	})

	group.POST("/edit", func(c *gin.Context) {
		var param common.ParamAccount
		if err := c.Bind(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		err := svc.EditAccountPassword(c.Request.Context(), param)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, nil)
	})
}

func RegisterOrderService(group gin.IRoutes, svc common.Order) {
	group.PUT("/new", func(c *gin.Context) {
		var param common.ParamNewOrder
		if err := c.Bind(&param); err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		get, _ := c.Get("user_id")
		userID, ok := get.(common.UserID)
		if !ok {
			c.AbortWithStatus(http.StatusUnauthorized)
			return
		}
		ret, err := svc.New(c.Request.Context(), param, userID)
		if err != nil {
			c.AbortWithStatus(http.StatusBadRequest)
			return
		}
		c.JSON(http.StatusOK, ret)
	})
}
