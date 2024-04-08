// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package services

import (
	"context"
	"fmt"
	"time"

	openai "github.com/sashabaranov/go-openai"
)

type chatService struct {
	oaClient *openai.Client
}

func NewChatService(oaClient *openai.Client) ChatService {
	return &chatService{
		oaClient: oaClient,
	}
}

func (c *chatService) CreateMessage(ctx context.Context, threadID string, content string) (string, error) {
	msg, err := c.oaClient.CreateMessage(ctx, threadID, openai.MessageRequest{
		Role:     "user",
		Content:  content,
		FileIds:  nil,
		Metadata: nil,
	})
	if err != nil {
		return "", err
	}
	return msg.ID, err
}

func (c *chatService) CreateThread(ctx context.Context, content string, create bool) (string, error) {
	threadRequest := openai.ThreadRequest{
		Messages: []openai.ThreadMessage{
			{
				Role:    openai.ThreadMessageRoleUser,
				Content: content,
			},
		},
	}

	if create {
		threadRequest = openai.ThreadRequest{}
	}

	resp, err := c.oaClient.CreateThread(ctx, threadRequest)
	if err != nil {
		return "", err
	}

	return resp.ID, nil
}

func (c *chatService) GetResponse(ctx context.Context, threadID string, messageID string) (string, error) {
	limit := 1
	order := "asc"
	after := messageID
	resp, err := c.oaClient.ListMessage(ctx, threadID, &limit, &order, &after, nil)
	if err != nil {
		return "", err
	}

	fmt.Println("resp: ", resp)

	return resp.Messages[0].Content[0].Text.Value, nil
}

func (c *chatService) CreateRun(ctx context.Context, threadID string, assistantID string) (string, error) {
	run, err := c.oaClient.CreateRun(context.Background(), threadID, openai.RunRequest{
		AssistantID: assistantID,
	})
	if err != nil {
		return "", err
	}
	return run.ID, nil
}

func (c *chatService) WaitOnRun(ctx context.Context, threadID string, runID string) error {
	for {
		run, err := c.oaClient.RetrieveRun(ctx, threadID, runID)
		if err != nil {
			return err
		}
		if run.Status == openai.RunStatusQueued || run.Status == openai.RunStatusInProgress {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		if run.Status == openai.RunStatusCompleted {
			return nil
		}
	}
}
