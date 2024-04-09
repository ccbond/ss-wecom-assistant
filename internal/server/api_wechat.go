// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

// WechatCheck wechat check
func (srv *Server) wechatCheck(ctx *gin.Context) {
	rs, err := srv.svcs.WechatService.Server(ctx.Request)
	if err != nil {
		panic(err)
	}
	text, _ := ioutil.ReadAll(rs.Body)
	srv.logger.Info(string(text))

	ctx.String(http.StatusOK, string(text))
}

// Reply reply text message
func (srv *Server) wechatReply(ctx *gin.Context) {
	content, toUser, msgID, openKFID, err := srv.svcs.WechatService.Notify(ctx.Request)
	if err != nil {
		panic(err)
	}

	go func() {
		threadID, err := srv.svcs.ChatService.CreateThread(ctx, content, true)
		if err != nil {
			fmt.Println("create thread error", err)
			panic(err)
		}
		fmt.Println("threadID", threadID)

		messageID, err := srv.svcs.ChatService.CreateMessage(ctx, threadID, content)
		if err != nil {
			fmt.Println("create message error", err)
			panic(err)
		}
		fmt.Println("messageID", messageID)

		runID, err := srv.svcs.ChatService.CreateRun(ctx, threadID, srv.config.OpenAIConfig.AssistantID)
		if err != nil {
			fmt.Println("create run error", err)
			panic(err)
		}
		fmt.Println("runID", runID)

		err = srv.svcs.ChatService.WaitOnRun(ctx, threadID, runID)
		if err != nil {
			fmt.Println("wait on run error", err)
			panic(err)
		}
		fmt.Println("wait on run success")

		reply, err := srv.svcs.ChatService.GetResponse(ctx, threadID, messageID)
		if err != nil {
			fmt.Println("get response error", err)
			panic(err)
		}
		fmt.Println("reply", reply)

		if reply == "请您稍等，马上给您安排。" {
			content := "请联系我们的客服，马上帮您安排。"
			err = srv.svcs.WechatService.SendMsg(ctx, content, toUser, openKFID, msgID)
			if err != nil {
				fmt.Println("senf msg error", err)
				panic(err)
			}
			err = srv.svcs.WechatService.TransMP(ctx, "ivychenChenXinYu", toUser, openKFID, msgID)
			if err != nil {
				panic(err)
			}
		} else {
			err = srv.svcs.WechatService.SendMsg(ctx, reply, toUser, openKFID, msgID)
			if err != nil {
				fmt.Println("senf msg error", err)
				panic(err)
			}
		}
	}()

	ctx.String(http.StatusOK, string(""))
}
