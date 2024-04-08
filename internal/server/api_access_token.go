// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package server

import (
	"errors"
	"ss-wecom-assistant/internal/util/api"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (srv *Server) getAccessToken(ctx *gin.Context) {
	srv.logger.Info("http: ", zap.Any("accseeToken", ctx.Request.URL.String()))
	api_key := ctx.Query("api_key")

	secret := srv.config.SecretConfig.Admin
	if api_key != secret {
		err := errors.New("admin secret is not correct")
		srv.logger.ErrorCtx(ctx, "Check secret", zap.Error(err))
		api.ResponseErrors(ctx, err)
		return
	}

	accessToken, err := srv.svcs.WechatService.GetAccessToken()
	if err != nil {
		api.ResponseErrors(ctx, err)
		return
	}

	api.ResponseWithSuccess(ctx, 200, accessToken)
}
