// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"os"
	"os/signal"
	"ss-wecom-assistant/internal/config"
	"ss-wecom-assistant/internal/datastore"
	"ss-wecom-assistant/internal/repo"
	"ss-wecom-assistant/internal/server"
	"ss-wecom-assistant/internal/services"
	"syscall"

	"github.com/ArtisanCloud/PowerWeChat/v3/src/work"
	"github.com/sashabaranov/go-openai"
)

func args_parse() {
	flag.StringVar(&config.Env, "e", "local", "Default using local environment configuration.")
	flag.Parse()
}

func handleSignals(sig os.Signal) (exitNow bool) {
	switch sig {
	case syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM:
		return true
	case syscall.SIGUSR1:
		return false
	}
	return false
}

func registerSignal(shutdown chan struct{}) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, []os.Signal{syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1}...)
	go func() {
		for sig := range c {
			if handleSignals(sig) {
				close(shutdown)
				return
			}
		}
	}()
}

func main() {
	shutdownChannel := make(chan struct{})
	registerSignal(shutdownChannel)
	args_parse()

	configFilePath := "/usr/conf/config.toml"
	config.Init(configFilePath)
	c := config.Get()

	datastore.Init(c)
	db := datastore.Get()

	openaiClient := openai.NewClient(c.OpenAIConfig.ApiKey)

	weComApp, err := work.NewWork(&work.UserConfig{
		CorpID:  c.WeChatConfig.AppID,     // 企业微信的app id，所有企业微信共用一个。
		AgentID: c.WeChatConfig.AgentID,   // 内部应用的app id
		Secret:  c.WeChatConfig.AppSecret, // 内部应用的app secret
		Token:   c.WeChatConfig.Token,
		AESKey:  c.WeChatConfig.EncodingAESKey,
		OAuth: work.OAuth{
			Callback: "https://wecom.artisan-cloud.com/callback", //
			Scopes:   nil,
		},
		HttpDebug: true,
	})

	svcs := &server.Services{
		WechatService: services.NewWeComService(weComApp, c.WeChatConfig.AgentID),
		ChatService:   services.NewChatService(openaiClient),
	}

	repos := &server.Repos{
		DB:          db.DB,
		SessionInfo: repo.NewSessionInfo(db.DB),
		User:        repo.NewUser(db.DB),
	}

	server, err := server.NewServer(c, svcs, repos)
	if err != nil {
		panic("Failed to build new server, err: " + err.Error())
	}

	if err := server.Run(); err != nil {
		panic("Failed to run the server, err: " + err.Error())
	}

	<-shutdownChannel
	server.Shutdown()
}
