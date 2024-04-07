// Copyright (c) Syntsugar Labs, Inc.
// SPDX-License-Identifier: MIT

package main

import (
	"flag"
	"os"
	"os/signal"
	"ss-assistant/internal/config"
	"ss-assistant/internal/datastore"
	"ss-assistant/internal/server"
	"ss-assistant/internal/services"
	"syscall"

	"github.com/ArtisanCloud/PowerWeChat/src/work"
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

	configFilePath := "./conf/config.toml"
	config.Init(configFilePath)
	c := config.Get()

	println("config: ", c)

	datastore.Init(c)
	defer datastore.Close()
	db := datastore.Get()

	openaiClient := openai.NewClient(c.OpenAIConfig.ApiKey)

	weComApp, err := work.NewWork(&work.UserConfig{
		CorpID:  c.WeChatConfig.AppID,       // 企业微信的app id，所有企业微信共用一个。
		AgentID: c.WeChatConfig.AgentID,     // 内部应用的app id
		Secret:  c.WeChatConfig.AgentSecret, // 内部应用的app secret
		OAuth: work.OAuth{
			Callback: "https://wecom.artisan-cloud.com/callback", //
			Scopes:   nil,
		},
		HttpDebug: true,
	})

	repos := &server.Repos{
		EventReply:       repo.NewEventReply(db.DB),
		SensitiveWords:   repo.NewSensitiveWord(db.DB),
		SpecifiedMessage: repo.NewSpecifiedMessage(db.DB),
		Users:            repo.NewUser(db.DB),
		Usage:            repo.NewUsage(db.DB),
		Subscription:     repo.NewSubscription(db.DB),
		JDK:              repo.NewJDK(db.DB),
		SystemPrompt:     repo.NewSystemPrompt(db.DB),
	}

	svcs := &server.Services{
		WechatService: services.NewWechatService(weComApp),
		ChatService:   services.NewChatService(openaiClient),
	}

	server, err := server.NewServer(c, db, repos, svcs)
	if err != nil {
		panic("Failed to build new server, err: " + err.Error())
	}

	if err := server.Run(); err != nil {
		panic("Failed to run the server, err: " + err.Error())
	}

	<-shutdownChannel
	server.Shutdown()
}
