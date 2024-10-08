// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package services

import (
	"context"
	"net/http"

	"github.com/ArtisanCloud/PowerLibs/fmt"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/kernel/contract"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/accountService/message/request"
	"github.com/ArtisanCloud/PowerWeChat/v3/src/work/server/handlers/models"
)

var cursor = ""

type wechatService struct {
	weCom   *work.Work
	agentID int
}

func NewWeComService(weCom *work.Work, agentID int) WechatService {
	return &wechatService{
		weCom:   weCom,
		agentID: agentID,
	}
}

func (w *wechatService) GetAccessToken() (string, error) {
	return "", nil
}

func (w *wechatService) Server(req *http.Request) (*http.Response, error) {
	rs, err := w.weCom.Server.Serve(req)
	if err != nil {
		panic(err)
	}
	return rs, nil
}

func (w *wechatService) Notify(req *http.Request) (string, string, string, string, error) {
	ctx := context.Background()
	openKFID := "wkiHeVBgAAckgxIrimL83KL9E10LtY0w"

	token := ""

	_, err := w.weCom.Server.Notify(req, func(event contract.EventInterface) interface{} {
		fmt.Dump("event", event)
		if event.GetEvent() == models.CALLBACK_EVENT_KF_MSG_OR_EVENT && event.GetMsgType() == "text" {
			msg := models.EventKFMsgOrEvent{}
			err := event.ReadMessage(&msg)
			if err != nil {
				println(err.Error())
				return "error"
			}
			// content := string(msg.Content)

			token = msg.Token
		}
		return kernel.SUCCESS_EMPTY_RESPONSE
	})

	if err != nil {
		panic(err)
	}

	findEndMsg := false
	lastContent := ""
	toUser := ""
	msgID := ""

	for !findEndMsg {
		msg, err := w.weCom.AccountServiceMessage.SyncMsg(ctx, cursor, token, 100, 0, openKFID)
		if err != nil {
			panic(err)
		}
		fmt.Dump("msg", msg)

		msgListLen := len(msg.MsgList)
		if msgListLen < 100 {
			findEndMsg = true
			toUser = msg.MsgList[msgListLen-1].Get("external_userid").(string)
			msgID = msg.MsgList[msgListLen-1].Get("msgid").(string)
			lastContentInterface := msg.MsgList[msgListLen-1].Get("text")
			if lastContentMap, ok := lastContentInterface.(map[string]interface{}); ok {
				if lastContent, ok = lastContentMap["content"].(string); ok {
					fmt.Dump("Content:", lastContent)
				} else {
					fmt.Dump("no content")
				}
			} else {
				fmt.Dump("error: 'text' not map[string]interface{}")
			}
		} else {
			cursor = msg.NextCursor
		}
	}

	return lastContent, toUser, msgID, openKFID, nil
}

type MsgBusinessCard struct {
	UserID string `json:"userid"`
}

func (w *wechatService) SendMsg(ctx context.Context, content string, toUser string, openKFID string, msgID string) error {
	messages := &request.RequestAccountServiceSendMsg{
		ToUser:   toUser,
		OpenKfid: openKFID,
		MsgID:    msgID,
		MsgType:  "text",
		Text: &request.RequestAccountServiceMsgText{
			Content: content,
		},
	}
	res, err := w.weCom.AccountServiceMessage.SendMsg(ctx, messages)
	fmt.Dump("res", res)
	return err
}

func (w *wechatService) TransKF(ctx context.Context, oldOpenKFID string, newOpenKFID string, externalUserID string) error {
	state, err := w.weCom.AccountServiceState.Get(ctx, oldOpenKFID, externalUserID)
	if err != nil {
		return err
	}
	fmt.Dump("state", state)

	_, err = w.weCom.AccountServiceState.Trans(ctx, newOpenKFID, externalUserID, 3, state.ServicerUserID)
	return err
}

func (w *wechatService) TransEWM(ctx context.Context, mediaID string, toUser string, openKFID string, msgID string) error {
	messages := &request.RequestAccountServiceSendMsg{
		ToUser:   toUser,
		OpenKfid: openKFID,
		MsgID:    msgID,
		MsgType:  "file",
		File: &request.RequestAccountServiceMsgFile{
			MediaID: mediaID,
		},
	}

	res, err := w.weCom.AccountServiceMessage.SendMsg(ctx, messages)
	fmt.Dump("medieID", mediaID)
	fmt.Dump("res", res)
	return err
}

func (w *wechatService) BatchGetUserInfo(ctx context.Context, externalUserIDList []string) (map[string]string, error) {
	userNickNameMap := make(map[string]string)

	resp, err := w.weCom.AccountServiceCustomer.BatchGet(ctx, externalUserIDList)
	if err != nil {
		return nil, err
	}

	for _, user := range resp.CustomerList {
		userNickNameMap[user.ToHashMap().Get("external_userid").(string)] = user.ToHashMap().Get("nickname").(string)
	}

	return userNickNameMap, nil
}
