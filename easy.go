package svrlessgin

import (
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
)

var _ GinIOController = &easyGinController{}

type (
	easyGinController struct {
		errorStatus int
		logger      errLogger
	}

	EasyOption func(io *easyGinController)

	errLogger interface {
		Printf(msg string, args ...interface{})
	}

	easyRet struct {
		Data interface{} `json:"data"`
		OK   bool        `json:"ok"`
		Err  string      `json:"err,omitempty"`
	}

	Render struct {
		Render render.Render
		Code   int
	}
)

func WithErrorStatus(st int) EasyOption {
	return func(io *easyGinController) {
		io.errorStatus = st
	}
}

func WithLogger(logger errLogger) EasyOption {
	return func(io *easyGinController) {
		io.logger = logger
	}
}

func NewEasyController(opts ...EasyOption) GinSvcHandler {
	eio := &easyGinController{
		errorStatus: http.StatusBadRequest,
		logger:      log.New(os.Stdout, "[GIN] ", log.LstdFlags),
	}
	for _, o := range opts {
		o(eio)
	}
	return CreateGinIOController(eio)
}

func (l *easyGinController) ParamHandler(c *gin.Context, params []interface{}) {
	paramLen := len(params)
	switch {
	case paramLen >= 1:
		if err := c.Bind(params[0]); err != nil {
			l.logger.Printf("param bind error:%v", err)
			c.AbortWithStatus(l.errorStatus)
			return
		}
	}
}

func (l *easyGinController) Response(c *gin.Context, ret interface{}, err error) {
	res := easyRet{
		Data: ret,
		OK:   err == nil,
	}
	if err != nil {
		l.logger.Printf("param bind error:%v", err)
		res.Err = err.Error()
		c.JSON(l.errorStatus, &res)
		return
	}
	switch r := ret.(type) {
	case Render:
		c.Render(r.Code, r.Render)
	default:
		c.JSON(http.StatusOK, &res)
	}
}
