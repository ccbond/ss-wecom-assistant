package services

import (
	"context"
	"net/http"
)

type ChatService interface {
	CreateMessage(context.Context, string, string) (string, error)
	CreateThread(context.Context, string, bool) (string, error)
	GetResponse(context.Context, string, string) (string, error)
	WaitOnRun(context.Context, string, string) error
	CreateRun(context.Context, string, string) (string, error)
}

type WechatService interface {
	Server(*http.Request) (*http.Response, error)
	GetAccessToken() (string, error)
	Notify(*http.Request) (string, string, string, string, error)
	SendMsg(context.Context, string, string, string, string) error
	TransKF(context.Context, string, string, string) error
	TransEWM(context.Context, string, string, string, string) error
}
