# Serverless-GIN
> provide your service easier, focus on service not server.

## Feature

- [x] easier to use
- [x] focus on service not server
- [x] help standardized project params handle and response

## Todo

- [ ] interface annotation code generate
- [ ] interface annotation doc generate

## Installation


Get by running:
```sh
go get github.com/Just-maple/serverless-gin
```

## Usage example

look at [example](./example/internal)

there is some service interface in [common](./example/internal/common) 

actually in ours real project service modules and interface will much more than example

if your want to provide them as api by gin

you can use the [raw](./example/internal/gin_raw/api.go) way

like

```go
// this example provide only one api router
package order

import (
	"net/http"

	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
)

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

```

or 
  
use the clean and simply [svrless](./example/internal/svrless/api.go) way


```go
// this example provide all api routers
package main

import (
	"github.com/Just-maple/serverless-gin/example/internal/common"
	"github.com/gin-gonic/gin"
)

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
```  

per request may lost `5µs` cause it based on `reflect`

but you can save time `(api-count - 1) * time-waist-in-api-code` in your life
 
just use these time to do more funny things


