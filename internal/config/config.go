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
	KFID           string
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

type Db struct {
	Port     int    `toml:"port"`
	User     string `toml:"user"`
	Type     string `toml:"type"`
	DbName   string `toml:"db_name"`
	InitConn int    `toml:"init_conn"`
	MaxConn  int    `toml:"max_conn"`
	Host     string
	Pwd      string
}

type Config struct {
	API   *server.APICfg
	Admin *server.AdminCfg

	ServerConfig Server `toml:"server"`
	LogConfig    Log    `toml:"log"`
	DbConfig     Db     `toml:"database"`
	WeChatConfig WeChat
	OpenAIConfig OpenAI
	SecretConfig Secret
}

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

	c.DbConfig.Host = os.Getenv("DATABASE_HOST")
	c.DbConfig.Pwd = os.Getenv("DATABASE_SECRET")

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
		KFID:           os.Getenv("WECHAT_KFID"),
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
