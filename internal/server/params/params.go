// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package params

import (
	"ss-assistant/internal/consts"
	"ss-assistant/internal/model"
)

type EmbeddingSearchRequest struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
}

type CreateCollectionRequest struct {
	CollectionName string `json:"collection_name" form:"collection_name" binding:"required"`
}

type CreateKeywordFieldIndexRequest struct {
	Keyword string `json:"keyword" form:"keyword" binding:"required"`
}

type InsertMsgss-assistantRequest struct {
	ss-assistant []float32 `json:"ss-assistant" form:"ss-assistant" binding:"required"`
	MsgID  int64     `json:"msg_id" form:"msg_id" binding:"required"`
	UserID string    `json:"user_id" form:"user_id" binding:"required"`
}

type ss-assistantResponse struct {
	ss-assistant int `json:"ss-assistant" form:"ss-assistant"`
}

type LongTermMemorySession struct {
	UserID  string `json:"user_id" form:"user_id" binding:"required"`
	Content string `json:"content" form:"content" binding:"required"`
	Days    int    `json:"days" form:"days" binding:"required"`
}

type WechatCheckRequest struct {
	Signature string `json:"signature" form:"signature" binding:"required"`
	Timestamp string `json:"timestamp" form:"timestamp" binding:"required"`
	Nonce     string `json:"nonce" form:"nonce" binding:"required"`
	Echostr   string `json:"echostr" form:"echostr" binding:"required"`
}

type WechatChatRequest = model.Message

type WechatChatResponse struct {
	WechatChatReq,
	MsgType model.MessageType
}

type WechatUserInfoResponse struct {
	Subscribe      uint64   `json:"subscribe"`
	OpenID         string   `json:"openid"`
	Language       string   `json:"language"`
	SubscribeTime  uint64   `json:"subscribe_time"`
	UnionID        string   `json:"unionid"`
	Remark         string   `json:"remark"`
	GroupID        uint64   `json:"groupid"`
	TagIDList      []uint64 `json:"tagid_list"`
	SubscribeScene string   `json:"subscribe_scene"`
	QrScene        uint64   `json:"qr_scene"`
	QrSceneStr     string   `json:"qr_scene_str"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int64  `json:"expires_in"`
}

type ChatContent struct {
	System string
	User   string
}

type TagOpenIDList struct {
	Count int `json:"count"`
	Data  struct {
		OpenIDs []string `json:"openid"`
	} `json:"data"`
	NextOpenID string `json:"next_openid"`
}

type ApiKeyHeader struct {
	ApiKey string `header:"api_key"`
}

type BetchCreateJDKRequest struct {
	Level  consts.SubscribeLevel `json:"level" form:"level" binding:"required"`
	Amount int                   `json:"amount" form:"amount" binding:"required"`
}

type ListJDKRequest struct {
	Level consts.SubscribeLevel `uri:"level" binding:"required"`
}

type MysqlDumpRequest struct {
	Start uint64 `json:"start" form:"start" binding:"required"`
	Limit int    `json:"limit" form:"limit" binding:"required"`
}

type SystemPromptRequest struct {
	Role   string `json:"role" form:"role" binding:"required"`
	Prompt string `json:"prompt" form:"prompt" binding:"required"`
}

type EventReplyRequest struct {
	Event string `json:"event" form:"event" binding:"required"`
	Reply string `json:"reply" form:"reply" binding:"required"`
}
