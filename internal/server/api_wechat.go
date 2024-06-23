// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package server

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"ss-wecom-assistant/internal/util/history"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/xuri/excelize/v2"
)

var globalThreadID map[string]string
var imageID string
var lastUpdateTime int64

func init() {
	globalThreadID = make(map[string]string)
	imageID = ""
	lastUpdateTime = time.Now().Unix()
}

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
		threadID, ok := globalThreadID[toUser]
		if !ok {
			threadID, err = srv.svcs.ChatService.CreateThread(ctx, content, true)
			if err != nil {
				panic(err)
			}
			globalThreadID[toUser] = threadID
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

		reg := regexp.MustCompile(`【.*?】| (\*\*.+?\*\*)`)
		cleanedReply := reg.ReplaceAllString(reply, "")
		err = srv.svcs.WechatService.SendMsg(ctx, cleanedReply, toUser, openKFID, msgID)
		if err != nil {
			fmt.Println("senf msg error", err)
		}

		filePath := "/data/message_history.json"
		file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			fmt.Println("打开文件时出错:", err)
		}
		defer file.Close()

		histories, err := history.GetMessageHistory(file)
		if err != nil {
			fmt.Println("获取历史消息出错:", err)
		}

		newHistory := &history.MessageHistory{
			Question: content,
			Answer:   reply,
			UserId:   toUser,
		}

		histories = append(histories, *newHistory)

		emptyNickNameUser := []string{}

		for _, history := range histories {
			if history.NickName == "" {
				emptyNickNameUser = append(emptyNickNameUser, history.UserId)
			}
		}

		nickNameMap, err := srv.svcs.WechatService.BatchGetUserInfo(ctx, emptyNickNameUser)
		if err != nil {
			fmt.Println("获取用户信息出错")
		}

		for i, history := range histories {
			if history.NickName == "" {
				if newNickName, ok := nickNameMap[history.UserId]; ok {
					histories[i].NickName = newNickName
				}
			}
		}

		err = history.SaveMessage(histories, file)
		if err != nil {
			fmt.Println("save message error", err)
		}
	}()

	ctx.String(http.StatusOK, string(""))
}

func (srv *Server) getHistoryJson(ctx *gin.Context) {
	filePath := "/data/message_history.json"
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

	for _, history := range histories {
		if history.NickName == "" {
			emptyNickNameUser = append(emptyNickNameUser, history.UserId)
		}
	}

	nickNameMap, err := srv.svcs.WechatService.BatchGetUserInfo(ctx, emptyNickNameUser)
	if err != nil {
		fmt.Println("获取用户信息出错")
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

// Reply reply text message
func (srv *Server) getHistory(ctx *gin.Context) {
	filePath := "/data/message_history.json"
	file, err := os.OpenFile(filePath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		fmt.Println("打开文件时出错:", err)
	}
	defer file.Close()

	histories, err := history.GetMessageHistory(file)
	if err != nil {
		fmt.Println("获取历史消息出错:", err)
	}

	fileBytes, err := genExcel(histories)
	if err != nil {
		ctx.String(http.StatusInternalServerError, "Error generating Excel file")
		return
	}

	ctx.Header("Content-Disposition", "attachment; filename=people.xlsx")
	ctx.Data(http.StatusOK, "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet", fileBytes)
}

func genExcel(histories []history.MessageHistory) ([]byte, error) {
	f := excelize.NewFile()
	sheetName := "Sheet1"
	f.SetSheetName(f.GetSheetName(1), sheetName)

	headers := []string{"UserId", "Question", "Answer", "NickName"}
	for col, header := range headers {
		cell, _ := excelize.CoordinatesToCellName(col+1, 1)
		f.SetCellValue(sheetName, cell, header)
	}

	for row, history := range histories {
		cell, _ := excelize.CoordinatesToCellName(1, row+2)
		f.SetCellValue(sheetName, cell, history.UserId)
		cell, _ = excelize.CoordinatesToCellName(2, row+2)
		f.SetCellValue(sheetName, cell, history.Question)
		cell, _ = excelize.CoordinatesToCellName(3, row+2)
		f.SetCellValue(sheetName, cell, history.Answer)
		cell, _ = excelize.CoordinatesToCellName(4, row+2)
		f.SetCellValue(sheetName, cell, history.NickName)
	}

	var buf bytes.Buffer
	if err := f.Write(&buf); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
