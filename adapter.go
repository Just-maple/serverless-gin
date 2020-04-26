package svrlessgin

import (
	"context"
	"reflect"
	"sync"

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
		svcV           reflect.Value
		funcNumIn      int
		funcNumOut     int
		firstIsContext bool
		types          []reflect.Type
		retFunc        func(c *gin.Context, values []reflect.Value)
		kinds          []reflect.Kind
		newParamPool   sync.Pool
	}

	newParam struct {
		i []interface{}
		v []reflect.Value
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
		types:      make([]reflect.Type, funcNumIn, funcNumIn),
		kinds:      make([]reflect.Kind, funcNumIn, funcNumIn),
	}
	if v.Kind() != reflect.Func {
		panic("invalid service func")
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
	for i := 0; i < funcNumIn; i++ {
		ad.types[i] = svcTyp.In(i)
		ad.kinds[i] = ad.types[i].Kind()
	}
	ad.firstIsContext = ad.types[0] == rTypeContext
	return ad.ginHandler(io)
}

func (ad *adapter) ginHandler(io IO) gin.HandlerFunc {
	ad.newParamPool = sync.Pool{
		New: func() interface{} {
			return &newParam{
				i: make([]interface{}, ad.funcNumIn, ad.funcNumIn),
				v: make([]reflect.Value, ad.funcNumIn, ad.funcNumIn),
			}
		},
	}
	return func(c *gin.Context) {
		newP := ad.newParamPool.Get().(*newParam)
		newParamV := newP.v
		newParam := newP.i
		for i := 0; i < ad.funcNumIn; i++ {
			if i == 0 && ad.firstIsContext {
				continue
			}
			typ := ad.types[i]
			param := reflect.New(typ)
			if ad.kinds[i] == reflect.Ptr {
				param.Elem().Set(reflect.New(typ.Elem()))
			}
			newParamV[i] = param.Elem()
			newParam[i] = param.Interface()
		}
		if ad.firstIsContext {
			newParam = newParam[1:]
		}
		io.ParamHandler(c, newParam)
		if c.IsAborted() {
			ad.newParamPool.Put(newP)
			return
		}
		if ad.firstIsContext {
			newParamV[0] = reflect.ValueOf(c.Request.Context())
		}
		ad.retFunc(c, ad.svcV.Call(newParamV))
		ad.newParamPool.Put(newP)
	}
}
