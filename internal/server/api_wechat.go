// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package server

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"ss-wecom-assistant/internal/model"
	"ss-wecom-assistant/internal/util/history"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var globalThreadID map[string]string

func init() {
	globalThreadID = make(map[string]string)
}

// WechatCheck wechat check
func (srv *Server) wechatCheck(ctx *gin.Context) {
	rs, err := srv.svcs.WechatService.Server(ctx.Request)
	if err != nil {
		srv.logger.ErrorCtx(ctx, "wechat check error", zap.Error(err))
		ctx.String(http.StatusInternalServerError, "wechat check error")
		return
	}
	text, _ := ioutil.ReadAll(rs.Body)
	srv.logger.Info(string(text))

	ctx.String(http.StatusOK, string(text))
}

// Reply reply text message
func (srv *Server) wechatReply(ctx *gin.Context) {
	question, toUser, msgID, openKFID, err := srv.svcs.WechatService.Notify(ctx.Request)
	if err != nil {
		srv.logger.ErrorCtx(ctx, "wechat notify error", zap.Error(err))
		ctx.String(http.StatusInternalServerError, "wechat notify error")
		return
	}

	go func() {
		tx := srv.repos.DB.Begin()

		threadID, ok := globalThreadID[toUser]
		if !ok {
			threadID, err = srv.svcs.ChatService.CreateThread(ctx, question, true)
			if err != nil {
				srv.logger.ErrorCtx(ctx, "create thread error", zap.Error(err))
				return
			}
			globalThreadID[toUser] = threadID
		}

		messageID, err := srv.svcs.ChatService.CreateMessage(ctx, threadID, question)
		if err != nil {
			srv.logger.ErrorCtx(ctx, "create message error", zap.Error(err))
			return
		}

		runID, err := srv.svcs.ChatService.CreateRun(ctx, threadID, srv.config.OpenAIConfig.AssistantID)
		if err != nil {
			srv.logger.ErrorCtx(ctx, "create run error", zap.Error(err))
			return
		}

		err = srv.svcs.ChatService.WaitOnRun(ctx, threadID, runID)
		if err != nil {
			srv.logger.ErrorCtx(ctx, "wait on run error", zap.Error(err))
			return
		}

		answer, err := srv.svcs.ChatService.GetResponse(ctx, threadID, messageID)
		if err != nil {
			srv.logger.ErrorCtx(ctx, "get response error", zap.Error(err))
			return
		}

		reg := regexp.MustCompile(`【.*?】| (\*\*.+?\*\*)`)
		cleanedReply := reg.ReplaceAllString(answer, "")
		err = srv.svcs.WechatService.SendMsg(ctx, cleanedReply, toUser, openKFID, msgID)
		if err != nil {
			srv.logger.ErrorCtx(ctx, "send msg error", zap.Error(err))
		}

		// get user by toUser
		var user model.User
		err = tx.Model(&model.User{}).Where("open_id = ?", toUser).First(&user).Error
		if err != nil {
			srv.logger.ErrorCtx(ctx, "get user error", zap.Error(err))

			emptyNickNameUser := []string{toUser}
			nickNameMap, err := srv.svcs.WechatService.BatchGetUserInfo(ctx, emptyNickNameUser)
			if err != nil {
				srv.logger.ErrorCtx(ctx, "get user info error", zap.Error(err))
			}

			if newNickName, ok := nickNameMap[toUser]; !ok {
				srv.logger.ErrorCtx(ctx, "get user info error", zap.Error(err))
			} else {
				user = model.User{
					OpenID:        toUser,
					NickName:      newNickName,
					TotalMessages: 0,
				}
				err = tx.Model(&model.User{}).Create(&user).Error
				if err != nil {
					srv.logger.ErrorCtx(ctx, "create user error", zap.Error(err))
				}
			}
		}

		srv.logger.InfoCtx(ctx, "chat success session info", zap.Any("user", user), zap.Any("threadID", threadID), zap.Any("messageID", messageID), zap.Any("runID", runID), zap.Any("question", question), zap.Any("reply", answer))

		newHistory := &model.SessionInfo{
			Question: question,
			Answer:   answer,
			UserID:   toUser,
			NickName: user.NickName,
		}

		if err := srv.repos.SessionInfo.Create(newHistory); err != nil {
			tx.Rollback()
			srv.logger.ErrorCtx(ctx, "create session info error", zap.Error(err))
			return
		}

		user.TotalMessages++
		if err := tx.Model(&model.User{}).Where("open_id = ?", toUser).Update("total_messages", user.TotalMessages).Error; err != nil {
			tx.Rollback()
			srv.logger.ErrorCtx(ctx, "update user error", zap.Error(err))
			return
		}

		tx.Commit()
	}()

	ctx.String(http.StatusOK, string(""))
}

func (srv *Server) getHistoryJson(ctx *gin.Context) {
	filePath := "./data/message_history.json"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
	}
	defer file.Close()

	histories, err := history.GetMessageHistory(file)
	if err != nil {
		fmt.Println("获取历史消息出错:", err)
	}

	emptyNickNameUser := []string{}
	seen := make(map[string]bool)

	for _, history := range histories {
		if history.NickName == "" {
			if !seen[history.UserId] {
				emptyNickNameUser = append(emptyNickNameUser, history.UserId)
				seen[history.UserId] = true
			}
		}
	}

	nickNameMap, err := srv.svcs.WechatService.BatchGetUserInfo(ctx, emptyNickNameUser)
	if err != nil {
		fmt.Println("获取用户信息出错", err)
	}

	for i, history := range histories {
		if history.NickName == "" {
			if newNickName, ok := nickNameMap[history.UserId]; ok {
				histories[i].NickName = newNickName
			}
		}
	}

	ctx.JSON(http.StatusOK, histories)
}
