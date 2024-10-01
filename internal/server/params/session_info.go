package params

import (
	"ss-wecom-assistant/internal/model"
)

type ListSessionInfoReq struct {
	StartDate string `uri:"start_date" binding:"required"`
	EndDate   string `uri:"end_date"`
}

type ListSessionInfoResp struct {
	Total int                  `json:"total"`
	Data  []*model.SessionInfo `json:"data"`
}
