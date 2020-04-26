package common

import (
	"context"
	"log"
)

type (
	Account interface {
		NewAccount(ctx context.Context, param ParamAccount, ip IP) (err error)

		EditAccountPassword(ctx context.Context, param ParamAccount) (err error)
	}

	AccountSvc struct {
	}

	ParamAccount struct {
		Username string
		Password string
	}
)

func (s AccountSvc) NewAccount(ctx context.Context, param ParamAccount, ip IP) (err error) {
	log.Print("new account", param.Password, param.Username, ip)
	return
}

func (s AccountSvc) EditAccountPassword(ctx context.Context, param ParamAccount) (err error) {
	log.Print("edit account", param.Password, param.Username)
	return
}
