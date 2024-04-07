// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package config

import (
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
	"github.com/SyntSugar/ss-infra-go/api/server"
	"github.com/joho/godotenv"
)

type Server struct {
	ServerName string `toml:"server_name"`
	ServerPort int    `toml:"server_port"`
}

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

type Redis struct {
	Port       int    `toml:"port"`
	User       string `toml:"user"`
	MaxRetries int    `toml:"max_retries"`
	Host       string
	Pwd        string
}

type Log struct {
	LogLevel string `toml:"log_level"`
}

type WeChat struct {
	AppID          string
	AgentID        string
	AgentSecret    string
	Token          string
	EncodingAESKey string
	AppSecret      string
}

type OpenAI struct {
	ApiKey string
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
	DbConfig     Db     `toml:"database"`
	RedisConfig  Redis  `toml:"redis"`
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

	c.DbConfig.Host = os.Getenv("DATABASE_HOST")
	c.DbConfig.Pwd = os.Getenv("DATABASE_SECRET")

	c.RedisConfig.Host = os.Getenv("REDIS_HOST")
	c.RedisConfig.Pwd = os.Getenv("REDIS_SECRET")

	c.WeChatConfig = WeChat{
		AppID:          os.Getenv("WECHAT_APP_ID"),
		Token:          os.Getenv("WECHAT_TOKEN"),
		EncodingAESKey: os.Getenv("WECHAT_ENCODEING_AES_KEY"),
		AppSecret:      os.Getenv("WECHAT_APP_SECRET"),
	}

	c.OpenAIConfig = OpenAI{
		ApiKey: os.Getenv("OPENAI_API_KEY"),
	}

	c.SecretConfig = Secret{
		Admin: os.Getenv("ADMIN_SECRET"),
	}
}

// Get get config.
func Get() *Config {
	return &c
}
