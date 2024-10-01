package server

import (
	"errors"
	"ss-wecom-assistant/internal/server/params"
	"ss-wecom-assistant/internal/util/api"
	"time"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (srv *Server) getMsgHistory(ctx *gin.Context) {
	srv.logger.Info("http: ", zap.Any("Get Msg History", ctx.Request.URL.String()))
	api_key := ctx.Query("api_key")

	secret := srv.config.SecretConfig.Admin
	if api_key != secret {
		err := errors.New("admin secret is not correct")
		api.ResponseErrors(ctx, err)
		return
	}

	req := params.ListSessionInfoReq{}
	if err := ctx.ShouldBindUri(&req); err != nil {
		api.ResponseErrors(ctx, err)
		return
	}

	// oarse 2024-01-01 to time.Time
	startTime, err := time.Parse("2006-01-02", req.StartDate)
	if err != nil {
		api.ResponseErrors(ctx, err)
		return
	}

	endTime, err := time.Parse("2006-01-02", req.EndDate)
	if err != nil {
		api.ResponseErrors(ctx, err)
		return
	}

	total, sessionInfos, err := srv.repos.SessionInfo.ListByTimeInterval(startTime, endTime)
	if err != nil {
		srv.logger.ErrorCtx(ctx, "GetAllEventReplys", zap.Error(err))
		api.ResponseErrors(ctx, err)
		return
	}

	response := params.ListSessionInfoResp{
		Total: total,
		Data:  sessionInfos,
	}

	api.ResponseWithOK(ctx, response)
}
