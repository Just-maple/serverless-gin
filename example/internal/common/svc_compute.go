package common

import (
	"context"
)

type (
	Compute interface {
		Add(ctx context.Context, param ParamIn) (ret ParamRet, err error)

		Dec(ctx context.Context, param ParamIn) (ret ParamRet, err error)
	}

	ComputeSvc struct {
	}

	ParamIn struct {
		A int `form:"a"`
		B int `form:"b"`
	}

	ParamRet struct {
		Ret int `json:"ret"`
	}
)

func (s ComputeSvc) Dec(ctx context.Context, param ParamIn) (ret ParamRet, err error) {
	ret.Ret = param.A - param.B
	return
}

func (s ComputeSvc) Add(ctx context.Context, param ParamIn) (ret ParamRet, err error) {
	ret.Ret = param.A + param.B
	return
}
