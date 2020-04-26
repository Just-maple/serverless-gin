package common

import (
	"context"
	"log"
)

type (
	Order interface {
		New(ctx context.Context, param ParamNewOrder, id UserID) (ret ParamRet, err error)
	}

	OrderSvc struct {
	}

	ParamNewOrder struct {
		Count  int
		ItemID int
	}
)

func (s OrderSvc) New(ctx context.Context, param ParamNewOrder, id UserID) (ret ParamRet, err error) {
	log.Print("new order", param.ItemID, param.Count, id)
	return
}
