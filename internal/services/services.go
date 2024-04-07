package services

type ChatService interface {
	// CreateChatCompletion(string, string, string, int) (string, int, error)
	// JudgmentUserMessage(context.Context, string) (int, error)
	// CreateTranscription(context.Context, string) (string, error)
}

type WechatService interface {
	// GetServer(*http.Request, gin.ResponseWriter) *wxServer.Server
	// GetAccessToken() (string, error)
	// SendCustomerMessage(string, []string) error
}
