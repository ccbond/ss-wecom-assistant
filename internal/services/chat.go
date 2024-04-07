// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package services

import (
	"github.com/sashabaranov/go-openai"
)

type chatService struct {
	oaClient *openai.Client
}

func NewChatService(oaClient *openai.Client) ChatService {
	return &chatService{
		oaClient: oaClient,
	}
}

// func (c *chatService) CreateChatCompletion(systemContent, userContent string, dialogues []*model.Dialogue, chatModel string, maxTokenLimit int) (string, int, error) {
// 	messages := make([]openai.ChatCompletionMessage, 0, len(dialogues)+2)
// 	messages = append(messages, openai.ChatCompletionMessage{
// 		Role:    openai.ChatMessageRoleSystem,
// 		Content: systemContent,
// 	})

// 	for i := len(dialogues) - 1; i >= 0; i-- {
// 		dialogue := dialogues[i]
// 		messages = append(messages, openai.ChatCompletionMessage{
// 			Role:    openai.ChatMessageRoleUser,
// 			Content: dialogue.UserContent,
// 		})

// 		messages = append(messages, openai.ChatCompletionMessage{
// 			Role:    openai.ChatMessageRoleAssistant,
// 			Content: dialogue.AssistantContent,
// 		})
// 	}

// 	messages = append(messages, openai.ChatCompletionMessage{
// 		Role:    openai.ChatMessageRoleUser,
// 		Content: userContent,
// 	})

// 	resp, err := c.oaClient.CreateChatCompletion(
// 		context.Background(),
// 		openai.ChatCompletionRequest{
// 			Model:       chatModel,
// 			Messages:    messages,
// 			Temperature: 1,
// 			MaxTokens:   maxTokenLimit,
// 		},
// 	)
// 	if err != nil {
// 		return "", 0, fmt.Errorf("%w: %v", ErrCreateChatCompletion, err)
// 	}

// 	return resp.Choices[0].Message.Content, resp.Usage.TotalTokens, nil
// }

// const judgmentSystemContent = `You are named nono.Analyze the following text sent by the user. \nReturn a json object with the following parameters:\n'\''mode'\''(string): There are two modes correspond to different situations respectively\n (1) '\''information_exchange'\'': When the user seeks knowledge and answers in a professional field.\n (2) '\''emotional_interaction'\'': When the conversation is primarily emotional in nature.\n'\''special_service'\''(string): There are five special_services:\n (1) '\''translate'\'': When the user indicates they want to translate some content.\n (2) '\''ielts'\'': When users seek help with the IELTS exam or IELTS study.\n (3) '\''paid'\'': When a user expresses concern about the details of the payment for your services. you are named nono.\n (4) '\''hello'\'': when user just say hello, hi etc.\n (5) '\''normal'\'': Others services all return normal.\n'\''difficult'\'' (bool): Analyze the difficulty level of the user'\''s question.\n (1) false: When user'\''s content just a simple knowledge-based question and answer or a casual conversation. \n (2) true: When user'\''s content involves complex knowledge reasoning and deeper level logical judgments.\n'\''need_context'\'' (bool): By analyzing the syntactic structure of this text, determine whether it requires contextual information:\n (1)Check Thematic Consistency: If the text revolves around a clear and self-contained theme without abrupt changes or unrelated information, additional context is usually not needed. However, if there are sudden shifts in topic or seemingly random information, context may be necessary.\n (2)Identify Implicit Assumptions: When a text uses specific terms or references specific events without explanation, it assumes prior knowledge, indicating a need for context. If the text provides sufficient explanations or background internally, no external context is needed.\n (3)Analyze Semantic Coherence: If sentences or paragraphs in the text are logically connected and understandable, no additional context is required. If sentences appear isolated without clear logical links, context might be needed.\n (4)Recognize Pronouns and Conjunctions: If the pronouns and conjunctions in the text clearly refer to elements within the text, no extra context is needed. If their references are unclear or point to external information, context is required.\n (5)Consider Genre and Format: Some genres (like academic papers, news reports) often assume specific background knowledge and might need context. In contrast, genres like novels, if they present a story independently and completely, don'\''t require additional context.\n'\''emotion'\''(string): Categorize the user'\''s message emotion into four types: positive, negative, neutral, or confused.\nThe final JSON output should be structured as:{mode: string, special_service: string, emotion: string, difficult: bool, need_context: bool}. `

// func (c *chatService) JudgmentUserMessage(ctx context.Context, content string) (*model.JudgmentResult, int, error) {
// 	resp, err := c.oaClient.CreateChatCompletion(
// 		ctx,
// 		openai.ChatCompletionRequest{
// 			Model: openai.GPT3Dot5Turbo1106,
// 			Messages: []openai.ChatCompletionMessage{
// 				{
// 					Role:    openai.ChatMessageRoleSystem,
// 					Content: judgmentSystemContent,
// 				},
// 				{
// 					Role:    openai.ChatMessageRoleUser,
// 					Content: content,
// 				},
// 			},
// 			Temperature: 0,
// 			ResponseFormat: &openai.ChatCompletionResponseFormat{
// 				Type: openai.ChatCompletionResponseFormatTypeJSONObject,
// 			},
// 		},
// 	)

// 	if err != nil {
// 		return nil, 0, fmt.Errorf("%w: %v", ErrCreateJudgementResponse, err)
// 	}

// 	assistantContent := resp.Choices[0].Message.Content

// 	var result model.JudgmentResult
// 	if err := json.Unmarshal([]byte(assistantContent), &result); err != nil {
// 		return nil, 0, fmt.Errorf("%w: %v", ErrUnmarshalJudgementResponse, err)
// 	}

// 	return &result, resp.Usage.TotalTokens, nil
// }

// func (c *chatService) CreateTranscription(ctx context.Context, filePath string) (string, error) {
// 	resp, err := c.oaClient.CreateTranscription(
// 		ctx,
// 		openai.AudioRequest{
// 			Model:       openai.Whisper1,
// 			FilePath:    filePath,
// 			Temperature: 0.7,
// 			Language:    "zh",
// 		},
// 	)

// 	if err != nil {
// 		return "", fmt.Errorf("%w: %v", ErrCreateTranscription, err)
// 	}

// 	return resp.Text, nil
// }
