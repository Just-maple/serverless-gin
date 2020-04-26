package svrlessgin

import (
	"context"
	"reflect"

	"github.com/gin-gonic/gin"
)

var (
	rTypeContext = reflect.TypeOf(new(context.Context)).Elem()
	rTypeError   = reflect.TypeOf(new(error)).Elem()
)

type (
	IOWrapper func(svc interface{}) gin.HandlerFunc

	IO interface {
		Response(c *gin.Context, ret interface{}, err error)

		ParamHandler(c *gin.Context, params []interface{})
	}

	adapter struct {
		svcV       reflect.Value
		funcNumIn  int
		funcNumOut int
		types      []reflect.Type
		retFunc    func(c *gin.Context, values []reflect.Value)
		kinds      []reflect.Kind
	}
)

func CreateIOWrapper(svr IO) IOWrapper {
	return func(svc interface{}) gin.HandlerFunc {
		return WrapSvcWithIO(svr, svc)
	}
}

func WrapSvcWithIO(io IO, svc interface{}) gin.HandlerFunc {
	v := reflect.ValueOf(svc)
	svcTyp := v.Type()
	funcNumOut := svcTyp.NumOut()
	funcNumIn := svcTyp.NumIn()
	ad := adapter{
		svcV:       v,
		funcNumIn:  funcNumIn,
		funcNumOut: svcTyp.NumOut(),
		types:      make([]reflect.Type, funcNumIn-1, funcNumIn-1),
		kinds:      make([]reflect.Kind, funcNumIn-1, funcNumIn-1),
	}
	if v.Kind() != reflect.Func {
		panic("invalid service func")
	}
	if funcNumIn < 1 {
		panic("service must has one or more param")
	}
	if svcTyp.In(0) != rTypeContext {
		panic("service first param must be context.Context")
	}
	switch funcNumOut {
	case 1:
		if svcTyp.Out(0) != rTypeError {
			panic("service last out must be error")
		}
		ad.retFunc = func(c *gin.Context, values []reflect.Value) {
			err, _ := values[0].Interface().(error)
			io.Response(c, nil, err)
		}
	case 2:
		if svcTyp.Out(1) != rTypeError {
			panic("service last out must be error")
		}
		ad.retFunc = func(c *gin.Context, values []reflect.Value) {
			err, _ := values[1].Interface().(error)
			io.Response(c, values[0].Interface(), err)
		}
	default:
		panic("service num out must be one or two")
	}
	for i := 1; i < funcNumIn; i++ {
		ad.types[i-1] = svcTyp.In(i)
		ad.kinds[i-1] = ad.types[i-1].Kind()
	}
	return ad.ginHandler(io)
}

func (ad *adapter) ginHandler(io IO) gin.HandlerFunc {
	return func(c *gin.Context) {
		newParamV := make([]reflect.Value, ad.funcNumIn, ad.funcNumIn)
		newParam := make([]interface{}, ad.funcNumIn-1, ad.funcNumIn-1)
		newParamV[0] = reflect.ValueOf(c.Request.Context())
		for i := 1; i < ad.funcNumIn; i++ {
			typ := ad.types[i-1]
			param := reflect.New(typ)
			if ad.kinds[i-1] == reflect.Ptr {
				param.Elem().Set(reflect.New(typ.Elem()))
			}
			newParamV[i] = param.Elem()
			newParam[i-1] = param.Interface()
		}
		io.ParamHandler(c, newParam)
		if c.IsAborted() {
			return
		}
		retValues := ad.svcV.Call(newParamV)
		ad.retFunc(c, retValues)
	}
}
