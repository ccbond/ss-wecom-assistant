// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/SyntSugar/ss-infra-go/api/server"
	"github.com/joho/godotenv"
)

type Server struct {
	ServerName string `toml:"server_name"`
	ServerPort int    `toml:"server_port"`
}

type Log struct {
	LogLevel string `toml:"log_level"`
}

type WeChat struct {
	AppID          string
	AgentID        int
	Token          string
	EncodingAESKey string
	AppSecret      string
	ZJKFID         string
}

type OpenAI struct {
	ApiKey      string
	AssistantID string
}

type AzureAI struct {
	ApiKey string
}

type Secret struct {
	Admin string
}

var Env string

var c Config

type Config struct {
	API   *server.APICfg
	Admin *server.AdminCfg

	ServerConfig Server `toml:"server"`
	LogConfig    Log    `toml:"log"`
	WeChatConfig WeChat
	OpenAIConfig OpenAI
	SecretConfig Secret
}

// Init init config.
func Init(path string) {
	_, err := toml.DecodeFile(path, &c)
	if err != nil {
		panic(fmt.Sprintf("config parse fail: %v", err))
	}

	if Env == "local" {
		err := godotenv.Load(".env")
		if err != nil {
			panic("Error loading .env file")
		}
	}

	wechatAgentIDStr := os.Getenv("WECHAT_AGENT_ID")

	wechatAgentID, err := strconv.Atoi(wechatAgentIDStr)
	if err != nil {
		fmt.Println("Error converting WECHAT_AGENT_ID to an integer:", err)
		return
	}

	c.WeChatConfig = WeChat{
		AppID:          os.Getenv("WECHAT_APP_ID"),
		AgentID:        wechatAgentID,
		Token:          os.Getenv("WECHAT_TOKEN"),
		EncodingAESKey: os.Getenv("WECHAT_ENCODEING_AES_KEY"),
		AppSecret:      os.Getenv("WECHAT_APP_SECRET"),
		ZJKFID:         os.Getenv("WECHAT_ZJKFID"),
	}

	c.OpenAIConfig = OpenAI{
		ApiKey:      os.Getenv("OPENAI_API_KEY"),
		AssistantID: os.Getenv("OPENAI_ASSISTANT_ID"),
	}

	c.SecretConfig = Secret{
		Admin: os.Getenv("ADMIN_SECRET"),
	}
}

// Get get config.
func Get() *Config {
	return &c
}
